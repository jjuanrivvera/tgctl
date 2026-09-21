package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeToken is an obvious fake with a real token's SHAPE (numeric id, long hash) so the
// redaction patterns are exercised exactly as they would be in production. Never put a live
// token in a test, a fixture or a bug report.
const fakeToken = "123456789:AAFfakeFAKEfakeFAKEfakeFAKEfake-01"

func TestRedactSecrets(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{{
		name: "method call URL",
		in:   `Post "https://api.telegram.org/bot` + fakeToken + `/sendMessage": read: connection timed out`,
		want: `Post "https://api.telegram.org/bot123456789:<redacted>/sendMessage": read: connection timed out`,
	}, {
		name: "file download URL",
		in:   "https://api.telegram.org/file/bot" + fakeToken + "/photos/file_0.jpg",
		want: "https://api.telegram.org/file/bot123456789:<redacted>/photos/file_0.jpg",
	}, {
		name: "bare token outside a URL",
		in:   "TGCTL_TOKEN=" + fakeToken + " is invalid",
		want: "TGCTL_TOKEN=123456789:<redacted> is invalid",
	}, {
		name: "already redacted text is left alone",
		in:   "https://api.telegram.org/bot123456789:<redacted>/getMe",
		want: "https://api.telegram.org/bot123456789:<redacted>/getMe",
	}, {
		name: "ordinary id:value text is not mangled",
		in:   "chat_id:12345 message_id:678 — 09:30",
		want: "chat_id:12345 message_id:678 — 09:30",
	}, {
		name: "empty",
		in:   "",
		want: "",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, RedactSecrets(tc.in))
		})
	}
}

// A self-hosted Local Bot API Server may hand out a credential the generic pattern does not
// recognize; the authenticator knows its own token and masks it regardless of shape.
func TestBotTokenAuth_RedactSecrets_UnusualShape(t *testing.T) {
	a, err := NewBotTokenAuth("42:short+odd/shape==")
	require.NoError(t, err)

	got := a.RedactSecrets("Post \"http://localhost:8081/bot42:short+odd/shape==/getMe\": EOF")
	assert.NotContains(t, got, "short+odd")
	assert.Contains(t, got, "42:<redacted>")
}

// dropConnServer answers every request by dropping the connection mid-flight — the local
// stand-in for the transport failure that leaked a live token (issue #21).
func dropConnServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		// ErrAbortHandler aborts without a response and without a panic trace in the log.
		panic(http.ErrAbortHandler)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func fakeTokenClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	a, err := NewBotTokenAuth(fakeToken)
	require.NoError(t, err)
	return New(a,
		WithBaseURL(baseURL),
		WithRPS(0),
		WithRetryPolicy(retryPolicy{
			maxAttempts: 2,
			base:        time.Millisecond,
			max:         time.Millisecond,
			rng:         func() float64 { return 0 },
		}),
	)
}

// The regression test for issue #21: a transport-level failure must not put the token in the
// error a caller prints, logs, or hands back as an MCP tool result.
func TestClient_TransportError_DoesNotLeakToken(t *testing.T) {
	c := fakeTokenClient(t, dropConnServer(t).URL)

	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.Error(t, err)

	msg := err.Error()
	assert.NotContains(t, msg, fakeToken, "the bot token must never reach an error message")
	assert.NotContains(t, msg, "AAFfakeFAKE", "not even the hash on its own")
	assert.Contains(t, msg, "123456789:<redacted>", "the non-secret bot id stays, so the error still says which bot")
	assert.Contains(t, msg, "sendMessage", "the failing method is still named")

	// Formatting the error through every verb a caller might reach for stays clean, and so
	// does the *url.Error a caller could pull out with errors.As.
	assert.NotContains(t, fmt.Sprintf("%v %s %q", err, err, err), fakeToken)
	var ue *url.Error
	require.True(t, errors.As(err, &ue))
	assert.NotContains(t, ue.URL, fakeToken)
	assert.NotContains(t, ue.Error(), fakeToken)
}

func TestClient_DownloadFile_TransportError_DoesNotLeakToken(t *testing.T) {
	c := fakeTokenClient(t, dropConnServer(t).URL)

	_, err := c.DownloadFile(t.Context(), "photos/file_0.jpg", io.Discard)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), fakeToken, "the /file/bot<token>/ URL leaks the same secret")
	assert.Contains(t, err.Error(), "123456789:<redacted>")
}

// Redaction must not break error inspection: a cancelled context is still recognizable.
func TestClient_RedactError_PreservesCause(t *testing.T) {
	c := fakeTokenClient(t, dropConnServer(t).URL)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := c.Call(ctx, "getMe", nil, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_RedactError_NilAndClean(t *testing.T) {
	c := fakeTokenClient(t, "https://example.invalid")
	assert.NoError(t, c.redactError(nil))

	clean := errors.New("plain failure, no secrets")
	assert.Equal(t, clean, c.redactError(clean), "an error with nothing to redact is returned untouched")
}

// An error that carries the token outside a *url.Error is wrapped, not rewritten in place:
// the message is redacted while errors.Is still reaches the cause.
func TestClient_RedactError_WrapsNonURLError(t *testing.T) {
	c := fakeTokenClient(t, "https://example.invalid")

	cause := errors.New("boom")
	err := c.redactError(fmt.Errorf("auth failed for %s: %w", fakeToken, cause))
	assert.NotContains(t, err.Error(), fakeToken)
	assert.Contains(t, err.Error(), "123456789:<redacted>")
	assert.ErrorIs(t, err, cause)
}

// An Authenticator that does not implement secretRedactor still gets pattern redaction.
func TestClient_Redact_FallsBackToPatterns(t *testing.T) {
	c := New(stubAuth{})
	got := c.redact("https://api.telegram.org/bot" + fakeToken + "/getMe")
	assert.Equal(t, "https://api.telegram.org/bot123456789:<redacted>/getMe", got)
	assert.False(t, strings.Contains(got, fakeToken))
}

// stubAuth is a minimal Authenticator with no redaction capability of its own.
type stubAuth struct{}

func (stubAuth) RequestURL(base, method string) string    { return base + "/" + method }
func (stubAuth) RedactedURL(base, method string) string   { return base + "/" + method }
func (stubAuth) FileURL(base, path string) string         { return base + "/" + path }
func (stubAuth) RedactedFileURL(base, path string) string { return base + "/" + path }
func (stubAuth) Method() string                           { return "stub" }
