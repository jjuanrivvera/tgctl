package commands

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

// renderCmd builds a minimal command carrying just the flags render() reads, with -o json.
func renderCmd(out *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "o", "json", "")
	cmd.Flags().StringSlice("columns", nil, "")
	cmd.Flags().Bool("no-color", true, "")
	cmd.Flags().Bool("quiet", false, "")
	cmd.Flags().String("jq", "", "")
	cmd.SetOut(out)
	return cmd
}

func TestWebhookUpdateHandler(t *testing.T) {
	var out bytes.Buffer
	h := webhookUpdateHandler(renderCmd(&out), "s3cr3t", nil)

	post := func(secret, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		if secret != "" {
			req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
		}
		rr := httptest.NewRecorder()
		h(rr, req)
		return rr
	}

	t.Run("valid with secret renders the update", func(t *testing.T) {
		out.Reset()
		rr := post("s3cr3t", `{"update_id":100,"message":{"message_id":7,"text":"hi"}}`)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, out.String(), `"update_id": 100`)
	})
	t.Run("wrong secret is 403", func(t *testing.T) {
		assert.Equal(t, http.StatusForbidden, post("nope", `{"update_id":1}`).Code)
	})
	t.Run("invalid JSON is 400", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, post("s3cr3t", `not json`).Code)
	})
	t.Run("GET is 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		h(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})
}

func TestRoutePath(t *testing.T) {
	assert.Equal(t, "/", routePath(""))
	assert.Equal(t, "/bot", routePath("/bot"))
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = l.Close() }()
	return l.Addr().(*net.TCPAddr).Port
}

// TestWebhookListen_Lifecycle runs the real command: it registers the webhook against a mock,
// serves a posted update, shuts down on context cancel, and deletes the webhook on exit.
func TestWebhookListen_Lifecycle(t *testing.T) {
	keyring.MockInit()
	t.Setenv("TGCTL_TOKEN", "123456:TESTHASH")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")

	var setCalled, delCalled bool
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/setWebhook"):
			setCalled = true
		case strings.HasSuffix(r.URL.Path, "/deleteWebhook"):
			delCalled = true
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer tg.Close()

	port := freePort(t)
	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetArgs([]string{
		"webhook", "listen",
		"--port", strconv.Itoa(port),
		"--secret-token", "s",
		"--set-url", "https://tunnel.example.com/bot",
		"--delete-on-exit",
		"-o", "json",
		"--base-url", tg.URL,
	})

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- root.ExecuteContext(ctx) }()

	// Wait for the listener to bind, then deliver one update.
	base := "http://127.0.0.1:" + strconv.Itoa(port)
	require.Eventually(t, func() bool {
		resp, err := http.Get(base) //nolint:noctx // simple readiness poll
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return true
	}, 3*time.Second, 20*time.Millisecond)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodPost, base,
		strings.NewReader(`{"update_id":5,"message":{"text":"yo"}}`))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "s")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cancel() // Ctrl-C
	require.NoError(t, <-done)

	assert.True(t, setCalled, "setWebhook should have been called for --set-url")
	assert.True(t, delCalled, "deleteWebhook should have been called for --delete-on-exit")
	assert.Contains(t, out.String(), `"update_id": 5`)
}

// TestWebhookListen_RecordsInbound drives the real `webhook listen` command end to end and
// verifies a delivered update lands in the local store with direction='in' — the webhook
// counterpart to TestUpdatesGet_RecordsInbound (commands/log_test.go).
func TestWebhookListen_RecordsInbound(t *testing.T) {
	keyring.MockInit()
	t.Setenv("TGCTL_TOKEN", "123456:TESTHASH")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")

	port := freePort(t)
	root := NewRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetArgs([]string{
		"webhook", "listen",
		"--port", strconv.Itoa(port),
		"-o", "json",
	})

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- root.ExecuteContext(ctx) }()

	base := "http://127.0.0.1:" + strconv.Itoa(port)
	require.Eventually(t, func() bool {
		resp, err := http.Get(base) //nolint:noctx // simple readiness poll
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return true
	}, 3*time.Second, 20*time.Millisecond)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodPost, base,
		strings.NewReader(`{"update_id":9,"message":{"message_id":3,"chat":{"id":555},"text":"webhook hello"}}`))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	_ = resp.Body.Close()

	cancel()
	require.NoError(t, <-done)

	// Same process, same XDG_CONFIG_HOME/profile: `log` reads the store the listener just wrote.
	logRoot := NewRootCmd()
	var logOut bytes.Buffer
	logRoot.SetOut(&logOut)
	logRoot.SetArgs([]string{"log", "--chat", "555", "-o", "json"})
	require.NoError(t, logRoot.ExecuteContext(t.Context()))
	mustContain(t, logOut.String(), "webhook hello")
	mustContain(t, logOut.String(), `"direction": "in"`)
}

func TestWebhookListenExcludedFromMCP(t *testing.T) {
	sel := ophis.ExcludeCmdsContaining(excludedFromMCP...)
	cmd := findCmd(NewRootCmd(), "webhook", "listen")
	require.NotNil(t, cmd, "webhook listen must exist")
	assert.False(t, sel(cmd), "the blocking listen server must be excluded from the MCP tool surface")
}
