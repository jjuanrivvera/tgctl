package commands

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

// The end-to-end shape of issue #21: a real command, a real transport failure, and the error
// a user (or an MCP client reading this process's stderr) ends up with. The leak that
// prompted the fix happened on exactly this path — a read timeout during `message send`.
func TestCommand_TransportError_DoesNotLeakToken(t *testing.T) {
	const fakeToken = "987654321:AAFfakeFAKEfakeFAKEfakeFAKEfake-02"

	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler) // drop the connection mid-flight
	}))
	t.Cleanup(srv.Close)
	keyring.MockInit()

	_, stderr, err := runIn(t, t.TempDir(), srv, fakeToken,
		"message", "send", "--chat", "1", "--text", "hi")

	require.Error(t, err)
	assert.NotContains(t, err.Error(), fakeToken, "the token must not reach the command's error")
	assert.NotContains(t, stderr, fakeToken, "nor anything the command wrote to stderr")
	assert.Contains(t, err.Error(), "987654321:<redacted>")
}
