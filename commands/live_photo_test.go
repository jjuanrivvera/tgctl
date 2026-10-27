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

func TestLivePhoto_BothLocalFiles(t *testing.T) {
	tmpDir := t.TempDir()
	photoPath := filepath.Join(tmpDir, "still.jpg")
	require.NoError(t, os.WriteFile(photoPath, []byte("fake-jpeg-bytes"), 0600))
	livePhotoPath := filepath.Join(tmpDir, "motion.mp4")
	require.NoError(t, os.WriteFile(livePhotoPath, []byte("fake-mp4-bytes"), 0600))

	var receivedPhoto, receivedLivePhoto string
	var chatID string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sendLivePhoto") {
			require.NoError(t, r.ParseMultipartForm(10<<20))
			chatID = r.FormValue("chat_id")

			photoFile, _, err := r.FormFile("photo")
			require.NoError(t, err)
			b1, _ := io.ReadAll(photoFile)
			receivedPhoto = string(b1)

			liveFile, _, err := r.FormFile("live_photo")
			require.NoError(t, err)
			b2, _ := io.ReadAll(liveFile)
			receivedLivePhoto = string(b2)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 50,
					"chat":       map[string]any{"id": 12345},
					"live_photo": map[string]any{"file_id": "lp_file_id"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "media", "live-photo",
		"--chat", "12345",
		"--photo", photoPath,
		"--live-photo", livePhotoPath,
	)
	require.NoError(t, err)
	assert.Contains(t, out, "50")
	assert.Equal(t, "12345", chatID)
	assert.Equal(t, "fake-jpeg-bytes", receivedPhoto)
	assert.Equal(t, "fake-mp4-bytes", receivedLivePhoto)
}

func TestLivePhoto_BothFileIDs(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sendLivePhoto") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 51,
					"chat":       map[string]any{"id": 12345},
					"live_photo": map[string]any{"file_id": "lp_fid_video"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "media", "live-photo",
		"--chat", "12345",
		"--photo", "fid_photo_123",
		"--live-photo", "fid_video_456",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "51")
	assert.Equal(t, "12345", got["chat_id"])
	assert.Equal(t, "fid_photo_123", got["photo"])
	assert.Equal(t, "fid_video_456", got["live_photo"])
}

func TestLivePhoto_MixedLocalAndFileID(t *testing.T) {
	tmpDir := t.TempDir()
	photoPath := filepath.Join(tmpDir, "still.jpg")
	require.NoError(t, os.WriteFile(photoPath, []byte("still-data"), 0600))

	var receivedPhoto, receivedLiveParam string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sendLivePhoto") {
			require.NoError(t, r.ParseMultipartForm(10<<20))
			receivedLiveParam = r.FormValue("live_photo")

			photoFile, _, err := r.FormFile("photo")
			require.NoError(t, err)
			b, _ := io.ReadAll(photoFile)
			receivedPhoto = string(b)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 52,
					"chat":       map[string]any{"id": 12345},
					"live_photo": map[string]any{"file_id": "lp_fid"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "media", "live-photo",
		"--chat", "12345",
		"--photo", photoPath,
		"--live-photo", "existing_file_id_video",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "52")
	assert.Equal(t, "still-data", receivedPhoto)
	assert.Equal(t, "existing_file_id_video", receivedLiveParam)
}

