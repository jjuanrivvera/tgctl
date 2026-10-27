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

func TestRichMessage_Send_WireParams(t *testing.T) {
	t.Run("minimal rich message send", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendRichMessage") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok": true,
					"result": map[string]any{
						"message_id": 101,
						"chat":       map[string]any{"id": 12345678},
						"date":       1700000000,
						"rich_message": map[string]any{
							"blocks": []any{
								map[string]any{"type": "paragraph", "text": "Hello world"},
							},
						},
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "message", "send-rich",
			"--chat", "12345678",
			"--rich-message", `{"blocks":[{"type":"paragraph","text":"Hello world"}]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "101")
		assert.Equal(t, "12345678", got["chat_id"])
		rm, ok := got["rich_message"].(map[string]any)
		require.True(t, ok)
		blocks, ok := rm["blocks"].([]any)
		require.True(t, ok)
		assert.Len(t, blocks, 1)
	})

	t.Run("rich message send with all optional flags", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendRichMessage") {
				_ = json.NewDecoder(r.Body).Decode(&got)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok": true,
					"result": map[string]any{
						"message_id": 102,
						"chat":       map[string]any{"id": 12345678},
						"date":       1700000000,
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "message", "send-rich",
			"--chat", "12345678",
			"--rich-message", `{"blocks":[{"type":"section_heading","text":"Title"}]}`,
			"--business-connection-id", "bc_999",
			"--thread", "42",
			"--direct-messages-topic-id", "77",
			"--silent",
			"--protect-content",
			"--allow-paid-broadcast",
			"--message-effect-id", "fx_123",
			"--suggested-post-parameters", `{"can_send_now":true}`,
			"--reply-parameters", `{"message_id":100}`,
			"--reply-markup", `{"inline_keyboard":[[{"text":"Click","url":"https://example.com"}]]}`,
			"--ephemeral-message-parameters", `{"time_to_live":60}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "102")
		assert.Equal(t, "12345678", got["chat_id"])
		assert.Equal(t, "bc_999", got["business_connection_id"])
		assert.EqualValues(t, 42, got["message_thread_id"])
		assert.EqualValues(t, 77, got["direct_messages_topic_id"])
		assert.Equal(t, true, got["disable_notification"])
		assert.Equal(t, true, got["protect_content"])
		assert.Equal(t, true, got["allow_paid_broadcast"])
		assert.Equal(t, "fx_123", got["message_effect_id"])
		assert.NotNil(t, got["suggested_post_parameters"])
		assert.NotNil(t, got["reply_parameters"])
		assert.NotNil(t, got["reply_markup"])
		assert.NotNil(t, got["ephemeral_message_parameters"])
	})

	t.Run("rich message send with multipart file upload and attach reference", func(t *testing.T) {
		dir := t.TempDir()
		picPath := filepath.Join(dir, "image.jpg")
		require.NoError(t, os.WriteFile(picPath, []byte("JPEG_CONTENT"), 0o600))

		var formChatID, formRichMsg string
		var uploadedBytes []byte
		var uploadedFilename string

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendRichMessage") {
				require.NoError(t, r.ParseMultipartForm(1<<20))
				formChatID = r.FormValue("chat_id")
				formRichMsg = r.FormValue("rich_message")

				f, hdr, err := r.FormFile("pic")
				require.NoError(t, err)
				defer func() { _ = f.Close() }()
				uploadedFilename = hdr.Filename
				uploadedBytes, _ = io.ReadAll(f)

				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok": true,
					"result": map[string]any{
						"message_id": 103,
						"chat":       map[string]any{"id": 12345678},
						"date":       1700000000,
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		out, _, err := run(t, srv, "message", "send-rich",
			"--chat", "12345678",
			"--rich-message", `{"blocks":[{"type":"photo","photo":"attach://pic"}]}`,
			"--file", "pic=@"+picPath,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "103")
		assert.Equal(t, "12345678", formChatID)
		assert.Contains(t, formRichMsg, "attach://pic")
		assert.Equal(t, "image.jpg", uploadedFilename)
		assert.Equal(t, "JPEG_CONTENT", string(uploadedBytes))
	})
}

