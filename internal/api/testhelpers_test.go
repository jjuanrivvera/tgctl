package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// newTestClient spins up an httptest server and returns a Client wired to it. Retries use a
// deterministic zero-jitter policy with millisecond backoff so failure-path tests are fast,
// and rate limiting is disabled.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth, err := NewBotTokenAuth("123456:AAH-fake-test-credential-value")
	require.NoError(t, err)

	return New(auth,
		WithBaseURL(srv.URL),
		WithRPS(0),
		WithRetryPolicy(retryPolicy{
			maxAttempts: 4,
			base:        time.Millisecond,
			max:         time.Millisecond,
			rng:         func() float64 { return 0 },
		}),
	)
}

// okJSON writes a successful Bot API envelope wrapping the given raw result JSON.
func okJSON(w http.ResponseWriter, result string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true,"result":` + result + `}`))
}

// errJSON writes a Bot API error envelope.
func errJSON(w http.ResponseWriter, status, code int, desc, paramsJSON string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body := `{"ok":false,"error_code":` + itoa(code) + `,"description":` + quote(desc)
	if paramsJSON != "" {
		body += `,"parameters":` + paramsJSON
	}
	body += `}`
	_, _ = w.Write([]byte(body))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

func quote(s string) string { return `"` + s + `"` }
