package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEphemeral_Delete(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/deleteEphemeralMessage") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok":     true,
				"result": true,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "ephemeral", "delete",
		"--chat", "-1001234567890",
		"--receiver", "12345",
		"--ephemeral-message-id", "42",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
	assert.Equal(t, "-1001234567890", got["chat_id"])
	assert.EqualValues(t, 12345, got["receiver_user_id"])
	assert.EqualValues(t, 42, got["ephemeral_message_id"])
}

func TestEphemeral_EditText(t *testing.T) {
	t.Run("with plain text", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageText") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":     true,
					"result": true,
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "ephemeral", "edit-text",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--text", "Updated text",
			"--parse-mode", "HTML",
			"--reply-markup", `{"inline_keyboard":[[{"text":"OK","callback_data":"ok"}]]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-1001234567890", got["chat_id"])
		assert.EqualValues(t, 12345, got["receiver_user_id"])
		assert.EqualValues(t, 42, got["ephemeral_message_id"])
		assert.Equal(t, "Updated text", got["text"])
		assert.Equal(t, "HTML", got["parse_mode"])
		assert.NotNil(t, got["reply_markup"])
	})

	t.Run("with rich message", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageText") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":     true,
					"result": true,
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "ephemeral", "edit-text",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--rich-message", `{"blocks":[{"type":"paragraph","text":"New block"}]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		rm, ok := got["rich_message"].(map[string]any)
		require.True(t, ok)
		assert.NotNil(t, rm["blocks"])
	})

	t.Run("error when neither text nor rich-message provided", func(t *testing.T) {
		_, _, err := run(t, nil, "ephemeral", "edit-text",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pass either --text or --rich-message")
	})

	t.Run("error when rich-message is not an object", func(t *testing.T) {
		_, _, err := run(t, nil, "ephemeral", "edit-text",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--rich-message", `["array"]`,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--rich-message must be a JSON object")
	})
}

func TestEphemeral_EditCaption(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageCaption") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok":     true,
				"result": true,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "ephemeral", "edit-caption",
		"--chat", "-1001234567890",
		"--receiver", "12345",
		"--ephemeral-message-id", "42",
		"--caption", "Updated caption",
		"--show-caption-above-media",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
	assert.Equal(t, "-1001234567890", got["chat_id"])
	assert.EqualValues(t, 12345, got["receiver_user_id"])
	assert.EqualValues(t, 42, got["ephemeral_message_id"])
	assert.Equal(t, "Updated caption", got["caption"])
	assert.Equal(t, true, got["show_caption_above_media"])
}

func TestEphemeral_EditMedia(t *testing.T) {
	t.Run("with URL in media", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageMedia") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":     true,
					"result": true,
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "ephemeral", "edit-media",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--media", `{"type":"photo","media":"https://example.com/photo.jpg"}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-1001234567890", got["chat_id"])
		assert.EqualValues(t, 12345, got["receiver_user_id"])
		assert.EqualValues(t, 42, got["ephemeral_message_id"])
		m, ok := got["media"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "photo", m["type"])
	})

	t.Run("with file upload and attach reference", func(t *testing.T) {
		dir := t.TempDir()
		picPath := filepath.Join(dir, "new.jpg")
		require.NoError(t, os.WriteFile(picPath, []byte("NEW_PIC"), 0o600))

		var formChatID string
		var uploadedBytes []byte

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageMedia") {
				require.NoError(t, r.ParseMultipartForm(1<<20))
				formChatID = r.FormValue("chat_id")
				f, _, err := r.FormFile("pic")
				require.NoError(t, err)
				defer func() { _ = f.Close() }()
				uploadedBytes, _ = io.ReadAll(f)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":     true,
					"result": true,
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "ephemeral", "edit-media",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--media", `{"type":"photo","media":"attach://pic"}`,
			"--file", "pic=@"+picPath,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-1001234567890", formChatID)
		assert.Equal(t, "NEW_PIC", string(uploadedBytes))
	})

	t.Run("error on missing attach upload", func(t *testing.T) {
		_, _, err := run(t, nil, "ephemeral", "edit-media",
			"--chat", "-1001234567890",
			"--receiver", "12345",
			"--ephemeral-message-id", "42",
			"--media", `{"type":"photo","media":"attach://missing"}`,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `missing upload for attach reference "missing"`)
	})
}