func TestRichMessage_Send_Errors(t *testing.T) {
	dir := t.TempDir()
	picPath := filepath.Join(dir, "image.jpg")
	require.NoError(t, os.WriteFile(picPath, []byte("IMG"), 0o600))

	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing chat",
			args:    []string{"message", "send-rich", "--rich-message", `{"blocks":[]}`},
			wantErr: "required",
		},
		{
			name:    "missing rich-message",
			args:    []string{"message", "send-rich", "--chat", "123"},
			wantErr: "required",
		},
		{
			name:    "invalid json",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", "not-json"},
			wantErr: "must be valid JSON",
		},
		{
			name:    "not a json object",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", `["array"]`},
			wantErr: "--rich-message must be a JSON object",
		},
		{
			name:    "blocks not an array",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", `{"blocks":"string"}`},
			wantErr: "--rich-message blocks must be a JSON array",
		},
		{
			name:    "missing upload for attach reference",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", `{"blocks":[{"type":"photo","photo":"attach://missing"}]}`},
			wantErr: `missing upload for attach reference "missing" (pass --file missing=@<path>)`,
		},
		{
			name:    "upload part collides with parameter",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", `{"blocks":[]}`, "--file", "chat_id=@" + picPath},
			wantErr: `upload part "chat_id" collides with parameter`,
		},
		{
			name:    "duplicate upload part",
			args:    []string{"message", "send-rich", "--chat", "123", "--rich-message", `{"blocks":[]}`, "--file", "pic=@" + picPath, "--file", "pic=@" + picPath},
			wantErr: `duplicate`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestRichMessage_Draft_WireParams(t *testing.T) {
	t.Run("minimal rich message draft", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendRichMessageDraft") {
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

		out, _, err := run(t, srv, "message", "draft-rich",
			"--chat", "12345",
			"--draft-id", "1",
			"--rich-message", `{"blocks":[{"type":"thinking","text":"Thinking..."}]}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.EqualValues(t, 12345, got["chat_id"])
		assert.EqualValues(t, 1, got["draft_id"])
		rm, ok := got["rich_message"].(map[string]any)
		require.True(t, ok)
		assert.NotNil(t, rm["blocks"])
	})

	t.Run("rich message draft with can_stop and keep_on_stop", func(t *testing.T) {
		var got map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/sendRichMessageDraft") {
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

		out, _, err := run(t, srv, "message", "draft-rich",
			"--chat", "12345",
			"--draft-id", "5",
			"--rich-message", `{"blocks":[{"type":"paragraph","text":"Progress..."}]}`,
			"--message-thread-id", "99",
			"--can-stop",
			"--keep-on-stop",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.EqualValues(t, 12345, got["chat_id"])
		assert.EqualValues(t, 5, got["draft_id"])
		assert.EqualValues(t, 99, got["message_thread_id"])
		assert.Equal(t, true, got["can_stop"])
		assert.Equal(t, true, got["keep_on_stop"])
	})
}

func TestRichMessage_Draft_Errors(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "zero draft id",
			args:    []string{"message", "draft-rich", "--chat", "123", "--draft-id", "0", "--rich-message", `{"blocks":[]}`},
			wantErr: "--draft-id must be non-zero",
		},
		{
			name:    "missing draft id",
			args:    []string{"message", "draft-rich", "--chat", "123", "--rich-message", `{"blocks":[]}`},
			wantErr: "required",
		},
		{
			name:    "draft rejects attach reference",
			args:    []string{"message", "draft-rich", "--chat", "123", "--draft-id", "1", "--rich-message", `{"blocks":[{"type":"photo","photo":"attach://pic"}]}`},
			wantErr: "drafts cannot reference uploaded files with attach://",
		},
		{
			name:    "draft rejects invalid json",
			args:    []string{"message", "draft-rich", "--chat", "123", "--draft-id", "1", "--rich-message", "bad-json"},
			wantErr: "must be valid JSON",
		},
		{
			name:    "draft rejects non-object",
			args:    []string{"message", "draft-rich", "--chat", "123", "--draft-id", "1", "--rich-message", `[1,2,3]`},
			wantErr: "--rich-message must be a JSON object",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestRichMessage_KindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	for _, tc := range []struct {
		path   string
		method string
	}{
		{"message send-rich", "sendRichMessage"},
		{"message draft-rich", "sendRichMessageDraft"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			var found bool
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

func TestRichMessage_RecorderAndDryRun(t *testing.T) {
	t.Run("recorder classification", func(t *testing.T) {
		assert.Equal(t, "rich", messageBearingMethods["sendRichMessage"])
		_, draftRecorded := messageBearingMethods["sendRichMessageDraft"]
		assert.False(t, draftRecorded, "drafts must never be recorded in SQLite message store")
	})

	t.Run("dry run send-rich", func(t *testing.T) {
		var hits int
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
		defer srv.Close()

		_, stderr, err := run(t, srv, "message", "send-rich",
			"--chat", "12345",
			"--rich-message", `{"blocks":[{"type":"paragraph","text":"test"}]}`,
			"--dry-run",
		)
		require.NoError(t, err)
		assert.Equal(t, 0, hits)
		assert.Contains(t, stderr, "sendRichMessage")
	})

	t.Run("dry run draft-rich", func(t *testing.T) {
		var hits int
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
		defer srv.Close()

		_, stderr, err := run(t, srv, "message", "draft-rich",
			"--chat", "12345",
			"--draft-id", "1",
			"--rich-message", `{"blocks":[]}`,
			"--dry-run",
		)
		require.NoError(t, err)
		assert.Equal(t, 0, hits)
		assert.Contains(t, stderr, "sendRichMessageDraft")
	})
}
