package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/tgctl/internal/store"
)

// TestOpenStoreForWrite_InvalidProfileWarnsButDoesNotFailSend pins DECISIONS.md's core
// guarantee: a store failure must never break a send. --bot names containing '/' fail
// config.ValidateProfileName (store.PathFor), which is exactly the case a crafted --bot value
// could hit since --bot/$TGCTL_BOT is never otherwise validated before reaching the client.
func TestOpenStoreForWrite_InvalidProfileWarnsButDoesNotFailSend(t *testing.T) {
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7}}`})
	out, errb, err := run(t, srv, "message", "send", "--chat", "7", "--text", "hi", "--bot", "a/b")
	require.NoError(t, err, "an unusable store must never fail the send")
	mustContain(t, out, "1")
	mustContain(t, errb, "message store unavailable")
}

func TestOpenStoreForWrite_QuietSuppressesWarning(t *testing.T) {
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7}}`})
	_, errb, err := run(t, srv, "message", "send", "--chat", "7", "--text", "hi", "--bot", "a/b", "--quiet")
	require.NoError(t, err)
	assert.NotContains(t, errb, "message store unavailable")
}

// TestLog_InvalidProfileIsARealError pins the read-path/write-path asymmetry: unlike a send,
// `tgctl log` exists purely to read the store, so an unopenable store must fail the command
// instead of silently reporting "no messages" (commands/log.go's withReadStore).
func TestLog_InvalidProfileIsARealError(t *testing.T) {
	_, _, err := run(t, nil, "log", "--bot", "a/b")
	require.Error(t, err)
}

// TestDryRun_NeverCreatesStore pins the Windows CI regression fix: --dry-run makes no API
// call, so clientFromCmd must skip opening the store entirely (not just skip recording). This
// is the one that mattered for TestAllCommands_DryRun on Windows — dozens of --dry-run
// subtests each opened (and never closed) a SQLite handle under the same shared temp config
// dir, and Windows' t.TempDir() cleanup can't unlink a file a still-open handle points at.
func TestDryRun_NeverCreatesStore(t *testing.T) {
	dir := t.TempDir()
	_, _, err := runIn(t, dir, nil, "123456:TESTHASH", "message", "send", "--chat", "7", "--text", "hi", "--dry-run")
	require.NoError(t, err)

	// XDG_CONFIG_HOME=dir → config.Dir() is <dir>/tgctl → the store would live under
	// <dir>/tgctl/messages/<profile>.db (internal/store.PathFor).
	messagesDir := filepath.Join(dir, "tgctl", "messages")
	_, statErr := os.Stat(messagesDir)
	assert.True(t, os.IsNotExist(statErr), "dry-run must never create the messages/ store directory")
}

// TestMessageSend_ClosesStoreHandle pins the Client.Close() plumbing end to end: after a real
// (non-dry-run) send through the full command tree, the store's file handle must already be
// released, so a second store.Open on the exact same path — simulating anything else (another
// tgctl invocation, or a test's tempdir cleanup) touching that file right after — succeeds
// immediately rather than blocking or erroring.
func TestMessageSend_ClosesStoreHandle(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7}}`})
	_, _, err := runIn(t, dir, srv, "123456:TESTHASH", "message", "send", "--chat", "7", "--text", "hi")
	require.NoError(t, err)

	dbPath := filepath.Join(dir, "tgctl", "messages", "default.db")
	st, err := store.Open(dbPath)
	require.NoError(t, err, "the previous command's store handle must already be closed")
	require.NoError(t, st.Close())
}
