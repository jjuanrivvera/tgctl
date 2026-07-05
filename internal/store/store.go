// Package store is tgctl's local message history: every outbound send (and, in polling/webhook
// mode, every inbound update) is recorded to a per-bot-profile SQLite database, so a restarted
// or compacted session — or any external tool — can answer "what did you send/receive, when, to
// whom" even though the Bot API itself exposes no history endpoint (issue #5).
//
// The driver is modernc.org/sqlite (pure Go, no cgo): tgctl has no cgo dependency today and
// GoReleaser cross-compiles linux/darwin/windows from a single toolchain (DECISIONS.md) — a
// cgo-based driver would break that. Store failures must never break a send: every write-path
// caller treats an error here as "log a warning and continue" (see commands/recorder.go).
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite" // registers the "sqlite" database/sql driver

	"github.com/jjuanrivvera/tgctl/internal/config"
)

// DefaultLimit is the default row cap for Query/Search (the issue's `tgctl log` default).
const DefaultLimit = 50

// Message is one recorded send/receive event — the row shape from issue #5's schema.
//
// Optional integer/string fields use the Go zero value (0 / "") to mean "not provided" rather
// than sql.Null*: every real Telegram id is positive and every real text/file_id is non-empty,
// so zero is an unambiguous "absent" sentinel and callers don't have to unwrap a Null* type.
type Message struct {
	ID               int64           `json:"id"`
	TS               time.Time       `json:"ts"`
	Direction        string          `json:"direction"` // "out" | "in"
	ChatID           int64           `json:"chat_id"`
	MessageID        int64           `json:"message_id,omitempty"`
	Kind             string          `json:"kind"` // text|photo|document|voice|edit|callback|...
	Text             string          `json:"text,omitempty"`
	FileID           string          `json:"file_id,omitempty"`
	ReplyToMessageID int64           `json:"reply_to_message_id,omitempty"`
	Raw              json.RawMessage `json:"raw,omitempty"`
}

// Filter narrows Query/Search. The zero value means "no constraint" on that field.
type Filter struct {
	ChatID int64
	Since  time.Time
	Kind   string
	Limit  int // <=0 → DefaultLimit
}

// Store is a per-profile SQLite message history.
type Store struct {
	db  *sql.DB
	fts bool // true when the SQLite build linked into the driver includes the FTS5 module
}

// PathFor returns the per-profile DB path under configDir: <configDir>/messages/<profile>.db.
//
// A profile name is validated by config.ValidateProfileName wherever it is CREATED (alias/config
// set), but the ACTIVE profile for a given invocation comes from --bot/$TGCTL_BOT, which is
// user input on every single call and is not otherwise validated. Re-validating here is the
// only thing standing between a crafted profile name and a path escape, so it is not optional:
// ValidateProfileName rejects '/' and '\', which is what actually makes the join below safe.
func PathFor(configDir, profile string) (string, error) {
	if err := config.ValidateProfileName(profile); err != nil {
		return "", fmt.Errorf("store path: %w", err)
	}
	return filepath.Join(configDir, "messages", profile+".db"), nil
}

// Open opens (creating if needed) the SQLite store at dbPath and initializes its schema
// idempotently. The parent directory is created 0700 and the file itself chmod'd 0600 — the
// same posture as config.Save, since this file can contain full message text.
func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return nil, fmt.Errorf("create store dir: %w", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	// One connection: modernc.org/sqlite serializes writers internally, and a single tgctl
	// invocation never needs concurrent connections. This avoids SQLITE_BUSY storms if two
	// tgctl processes touch the same profile at once — busy_timeout below handles the rest.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure store: %w", err)
	}
	// Exec above forces the file into existence (sql.Open itself is lazy), so chmod now.
	if err := os.Chmod(dbPath, 0o600); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("chmod store: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

// FTSEnabled reports whether Search uses FTS5 MATCH (true) or a LIKE scan fallback (false).
func (s *Store) FTSEnabled() bool { return s.fts }

const schema = `
CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ts TEXT NOT NULL,
	direction TEXT NOT NULL,
	chat_id INTEGER NOT NULL,
	message_id INTEGER,
	kind TEXT NOT NULL,
	text TEXT,
	file_id TEXT,
	reply_to_message_id INTEGER,
	raw TEXT
);
CREATE INDEX IF NOT EXISTS idx_messages_chat_ts ON messages(chat_id, ts);
`

func (s *Store) migrate() error {
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}
	s.fts = s.tryEnableFTS()
	return nil
}

