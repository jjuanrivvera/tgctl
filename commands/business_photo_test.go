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

type businessPhotoProbe struct {
	photoField string
	connID     string
	isPublic   string
	partBytes  string
	partName   string
}

func businessPhotoServer(t *testing.T, probe *businessPhotoProbe) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/setBusinessAccountProfilePhoto") {
			require.NoError(t, r.ParseMultipartForm(10<<20))
			probe.photoField = r.FormValue("photo")
			probe.connID = r.FormValue("business_connection_id")
			probe.isPublic = r.FormValue("is_public")
			for name := range r.MultipartForm.File {
				probe.partName = name
				f, _, err := r.FormFile(name)
				require.NoError(t, err)
				b, _ := io.ReadAll(f)
				_ = f.Close()
				probe.partBytes = string(b)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestBusinessSetPhoto_Static(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "business_logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("JPEG_BIZ"), 0600))

	var probe businessPhotoProbe
	srv := businessPhotoServer(t, &probe)

	out, _, err := run(t, srv, "business", "set-photo",
		"--business-connection-id", "bc_999",
		"--photo", pic,
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")

	assert.Equal(t, "bc_999", probe.connID)
	var photo map[string]any
	require.NoError(t, json.Unmarshal([]byte(probe.photoField), &photo))
	assert.Equal(t, "static", photo["type"])
	assert.Equal(t, "attach://profile_photo", photo["photo"])
	assert.Equal(t, "profile_photo", probe.partName, "the uploaded part must not collide with the photo param name")
	assert.Equal(t, "JPEG_BIZ", probe.partBytes)
}

func TestBusinessSetPhoto_AnimatedWithPublic(t *testing.T) {
	dir := t.TempDir()
	clip := filepath.Join(dir, "business_intro.mp4")
	require.NoError(t, os.WriteFile(clip, []byte("MP4_BIZ"), 0600))

	var probe businessPhotoProbe
	srv := businessPhotoServer(t, &probe)

	out, _, err := run(t, srv, "business", "set-photo",
		"--business-connection-id", "bc_999",
		"--animated", clip,
		"--main-frame-timestamp", "2.5",
		"--public",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")

	assert.Equal(t, "bc_999", probe.connID)
	assert.Equal(t, "true", probe.isPublic)
	var photo map[string]any
	require.NoError(t, json.Unmarshal([]byte(probe.photoField), &photo))
	assert.Equal(t, "animated", photo["type"])
	assert.Equal(t, "attach://profile_photo", photo["animation"])
	assert.InDelta(t, 2.5, photo["main_frame_timestamp"], 0.0001)
	assert.Equal(t, "MP4_BIZ", probe.partBytes)
}

func TestBusinessSetPhoto_Rejections(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	clip := filepath.Join(dir, "intro.mp4")
	require.NoError(t, os.WriteFile(pic, []byte("J"), 0600))
	require.NoError(t, os.WriteFile(clip, []byte("M"), 0600))

	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing business connection id",
			args:    []string{"business", "set-photo", "--photo", pic},
			wantErr: "required flag(s) \"business-connection-id\"",
		},
		{
			name:    "neither photo nor animated passed",
			args:    []string{"business", "set-photo", "--business-connection-id", "bc"},
			wantErr: "one of --photo <file.jpg> or --animated <file.mp4> is required",
		},
		{
			name:    "both photo and animated passed",
			args:    []string{"business", "set-photo", "--business-connection-id", "bc", "--photo", pic, "--animated", clip},
			wantErr: "pass either --photo or --animated, not both",
		},
		{
			name:    "main-frame-timestamp with static photo",
			args:    []string{"business", "set-photo", "--business-connection-id", "bc", "--photo", pic, "--main-frame-timestamp", "1.0"},
			wantErr: "--main-frame-timestamp applies to --animated, not to a static --photo",
		},
		{
			name:    "photo is a URL",
			args:    []string{"business", "set-photo", "--business-connection-id", "bc", "--photo", "https://example.com/logo.jpg"},
			wantErr: "a profile photo must be an uploadable local file — Telegram accepts no URL or file_id here",
		},
		{
			name:    "photo is a non-existent file / file_id",
			args:    []string{"business", "set-photo", "--business-connection-id", "bc", "--photo", "AgACAgIAAxkBA..."},
			wantErr: "a profile photo must be an uploadable local file — Telegram accepts no URL or file_id here",
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

func TestBusinessSetPhoto_KindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range APICommands() {
		if c.Path == "business set-photo" {
			found = true
			assert.Equal(t, kindWrite, c.Kind)
			assert.Equal(t, "setBusinessAccountProfilePhoto", c.Method)
		}
	}
	assert.True(t, found, "business set-photo must be in APICommands")

	cmd, _, err := root.Find(strings.Fields("business set-photo"))
	require.NoError(t, err)
	assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
	assert.Empty(t, cmd.Annotations[annReadOnly])
	assert.Empty(t, cmd.Annotations[annDestructive])
}
