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

func TestStory_Post_Photo(t *testing.T) {
	dir := t.TempDir()
	photoPath := filepath.Join(dir, "photo.jpg")
	require.NoError(t, os.WriteFile(photoPath, []byte("PHOTO_DATA"), 0o600))

	var formBCID string
	var formActivePeriod string
	var formCaption string
	var contentJSON map[string]any
	var uploadedBytes []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/postStory") {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			formBCID = r.FormValue("business_connection_id")
			formActivePeriod = r.FormValue("active_period")
			formCaption = r.FormValue("caption")
			require.NoError(t, json.Unmarshal([]byte(r.FormValue("content")), &contentJSON))

			f, _, err := r.FormFile("story_content")
			require.NoError(t, err)
			defer func() { _ = f.Close() }()
			uploadedBytes, _ = io.ReadAll(f)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"id":   1001,
					"chat": map[string]any{"id": 12345},
					"date": 1700000000,
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "story", "post",
		"--business-connection-id", "bc_test",
		"--photo", photoPath,
		"--active-period", "86400",
		"--caption", "Daily update",
		"--post-to-chat-page",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "1001")
	assert.Equal(t, "bc_test", formBCID)
	assert.Equal(t, "86400", formActivePeriod)
	assert.Equal(t, "Daily update", formCaption)
	assert.Equal(t, "photo", contentJSON["type"])
	assert.Equal(t, "attach://story_content", contentJSON["photo"])
	assert.Equal(t, "PHOTO_DATA", string(uploadedBytes))
}

func TestStory_Post_VideoWithCover(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "video.mp4")
	coverPath := filepath.Join(dir, "cover.jpg")
	require.NoError(t, os.WriteFile(videoPath, []byte("VIDEO_DATA"), 0o600))
	require.NoError(t, os.WriteFile(coverPath, []byte("COVER_DATA"), 0o600))

	var contentJSON map[string]any
	var videoBytes, coverBytes []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/postStory") {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			require.NoError(t, json.Unmarshal([]byte(r.FormValue("content")), &contentJSON))

			vf, _, err := r.FormFile("story_content")
			require.NoError(t, err)
			defer func() { _ = vf.Close() }()
			videoBytes, _ = io.ReadAll(vf)

			cf, _, err := r.FormFile("story_cover")
			require.NoError(t, err)
			defer func() { _ = cf.Close() }()
			coverBytes, _ = io.ReadAll(cf)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"id":   1002,
					"chat": map[string]any{"id": 12345},
					"date": 1700000000,
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "story", "post",
		"--business-connection-id", "bc_test",
		"--video", videoPath,
		"--cover", coverPath,
		"--active-period", "43200",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "1002")
	assert.Equal(t, "video", contentJSON["type"])
	assert.Equal(t, "attach://story_content", contentJSON["video"])
	assert.Equal(t, "attach://story_cover", contentJSON["cover"])
	assert.Equal(t, "VIDEO_DATA", string(videoBytes))
	assert.Equal(t, "COVER_DATA", string(coverBytes))
}

func TestStory_Post_WithContentJSONAndFile(t *testing.T) {
	dir := t.TempDir()
	customPath := filepath.Join(dir, "custom.jpg")
	require.NoError(t, os.WriteFile(customPath, []byte("CUSTOM_DATA"), 0o600))

	var customBytes []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/postStory") {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			f, _, err := r.FormFile("custom")
			require.NoError(t, err)
			defer func() { _ = f.Close() }()
			customBytes, _ = io.ReadAll(f)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"id":   1003,
					"chat": map[string]any{"id": 12345},
					"date": 1700000000,
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "story", "post",
		"--business-connection-id", "bc_test",
		"--content", `{"type":"photo","photo":"attach://custom"}`,
		"--file", "custom=@"+customPath,
		"--active-period", "21600",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "1003")
	assert.Equal(t, "CUSTOM_DATA", string(customBytes))
}