// tryEnableFTS creates the FTS5 side table used by Search, if this build of modernc.org/sqlite
// includes the FTS5 module. Some minimal SQLite builds omit it; rather than fail Open entirely,
// Search degrades to a LIKE scan (issue #5's explicit fallback requirement) and FTSEnabled()
// reports which mode is active (surfaced by `tgctl doctor` / `tgctl log search --help`).
func (s *Store) tryEnableFTS() bool {
	_, err := s.db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(text)`)
	return err == nil
}

// selectMessage is the shared column list for every read query. Columns are always qualified
// with the messages. prefix because Search's FTS5 path joins against messages_fts, which has
// its own (unrelated) `text` column — qualifying avoids an ambiguous-column error there and
// costs nothing in the unjoined queries.
const selectMessage = `SELECT messages.id, messages.ts, messages.direction, messages.chat_id,
	messages.message_id, messages.kind, messages.text, messages.file_id,
	messages.reply_to_message_id, messages.raw`

// Record inserts one message row. ts defaults to now (UTC) when zero. Every write goes through
// this one parameterized statement — no query is ever built by concatenating a value into SQL
// text, so there is no injection surface here regardless of what a chat's text/title contains.
func (s *Store) Record(ctx context.Context, m Message) error {
	if m.TS.IsZero() {
		m.TS = time.Now().UTC()
	}
	var raw any
	if len(m.Raw) > 0 {
		raw = string(m.Raw)
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO messages (ts, direction, chat_id, message_id, kind, text, file_id, reply_to_message_id, raw)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.TS.UTC().Format(time.RFC3339Nano), m.Direction, m.ChatID,
		nullableInt(m.MessageID), m.Kind, nullableString(m.Text), nullableString(m.FileID),
		nullableInt(m.ReplyToMessageID), raw,
	)
	if err != nil {
		return fmt.Errorf("record message: %w", err)
	}
	if s.fts && m.Text != "" {
		if id, idErr := res.LastInsertId(); idErr == nil {
			// Best-effort: FTS indexing is a search convenience, never a reason to fail Record.
			_, _ = s.db.ExecContext(ctx, `INSERT INTO messages_fts (rowid, text) VALUES (?, ?)`, id, m.Text)
		}
	}
	return nil
}

// Query lists messages matching f, newest first.
func (s *Store) Query(ctx context.Context, f Filter) ([]Message, error) {
	where, args := f.whereClause()
	parts := []string{selectMessage, "FROM messages"}
	if where != "" {
		parts = append(parts, "WHERE", where)
	}
	parts = append(parts, "ORDER BY messages.ts DESC, messages.id DESC LIMIT ?")
	args = append(args, effectiveLimit(f.Limit))

	rows, err := s.db.QueryContext(ctx, strings.Join(parts, " "), args...)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanMessages(rows)
}

// Search full-text searches recorded text, newest first. It uses FTS5 MATCH when available
// (query supports FTS5's operators: AND/OR/NOT, prefix*, "phrases"); otherwise a plain
// substring LIKE scan, which has no operator support but always works.
//
// The query text is always assembled from a []string of static SQL fragments joined with
// strings.Join, never by concatenating a dynamic value into the SQL text itself — every
// caller-supplied value (q, and whereArgs from f) travels only as a `?` bind argument.
func (s *Store) Search(ctx context.Context, q string, f Filter) ([]Message, error) {
	where, whereArgs := f.whereClause()

	var parts []string
	var args []any
	if s.fts {
		parts = []string{selectMessage, "FROM messages",
			"JOIN messages_fts ON messages_fts.rowid = messages.id", "WHERE messages_fts MATCH ?"}
		args = append(args, q)
	} else {
		parts = []string{selectMessage, "FROM messages", "WHERE messages.text LIKE ?"}
		args = append(args, "%"+q+"%")
	}
	if where != "" {
		parts = append(parts, "AND", where)
		args = append(args, whereArgs...)
	}
	parts = append(parts, "ORDER BY messages.ts DESC, messages.id DESC LIMIT ?")
	args = append(args, effectiveLimit(f.Limit))

	rows, err := s.db.QueryContext(ctx, strings.Join(parts, " "), args...)
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanMessages(rows)
}

// Show returns the most recently recorded row for a Telegram message_id, and false if none
// exists. message_id is only unique per-chat in the Bot API, not globally, so when the same
// numeric id was used in two different chats this returns the newest match — good enough for
// tgctl's single-operator use case, and callers can disambiguate with Query(Filter{ChatID}).
func (s *Store) Show(ctx context.Context, messageID int64) (Message, bool, error) {
	row := s.db.QueryRowContext(ctx,
		selectMessage+" FROM messages WHERE messages.message_id = ? ORDER BY messages.id DESC LIMIT 1", messageID)
	m, err := scanMessage(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, false, nil
	}
	if err != nil {
		return Message{}, false, fmt.Errorf("show message: %w", err)
	}
	return m, true, nil
}

// Prune deletes messages recorded before now-olderThan and returns the row count removed.
func (s *Store) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan).Format(time.RFC3339Nano)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("prune: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // no-op once Commit succeeds

	if s.fts {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM messages_fts WHERE rowid IN (SELECT id FROM messages WHERE ts < ?)`, cutoff); err != nil {
			return 0, fmt.Errorf("prune fts index: %w", err)
		}
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM messages WHERE ts < ?`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("prune messages: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("prune: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("prune: %w", err)
	}
	return n, nil
}

// whereClause renders f as a parameterized SQL fragment (no leading "WHERE") plus its bind
// args. Only static column-name/operator fragments are ever concatenated; every value is a
// placeholder, so this stays injection-safe regardless of what a filter value contains.
func (f Filter) whereClause() (string, []any) {
	var parts []string
	var args []any
	if f.ChatID != 0 {
		parts = append(parts, "messages.chat_id = ?")
		args = append(args, f.ChatID)
	}
	if !f.Since.IsZero() {
		parts = append(parts, "messages.ts >= ?")
		args = append(args, f.Since.UTC().Format(time.RFC3339Nano))
	}
	if f.Kind != "" {
		parts = append(parts, "messages.kind = ?")
		args = append(args, f.Kind)
	}
	return strings.Join(parts, " AND "), args
}

func effectiveLimit(n int) int {
	if n <= 0 {
		return DefaultLimit
	}
	return n
}

// nullableInt turns the zero-value sentinel into a real SQL NULL, so an absent message_id/
// reply_to_message_id reads back as 0 (via scanMessage) rather than round-tripping a stray 0.
func nullableInt(v int64) any {
	if v == 0 {
		return nil
	}
	return v
}

func nullableString(v string) any {
	if v == "" {
		return nil
	}
	return v
}

// rowScanner is the common method both *sql.Row and *sql.Rows implement, so scanMessage serves
// Query/Search (many rows) and Show (one row) without duplicating the column list.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanMessage(rs rowScanner) (Message, error) {
	var (
		m         Message
		ts        string
		messageID sql.NullInt64
		text      sql.NullString
		fileID    sql.NullString
		replyTo   sql.NullInt64
		raw       sql.NullString
	)
	if err := rs.Scan(&m.ID, &ts, &m.Direction, &m.ChatID, &messageID, &m.Kind, &text, &fileID, &replyTo, &raw); err != nil {
		return Message{}, err
	}
	parsed, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return Message{}, fmt.Errorf("parse stored ts %q: %w", ts, err)
	}
	m.TS = parsed
	m.MessageID = messageID.Int64
	m.Text = text.String
	m.FileID = fileID.String
	m.ReplyToMessageID = replyTo.Int64
	if raw.Valid {
		m.Raw = json.RawMessage(raw.String)
	}
	return m, nil
}

func scanMessages(rows *sql.Rows) ([]Message, error) {
	out := []Message{}
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
