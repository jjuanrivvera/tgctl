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

// profilePhotoProbe captures what setMyProfilePhoto actually received: the InputProfilePhoto
// object and the uploaded part it points at.
type profilePhotoProbe struct {
	photoField string
	partBytes  string
	partName   string
}

func profilePhotoServer(t *testing.T, probe *profilePhotoProbe) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/setMyProfilePhoto") {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			probe.photoField = r.FormValue("photo")
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

// The static case: the file goes up as its own part and the photo parameter is the
// InputProfilePhotoStatic object naming it (issue #23).
func TestBotSetPhoto_Static(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("JPEGDATA"), 0o600))

	var probe profilePhotoProbe
	srv := profilePhotoServer(t, &probe)

	out, _, err := run(t, srv, "bot", "set-photo", "--photo", pic)
	require.NoError(t, err)
	assert.Contains(t, out, "true")

	var photo map[string]any
	require.NoError(t, json.Unmarshal([]byte(probe.photoField), &photo))
	assert.Equal(t, "static", photo["type"])
	assert.Equal(t, "attach://profile_photo", photo["photo"])
	assert.Equal(t, "profile_photo", probe.partName, "the part must not be named photo — that is the parameter")
	assert.Equal(t, "JPEGDATA", probe.partBytes)
}

// The animated case uses "animation", not "photo", and carries the optional frame timestamp.
func TestBotSetPhoto_Animated(t *testing.T) {
	dir := t.TempDir()
	clip := filepath.Join(dir, "intro.mp4")
	require.NoError(t, os.WriteFile(clip, []byte("MP4DATA"), 0o600))

	var probe profilePhotoProbe
	srv := profilePhotoServer(t, &probe)

	_, _, err := run(t, srv, "bot", "set-photo", "--animated", clip, "--main-frame-timestamp", "1.5")
	require.NoError(t, err)

	var photo map[string]any
	require.NoError(t, json.Unmarshal([]byte(probe.photoField), &photo))
	assert.Equal(t, "animated", photo["type"])
	assert.Equal(t, "attach://profile_photo", photo["animation"])
	assert.InDelta(t, 1.5, photo["main_frame_timestamp"], 0.0001)
	assert.Equal(t, "MP4DATA", probe.partBytes)
}

func TestBotSetPhoto_Rejections(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	clip := filepath.Join(dir, "intro.mp4")
	require.NoError(t, os.WriteFile(pic, []byte("J"), 0o600))
	require.NoError(t, os.WriteFile(clip, []byte("M"), 0o600))
	srv := newServer(t, routes{"setMyProfilePhoto": `true`})

	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{{
		name:    "neither flag",
		args:    []string{"bot", "set-photo"},
		wantErr: "one of --photo",
	}, {
		name:    "both flags",
		args:    []string{"bot", "set-photo", "--photo", pic, "--animated", clip},
		wantErr: "not both",
	}, {
		// A URL is fine for sendPhoto and impossible here: a profile photo can't be reused.
		name:    "url instead of a file",
		args:    []string{"bot", "set-photo", "--photo", "https://example.com/logo.jpg"},
		wantErr: "uploadable local file",
	}, {
		name:    "file_id instead of a file",
		args:    []string{"bot", "set-photo", "--photo", "AgACAgEAAxkBAAI"},
		wantErr: "uploadable local file",
	}, {
		name:    "timestamp without an animation",
		args:    []string{"bot", "set-photo", "--photo", pic, "--main-frame-timestamp", "2"},
		wantErr: "applies to --animated",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, srv, tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestBotSetPhoto_DryRunRedactsToken(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("JPEGDATA"), 0o600))

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
	t.Cleanup(srv.Close)

	_, stderr, err := run(t, srv, "bot", "set-photo", "--photo", pic, "--dry-run")
	require.NoError(t, err)
	assert.Equal(t, 0, hits)
	assert.Contains(t, stderr, "-F profile_photo=@")
	assert.Contains(t, stderr, "<redacted>")
	assert.NotContains(t, stderr, "TESTHASHVALUE")
}

func TestBotRemovePhoto(t *testing.T) {
	srv := newServer(t, routes{"removeMyProfilePhoto": `true`})
	out, _, err := run(t, srv, "bot", "remove-photo")
	require.NoError(t, err)
	assert.Contains(t, out, "true")
}

// bot photo asks getMe for the bot's own id, then getUserProfilePhotos for it.
func TestBotPhoto_LooksUpOwnID(t *testing.T) {
	var gotUserID string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/getMe"):
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":123456,"is_bot":true,"username":"testbot"}}`))
		case strings.HasSuffix(r.URL.Path, "/getUserProfilePhotos"):
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			gotUserID, _ = body["user_id"].(string)
			_, _ = w.Write([]byte(`{"ok":true,"result":{"total_count":1,"photos":[[]]}}`))
		default:
			t.Errorf("unexpected call to %s", r.URL.Path)
		}
	}))
	t.Cleanup(srv.Close)

	out, _, err := run(t, srv, "bot", "photo")
	require.NoError(t, err)
	assert.Equal(t, "123456", gotUserID, "the bot's own id must come from getMe")
	assert.Contains(t, out, "1", "total_count is rendered")
}

func TestBotPhoto_JSON(t *testing.T) {
	srv := newServer(t, routes{
		"getMe":                `{"id":7,"is_bot":true,"username":"b"}`,
		"getUserProfilePhotos": `{"total_count":0,"photos":[]}`,
	})
	out, _, err := run(t, srv, "bot", "photo", "-o", "json")
	require.NoError(t, err)
	assert.Contains(t, out, `"total_count": 0`)
}
