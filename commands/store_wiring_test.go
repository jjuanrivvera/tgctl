package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
