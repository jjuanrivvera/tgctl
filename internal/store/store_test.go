package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "profile.db")
	s, err := Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestOpen_CreatesDirAndPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not meaningful on Windows")
	}
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "messages", "default.db")
	s, err := Open(dbPath)
	require.NoError(t, err)
	defer func() { _ = s.Close() }()

	info, err := os.Stat(dbPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	dirInfo, err := os.Stat(filepath.Dir(dbPath))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm())
}

func TestOpen_IdempotentSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default.db")
	s1, err := Open(path)
	require.NoError(t, err)
	require.NoError(t, s1.Close())

	// Reopening an existing store must not fail or wipe the schema.
	s2, err := Open(path)
	require.NoError(t, err)
	defer func() { _ = s2.Close() }()
	require.NoError(t, s2.Record(t.Context(), Message{Direction: "out", ChatID: 1, Kind: "text", Text: "hi"}))
}

func TestPathFor(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		wantErr bool
	}{
		{"simple name", "default", false},
		{"hyphenated", "my-bot", false},
		{"empty rejected", "", true},
		{"traversal rejected", "../../etc/passwd", true},
		{"slash rejected", "a/b", true},
		{"dot rejected", ".", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PathFor("/config", tc.profile)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, filepath.Join("/config", "messages", tc.profile+".db"), got)
		})
	}
}

func TestRecord_And_Query(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	require.NoError(t, s.Record(ctx, Message{TS: base, Direction: "out", ChatID: 100, Kind: "text", Text: "first"}))
	require.NoError(t, s.Record(ctx, Message{TS: base.Add(time.Hour), Direction: "out", ChatID: 100, Kind: "photo", FileID: "ABC", MessageID: 2}))
	require.NoError(t, s.Record(ctx, Message{TS: base.Add(2 * time.Hour), Direction: "in", ChatID: 200, Kind: "text", Text: "other chat"}))

	t.Run("no filter returns all, newest first", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{})
		require.NoError(t, err)
		require.Len(t, got, 3)
		assert.Equal(t, "other chat", got[0].Text)
		assert.Equal(t, "first", got[2].Text)
	})

	t.Run("filter by chat", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{ChatID: 100})
		require.NoError(t, err)
		require.Len(t, got, 2)
		for _, m := range got {
			assert.EqualValues(t, 100, m.ChatID)
		}
	})

	t.Run("filter by kind", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{Kind: "photo"})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "ABC", got[0].FileID)
		assert.EqualValues(t, 2, got[0].MessageID)
	})

	t.Run("filter by since", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{Since: base.Add(90 * time.Minute)})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "other chat", got[0].Text)
	})

	t.Run("limit caps rows and defaults when <=0", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{Limit: 1})
		require.NoError(t, err)
		require.Len(t, got, 1)

		got, err = s.Query(ctx, Filter{Limit: -5})
		require.NoError(t, err)
		require.Len(t, got, 3) // negative limit falls back to DefaultLimit (50), not zero rows
	})

	t.Run("absent optional fields round-trip as zero values", func(t *testing.T) {
		got, err := s.Query(ctx, Filter{Kind: "text", ChatID: 100})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Zero(t, got[0].MessageID)
		assert.Empty(t, got[0].FileID)
		assert.Zero(t, got[0].ReplyToMessageID)
	})
}

func TestRecord_RawPayloadRoundTrips(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	raw := json.RawMessage(`{"message_id":7,"chat":{"id":1}}`)
	require.NoError(t, s.Record(ctx, Message{Direction: "out", ChatID: 1, Kind: "text", Text: "x", Raw: raw}))

	got, err := s.Query(ctx, Filter{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.JSONEq(t, string(raw), string(got[0].Raw))
}

func TestSearch_FTS(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	require.True(t, s.FTSEnabled(), "modernc.org/sqlite in use is expected to ship FTS5")

	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "deploy failed on staging"}))
	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "all systems nominal"}))
	require.NoError(t, s.Record(ctx, Message{ChatID: 2, Kind: "text", Text: "deploy succeeded"}))

	got, err := s.Search(ctx, "deploy", Filter{})
	require.NoError(t, err)
	require.Len(t, got, 2)

	got, err = s.Search(ctx, "deploy", Filter{ChatID: 1})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Contains(t, got[0].Text, "staging")
}

func TestSearch_LIKEFallback(t *testing.T) {
	s := openTest(t)
	s.fts = false // force the fallback path regardless of this build's FTS5 support
	ctx := t.Context()

	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "deploy failed on staging"}))
	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "all systems nominal"}))

	got, err := s.Search(ctx, "deploy", Filter{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Contains(t, got[0].Text, "staging")

	got, err = s.Search(ctx, "nomatch", Filter{})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestShow(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "hi", MessageID: 42}))

	m, ok, err := s.Show(ctx, 42)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "hi", m.Text)

	_, ok, err = s.Show(ctx, 999)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestShow_MostRecentWins(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	base := time.Now().UTC()
	require.NoError(t, s.Record(ctx, Message{TS: base, ChatID: 1, Kind: "text", Text: "old", MessageID: 5}))
	require.NoError(t, s.Record(ctx, Message{TS: base.Add(time.Minute), ChatID: 2, Kind: "text", Text: "new", MessageID: 5}))

	m, ok, err := s.Show(ctx, 5)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "new", m.Text, "message_id is only unique per-chat; Show returns the newest match")
}

func TestPrune(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	now := time.Now().UTC()
	require.NoError(t, s.Record(ctx, Message{TS: now.Add(-100 * time.Hour), ChatID: 1, Kind: "text", Text: "old"}))
	require.NoError(t, s.Record(ctx, Message{TS: now, ChatID: 1, Kind: "text", Text: "recent"}))

	n, err := s.Prune(ctx, 24*time.Hour)
	require.NoError(t, err)
	assert.EqualValues(t, 1, n)

	got, err := s.Query(ctx, Filter{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "recent", got[0].Text)

	// The FTS side table must be pruned too, or a stale rowid would dangle.
	got, err = s.Search(ctx, "old", Filter{})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPrune_NoneOld(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	require.NoError(t, s.Record(ctx, Message{Direction: "out", ChatID: 1, Kind: "text", Text: "recent"}))

	n, err := s.Prune(ctx, 24*time.Hour)
	require.NoError(t, err)
	assert.Zero(t, n)
}

func TestRecord_DefaultsTSWhenZero(t *testing.T) {
	s := openTest(t)
	ctx := t.Context()
	before := time.Now().UTC()
	require.NoError(t, s.Record(ctx, Message{ChatID: 1, Kind: "text", Text: "x"}))
	after := time.Now().UTC()

	got, err := s.Query(ctx, Filter{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.False(t, got[0].TS.Before(before.Add(-time.Second)))
	assert.False(t, got[0].TS.After(after.Add(time.Second)))
}

func TestOpen_InvalidDir(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("permission-denied simulation needs a non-root POSIX user")
	}
	parent := t.TempDir()
	blocked := filepath.Join(parent, "blocked")
	require.NoError(t, os.MkdirAll(blocked, 0o500)) // no write perm: MkdirAll of a child must fail
	defer func() { _ = os.Chmod(blocked, 0o700) }() // let t.TempDir cleanup remove it

	_, err := Open(filepath.Join(blocked, "sub", "default.db"))
	require.Error(t, err)
}
