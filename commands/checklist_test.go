package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecklistSend_WireParams(t *testing.T) {
	t.Run("minimal checklist send", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendChecklist") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok": true,
					"result": map[string]any{
						"message_id": 42,
						"chat":       map[string]any{"id": 12345678},
						"checklist":  map[string]any{"title": "Groceries"},
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "checklist", "send",
			"--business-connection-id", "bc_123",
			"--chat", "12345678",
			"--checklist", `{"title":"Groceries","tasks":[{"id":1,"text":"Apples"}]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "42")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "12345678", got["chat_id"])
		cl, ok := got["checklist"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "Groceries", cl["title"])
		tasks, ok := cl["tasks"].([]any)
		require.True(t, ok)
		assert.Len(t, tasks, 1)
	})

	t.Run("checklist send with optional flags", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendChecklist") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok": true,
					"result": map[string]any{
						"message_id": 43,
						"chat":       map[string]any{"id": 12345678},
						"checklist":  map[string]any{"title": "Sprint Tasks"},
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "checklist", "send",
			"--business-connection-id", "bc_123",
			"--chat", "12345678",
			"--checklist", `{"title":"Sprint Tasks","tasks":[{"id":1,"text":"Issue #29"},{"id":2,"text":"Verify"}]}`,
			"--thread", "5",
			"--silent",
			"--protect-content",
			"--message-effect-id", "eff_1",
			"--reply-parameters", `{"message_id":10}`,
			"--reply-markup", `{"inline_keyboard":[[{"text":"Help","url":"https://example.com"}]]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "43")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "12345678", got["chat_id"])
		assert.EqualValues(t, 5, got["message_thread_id"])
		assert.Equal(t, true, got["disable_notification"])
		assert.Equal(t, true, got["protect_content"])
		assert.Equal(t, "eff_1", got["message_effect_id"])
		assert.NotNil(t, got["reply_parameters"])
		assert.NotNil(t, got["reply_markup"])
	})
}

func TestChecklistEdit_WireParams(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/editMessageChecklist") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 42,
					"chat":       map[string]any{"id": 12345678},
					"checklist":  map[string]any{"title": "Updated Groceries"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "checklist", "edit",
		"--business-connection-id", "bc_123",
		"--chat", "12345678",
		"--message-id", "42",
		"--checklist", `{"title":"Updated Groceries","tasks":[{"id":1,"text":"Apples"},{"id":2,"text":"Oranges"}]}`,
		"--reply-markup", `{"inline_keyboard":[[{"text":"Done","callback_data":"done"}]]}`,
	)
	require.NoError(t, err)
	assert.Contains(t, out, "42")
	assert.Equal(t, "bc_123", got["business_connection_id"])
	assert.Equal(t, "12345678", got["chat_id"])
	assert.EqualValues(t, 42, got["message_id"])
	cl, ok := got["checklist"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Updated Groceries", cl["title"])
	assert.NotNil(t, got["reply_markup"])
}

func TestChecklistValidation(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"send missing connection", []string{"checklist", "send", "--chat", "123", "--checklist", `{"title":"T","tasks":[{"id":1,"text":"A"}]}`}},
		{"send missing chat", []string{"checklist", "send", "--business-connection-id", "bc", "--checklist", `{"title":"T","tasks":[{"id":1,"text":"A"}]}`}},
		{"send missing checklist", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123"}},
		{"send invalid json", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", "not-json"}},
		{"send not object", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `["array"]`}},
		{"send empty title", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"","tasks":[{"id":1,"text":"A"}]}`}},
		{"send too long title", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"` + strings.Repeat("x", 256) + `","tasks":[{"id":1,"text":"A"}]}`}},
		{"send empty tasks", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"T","tasks":[]}`}},
		{"send task not object", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"T","tasks":["not-obj"]}`}},
		{"send task missing id", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"T","tasks":[{"text":"A"}]}`}},
		{"send task missing text", []string{"checklist", "send", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"T","tasks":[{"id":1,"text":""}]}`}},
		{"edit missing message-id", []string{"checklist", "edit", "--business-connection-id", "bc", "--chat", "123", "--checklist", `{"title":"T","tasks":[{"id":1,"text":"A"}]}`}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
		})
	}
}

func TestChecklistKindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	for _, tc := range []struct {
		path   string
		method string
	}{
		{"checklist send", "sendChecklist"},
		{"checklist edit", "editMessageChecklist"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			found := false
			for _, c := range APICommands() {
				if c.Path == tc.path {
					found = true
					assert.Equal(t, kindWrite, c.Kind)
					assert.Equal(t, tc.method, c.Method)
				}
			}
			assert.True(t, found, "path must be registered in APICommands")

			cmd, _, err := root.Find(strings.Fields(tc.path))
			require.NoError(t, err)
			assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
			assert.Empty(t, cmd.Annotations[annReadOnly])
			assert.Empty(t, cmd.Annotations[annDestructive])
		})
	}
}