func TestEphemeral_EditReplyMarkup(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/editEphemeralMessageReplyMarkup") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok":     true,
				"result": true,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "ephemeral", "edit-reply-markup",
		"--chat", "-1001234567890",
		"--receiver", "12345",
		"--ephemeral-message-id", "42",
		"--reply-markup", `{"inline_keyboard":[[{"text":"Refresh","callback_data":"ref"}]]}`,
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
	assert.Equal(t, "-1001234567890", got["chat_id"])
	assert.EqualValues(t, 12345, got["receiver_user_id"])
	assert.EqualValues(t, 42, got["ephemeral_message_id"])
	assert.NotNil(t, got["reply_markup"])
}

func TestEphemeral_RequiredTripleValidation(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"delete missing chat", []string{"ephemeral", "delete", "--receiver", "123", "--ephemeral-message-id", "1"}},
		{"delete missing receiver", []string{"ephemeral", "delete", "--chat", "123", "--ephemeral-message-id", "1"}},
		{"delete missing id", []string{"ephemeral", "delete", "--chat", "123", "--receiver", "123"}},
		{"edit-text missing chat", []string{"ephemeral", "edit-text", "--receiver", "123", "--ephemeral-message-id", "1", "--text", "hi"}},
		{"edit-text missing receiver", []string{"ephemeral", "edit-text", "--chat", "123", "--ephemeral-message-id", "1", "--text", "hi"}},
		{"edit-text missing id", []string{"ephemeral", "edit-text", "--chat", "123", "--receiver", "123", "--text", "hi"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "required")
		})
	}
}

func TestEphemeral_KindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	cases := []struct {
		path   string
		method string
		kind   cmdKind
	}{
		{"ephemeral delete", "deleteEphemeralMessage", kindDestructive},
		{"ephemeral edit-text", "editEphemeralMessageText", kindWrite},
		{"ephemeral edit-caption", "editEphemeralMessageCaption", kindWrite},
		{"ephemeral edit-media", "editEphemeralMessageMedia", kindWrite},
		{"ephemeral edit-reply-markup", "editEphemeralMessageReplyMarkup", kindWrite},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			var found bool
			for _, c := range APICommands() {
				if c.Path == tc.path {
					found = true
					assert.Equal(t, tc.kind, c.Kind)
					assert.Equal(t, tc.method, c.Method)
				}
			}
			assert.True(t, found, "path must be in APICommands")

			cmd, _, err := root.Find(strings.Fields(tc.path))
			require.NoError(t, err)
			assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
			if tc.kind == kindDestructive {
				assert.Equal(t, "true", cmd.Annotations[annDestructive])
			} else {
				assert.Empty(t, cmd.Annotations[annDestructive])
			}
			assert.Empty(t, cmd.Annotations[annReadOnly])
		})
	}
}

func TestEphemeral_DryRun(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantMethod string
	}{
		{"delete dry run", []string{"ephemeral", "delete", "--chat", "123", "--receiver", "456", "--ephemeral-message-id", "1", "--dry-run"}, "deleteEphemeralMessage"},
		{"edit-text dry run", []string{"ephemeral", "edit-text", "--chat", "123", "--receiver", "456", "--ephemeral-message-id", "1", "--text", "hi", "--dry-run"}, "editEphemeralMessageText"},
		{"edit-caption dry run", []string{"ephemeral", "edit-caption", "--chat", "123", "--receiver", "456", "--ephemeral-message-id", "1", "--caption", "cap", "--dry-run"}, "editEphemeralMessageCaption"},
		{"edit-media dry run", []string{"ephemeral", "edit-media", "--chat", "123", "--receiver", "456", "--ephemeral-message-id", "1", "--media", `{"type":"photo"}`, "--dry-run"}, "editEphemeralMessageMedia"},
		{"edit-reply-markup dry run", []string{"ephemeral", "edit-reply-markup", "--chat", "123", "--receiver", "456", "--ephemeral-message-id", "1", "--dry-run"}, "editEphemeralMessageReplyMarkup"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var hits int
			srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
			defer srv.Close()

			_, stderr, err := run(t, srv, tc.args...)
			require.NoError(t, err)
			assert.Equal(t, 0, hits)
			assert.Contains(t, stderr, tc.wantMethod)
		})
	}
}