func TestLivePhoto_RejectsURLs(t *testing.T) {
	t.Run("photo URL rejected", func(t *testing.T) {
		_, _, err := run(t, nil, "media", "live-photo",
			"--chat", "12345",
			"--photo", "https://example.com/pic.jpg",
			"--live-photo", "fid_video",
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--photo does not support HTTP(S) URLs")
	})

	t.Run("live-photo URL rejected", func(t *testing.T) {
		_, _, err := run(t, nil, "media", "live-photo",
			"--chat", "12345",
			"--photo", "fid_photo",
			"--live-photo", "http://example.com/video.mp4",
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--live-photo does not support HTTP(S) URLs")
	})
}

func TestLivePhoto_FullParametersAndPreservation(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sendLivePhoto") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 53,
					"chat":       map[string]any{"id": 12345},
					"live_photo": map[string]any{"file_id": "lp_fid"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "media", "live-photo",
		"--chat", "12345",
		"--photo", "fid_p",
		"--live-photo", "fid_l",
		"--business-connection-id", "bc_123",
		"--thread", "42",
		"--caption", "Sunset in Hawaii",
		"--parse-mode", "MarkdownV2",
		"--caption-entities", `[{"type":"bold","offset":0,"length":6}]`,
		"--show-caption-above-media",
		"--has-spoiler",
		"--silent",
		"--protect-content",
		"--allow-paid-broadcast",
		"--message-effect-id", "eff_5",
		"--reply-parameters", `{"message_id":1}`,
		"--reply-markup", `{"inline_keyboard":[[{"text":"Button","url":"https://t.me"}]]}`,
	)
	require.NoError(t, err)
	assert.Contains(t, out, "53")
	assert.Equal(t, "12345", got["chat_id"])
	assert.Equal(t, "fid_p", got["photo"])
	assert.Equal(t, "fid_l", got["live_photo"])
	assert.Equal(t, "bc_123", got["business_connection_id"])
	assert.EqualValues(t, 42, got["message_thread_id"])
	assert.Equal(t, "Sunset in Hawaii", got["caption"])
	assert.Equal(t, "MarkdownV2", got["parse_mode"])
	assert.NotNil(t, got["caption_entities"])
	assert.Equal(t, true, got["show_caption_above_media"])
	assert.Equal(t, true, got["has_spoiler"])
	assert.Equal(t, true, got["disable_notification"])
	assert.Equal(t, true, got["protect_content"])
	assert.Equal(t, true, got["allow_paid_broadcast"])
	assert.Equal(t, "eff_5", got["message_effect_id"])
	assert.NotNil(t, got["reply_parameters"])
	assert.NotNil(t, got["reply_markup"])
}

func TestLivePhoto_ExplicitFalseAndOmittedFlags(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sendLivePhoto") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"message_id": 54,
					"chat":       map[string]any{"id": 12345},
					"live_photo": map[string]any{"file_id": "lp_fid"},
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "media", "live-photo",
		"--chat", "12345",
		"--photo", "fid_p",
		"--live-photo", "fid_l",
		"--has-spoiler=false",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "54")
	assert.Equal(t, false, got["has_spoiler"])
	// Unused optional flags must be omitted entirely
	assert.NotContains(t, got, "business_connection_id")
	assert.NotContains(t, got, "message_thread_id")
	assert.NotContains(t, got, "caption")
	assert.NotContains(t, got, "disable_notification")
	assert.NotContains(t, got, "protect_content")
	assert.NotContains(t, got, "allow_paid_broadcast")
}

func TestLivePhoto_Validation(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"missing chat", []string{"media", "live-photo", "--photo", "p", "--live-photo", "l"}},
		{"missing photo", []string{"media", "live-photo", "--chat", "123", "--live-photo", "l"}},
		{"missing live-photo", []string{"media", "live-photo", "--chat", "123", "--photo", "p"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
		})
	}
}

func TestLivePhoto_KindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range APICommands() {
		if c.Path == "media live-photo" {
			found = true
			assert.Equal(t, kindWrite, c.Kind)
			assert.Equal(t, "sendLivePhoto", c.Method)
		}
	}
	assert.True(t, found, "media live-photo must be in APICommands")

	cmd, _, err := root.Find(strings.Fields("media live-photo"))
	require.NoError(t, err)
	assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
	assert.Empty(t, cmd.Annotations[annReadOnly])
	assert.Empty(t, cmd.Annotations[annDestructive])
}
