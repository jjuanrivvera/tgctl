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

func TestWebAppRequests(t *testing.T) {
	cases := []struct {
		method, result string
		args           []string
		expected       map[string]any
	}{
		{"answerWebAppQuery", `{"inline_message_id":"inline-1"}`, []string{"inline", "answer-web-app", "--query-id", "q1", "--result", `{"type":"article","id":"a"}`}, map[string]any{"web_app_query_id": "q1", "result": map[string]any{"type": "article", "id": "a"}}},
		{"savePreparedInlineMessage", `{"id":"prepared","expiration_date":12345}`, []string{"inline", "prepare-message", "--user", "12345", "--result", `{"type":"article","id":"a"}`, "--allow-user-chats=false", "--allow-group-chats"}, map[string]any{"user_id": float64(12345), "result": map[string]any{"type": "article", "id": "a"}, "allow_user_chats": false, "allow_group_chats": true}},
		{"savePreparedKeyboardButton", `{"id":"button-id"}`, []string{"inline", "prepare-button", "--user", "12345", "--button", `{"text":"Choose","request_chat":{"request_id":1}}`}, map[string]any{"user_id": float64(12345), "button": map[string]any{"text": "Choose", "request_chat": map[string]any{"request_id": float64(1)}}}},
		{"sendChatJoinRequestWebApp", `true`, []string{"member", "join-web-app", "--query-id", "join-q", "--web-app-url", "https://example.org/app"}, map[string]any{"chat_join_request_query_id": "join-q", "web_app_url": "https://example.org/app"}},
	}
	for _, tc := range cases {
		t.Run(tc.method, func(t *testing.T) {
			var got map[string]any
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/bot123456:TESTHASHVALUE/"+tc.method, r.URL.Path)
				require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
				_, _ = w.Write([]byte(`{"ok":true,"result":` + tc.result + `}`))
			}))
			defer srv.Close()
			out, _, err := run(t, srv, tc.args...)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
			if tc.method == "savePreparedInlineMessage" {
				assert.Contains(t, out, "prepared")
				assert.Contains(t, out, "12345")
			}
			if tc.method == "savePreparedKeyboardButton" {
				assert.Contains(t, out, "button-id")
			}
			if tc.method == "answerWebAppQuery" {
				assert.Contains(t, out, "inline-1")
			}
		})
	}
}

func TestWebAppValidationAndKinds(t *testing.T) {
	for _, args := range [][]string{
		{"inline", "answer-web-app", "--query-id", "q", "--result", "[]"},
		{"inline", "prepare-message", "--user", "1", "--result", "null"},
		{"inline", "prepare-button", "--user", "1", "--button", "[]"},
		{"inline", "prepare-button", "--user", "1", "--button", `{"text":"x"}`},
		{"inline", "prepare-button", "--user", "1", "--button", `{"request_chat":{},"request_users":{}}`},
		{"inline", "prepare-button", "--user", "1", "--button", `{"request_chat":null}`},
		{"inline", "prepare-button", "--user", "1", "--button", `{"request_chat":{},"web_app":{}}`},
	} {
		_, _, err := run(t, nil, args...)
		require.Error(t, err)
	}
	root := NewRootCmd()
	for _, path := range []string{"inline answer-web-app", "inline prepare-message", "inline prepare-button", "member join-web-app"} {
		found := false
		for _, c := range APICommands() {
			if c.Path == path {
				found = true
				assert.Equal(t, kindWrite, c.Kind)
			}
		}
		assert.True(t, found, path)
		cmd, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
		assert.Empty(t, cmd.Annotations[annReadOnly])
		assert.Empty(t, cmd.Annotations[annDestructive])
	}
}