func TestStory_Post_Errors(t *testing.T) {
	dir := t.TempDir()
	photoPath := filepath.Join(dir, "photo.jpg")
	videoPath := filepath.Join(dir, "video.mp4")
	require.NoError(t, os.WriteFile(photoPath, []byte("P"), 0o600))
	require.NoError(t, os.WriteFile(videoPath, []byte("V"), 0o600))

	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing connection id",
			args:    []string{"story", "post", "--photo", photoPath, "--active-period", "86400"},
			wantErr: "required",
		},
		{
			name:    "missing active period",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--photo", photoPath},
			wantErr: "required",
		},
		{
			name:    "invalid active period",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--photo", photoPath, "--active-period", "1000"},
			wantErr: "--active-period must be 21600 (6h), 43200 (12h), 86400 (24h), or 172800 (48h)",
		},
		{
			name:    "missing content or photo/video",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--active-period", "86400"},
			wantErr: "story content is required",
		},
		{
			name:    "both photo and video",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--photo", photoPath, "--video", videoPath, "--active-period", "86400"},
			wantErr: "pass either --photo or --video, not both",
		},
		{
			name:    "photo is a URL",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--photo", "https://example.com/p.jpg", "--active-period", "86400"},
			wantErr: "Telegram accepts no URL or file_id for stories",
		},
		{
			name:    "content JSON has URL",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--content", `{"type":"photo","photo":"https://example.com/pic.jpg"}`, "--active-period", "86400"},
			wantErr: "URLs and file_ids are not supported for stories",
		},
		{
			name:    "content JSON missing attach upload",
			args:    []string{"story", "post", "--business-connection-id", "bc", "--content", `{"type":"photo","photo":"attach://missing"}`, "--active-period", "86400"},
			wantErr: `missing upload for attach reference "missing"`,
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

func TestStory_Edit(t *testing.T) {
	dir := t.TempDir()
	photoPath := filepath.Join(dir, "new.jpg")
	require.NoError(t, os.WriteFile(photoPath, []byte("NEW_DATA"), 0o600))

	var formBCID, formStoryID, formCaption string
	var contentJSON map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/editStory") {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			formBCID = r.FormValue("business_connection_id")
			formStoryID = r.FormValue("story_id")
			formCaption = r.FormValue("caption")
			if c := r.FormValue("content"); c != "" {
				_ = json.Unmarshal([]byte(c), &contentJSON)
			}

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"id":   42,
					"chat": map[string]any{"id": 12345},
					"date": 1700000000,
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "story", "edit",
		"--business-connection-id", "bc_test",
		"--story-id", "42",
		"--photo", photoPath,
		"--caption", "Edited caption",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "42")
	assert.Equal(t, "bc_test", formBCID)
	assert.Equal(t, "42", formStoryID)
	assert.Equal(t, "Edited caption", formCaption)
	assert.Equal(t, "photo", contentJSON["type"])
}

func TestStory_Repost(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/repostStory") {
			_ = json.NewDecoder(r.Body).Decode(&got)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"result": map[string]any{
					"id":   2001,
					"chat": map[string]any{"id": 12345},
					"date": 1700000000,
				},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	out, _, err := run(t, srv, "story", "repost",
		"--business-connection-id", "dest_bc",
		"--from-business-connection-id", "src_bc",
		"--story-id", "42",
		"--active-period", "172800",
		"--caption", "Reposted!",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "2001")
	assert.Equal(t, "dest_bc", got["business_connection_id"])
	assert.Equal(t, "src_bc", got["from_business_connection_id"])
	assert.EqualValues(t, 42, got["story_id"])
	assert.EqualValues(t, 172800, got["active_period"])
	assert.Equal(t, "Reposted!", got["caption"])
}

func TestStory_Delete(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/deleteStory") {
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

	out, _, err := run(t, srv, "story", "delete",
		"--business-connection-id", "bc_test",
		"--story-id", "42",
	)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
	assert.Equal(t, "bc_test", got["business_connection_id"])
	assert.EqualValues(t, 42, got["story_id"])
}

func TestStory_KindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	cases := []struct {
		path   string
		method string
		kind   cmdKind
	}{
		{"story post", "postStory", kindWrite},
		{"story edit", "editStory", kindWrite},
		{"story repost", "repostStory", kindWrite},
		{"story delete", "deleteStory", kindDestructive},
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

func TestStory_DryRun(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantMethod string
	}{
		{"post dry run", []string{"story", "post", "--business-connection-id", "bc", "--content", `{"type":"photo","photo":"attach://p"}`, "--active-period", "86400", "--file", "p=@" + filepath.Join(t.TempDir(), "f.jpg"), "--dry-run"}, "postStory"},
		{"edit dry run", []string{"story", "edit", "--business-connection-id", "bc", "--story-id", "42", "--caption", "new", "--dry-run"}, "editStory"},
		{"repost dry run", []string{"story", "repost", "--business-connection-id", "bc", "--from-business-connection-id", "src", "--story-id", "42", "--active-period", "86400", "--dry-run"}, "repostStory"},
		{"delete dry run", []string{"story", "delete", "--business-connection-id", "bc", "--story-id", "42", "--dry-run"}, "deleteStory"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// create dummy file for post dry run if needed
			for i, arg := range tc.args {
				if strings.HasPrefix(arg, "p=@") {
					path := strings.TrimPrefix(arg, "p=@")
					require.NoError(t, os.WriteFile(path, []byte("x"), 0o600))
					tc.args[i] = "p=@" + path
				}
			}

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

func TestStory_NotRecordedInStore(t *testing.T) {
	for _, m := range []string{"postStory", "editStory", "repostStory", "deleteStory"} {
		_, recorded := messageBearingMethods[m]
		assert.False(t, recorded, "%s must not be treated as chat Message in recorder", m)
	}
}
