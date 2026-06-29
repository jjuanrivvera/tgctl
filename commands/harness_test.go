package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

// method routes a Bot API method name to a raw `result` JSON (or an error envelope).
type routes map[string]string

// newServer returns an httptest server that answers Bot API calls from the routes map. A
// method present in the map returns {"ok":true,"result":<value>}; an unknown method 404s
// with an error envelope so tests exercise the error path too.
func newServer(t *testing.T, r routes) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		method := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
		w.Header().Set("Content-Type", "application/json")
		if body, ok := r[method]; ok {
			_, _ = w.Write([]byte(`{"ok":true,"result":` + body + `}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"ok":false,"error_code":404,"description":"Not Found: method ` + method + `"}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// run executes tgctl with args against the given server, isolating config (a temp XDG dir)
// and the keyring (in-memory). A bot token is provided via env so commands authenticate.
func run(t *testing.T, srv *httptest.Server, args ...string) (string, string, error) {
	t.Helper()
	keyring.MockInit()
	t.Setenv("TGCTL_TOKEN", "123456:TESTHASHVALUE")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")

	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	full := append([]string{}, args...)
	if srv != nil {
		full = append(full, "--base-url", srv.URL)
	}
	root.SetArgs(full)
	err := root.ExecuteContext(t.Context())
	return out.String(), errb.String(), err
}

// runNoToken is like run but without a token in the environment (for auth/login tests).
func runNoToken(t *testing.T, srv *httptest.Server, stdin string, args ...string) (string, string, error) {
	t.Helper()
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")

	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetIn(strings.NewReader(stdin))
	full := append([]string{}, args...)
	if srv != nil {
		full = append(full, "--base-url", srv.URL)
	}
	root.SetArgs(full)
	err := root.ExecuteContext(t.Context())
	return out.String(), errb.String(), err
}

// runIn executes tgctl against a SHARED config dir (so multiple calls see each other's
// writes) without resetting the keyring — the caller controls keyring setup. token, if
// non-empty, is provided via env.
func runIn(t *testing.T, dir string, srv *httptest.Server, token string, args ...string) (string, string, error) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TGCTL_TOKEN", token)

	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetIn(strings.NewReader(""))
	full := append([]string{}, args...)
	if srv != nil {
		full = append(full, "--base-url", srv.URL)
	}
	root.SetArgs(full)
	err := root.ExecuteContext(t.Context())
	return out.String(), errb.String(), err
}

func mustContain(t *testing.T, s, sub string) {
	t.Helper()
	require.Contains(t, s, sub)
}
