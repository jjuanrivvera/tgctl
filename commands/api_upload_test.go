package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The escape hatch has to reach the methods that refuse anything but an uploaded file
// (setMyProfilePhoto, setChatPhoto, uploadStickerFile, ...), which means multipart with a
// JSON field that points at the upload through attach:// (issue #24).
func TestAPI_MultipartUpload_SendsFileAndFields(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("JPEGDATA"), 0o600))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseMultipartForm(1<<20))

		// The JSON body rides along as a field, serialized, so it can name the part.
		var photo struct {
			Type  string `json:"type"`
			Photo string `json:"photo"`
		}
		require.NoError(t, json.Unmarshal([]byte(r.FormValue("photo")), &photo))
		assert.Equal(t, "static", photo.Type)
		assert.Equal(t, "attach://pic", photo.Photo)

		f, hdr, err := r.FormFile("pic")
		require.NoError(t, err)
		defer func() { _ = f.Close() }()
		assert.Equal(t, "logo.jpg", hdr.Filename)
		data, _ := io.ReadAll(f)
		assert.Equal(t, "JPEGDATA", string(data))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	t.Cleanup(srv.Close)

	out, _, err := run(t, srv, "api", "setMyProfilePhoto",
		"-d", `{"photo":{"type":"static","photo":"attach://pic"}}`, "-F", "pic=@"+pic)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
}

// -q fields and several parts in one call, the shape setChatPhoto needs.
func TestAPI_MultipartUpload_QueryFieldsAndSeveralParts(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.png")
	b := filepath.Join(dir, "b.png")
	require.NoError(t, os.WriteFile(a, []byte("A"), 0o600))
	require.NoError(t, os.WriteFile(b, []byte("B"), 0o600))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "@mygroup", r.FormValue("chat_id"))
		for part, want := range map[string]string{"photo": "A", "thumb": "B"} {
			f, _, err := r.FormFile(part)
			require.NoError(t, err, "part %s", part)
			data, _ := io.ReadAll(f)
			_ = f.Close()
			assert.Equal(t, want, string(data), "part %s", part)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	t.Cleanup(srv.Close)

	_, _, err := run(t, srv, "api", "setChatPhoto", "-q", "chat_id=@mygroup",
		"-F", "photo=@"+a, "-F", "thumb=@"+b)
	require.NoError(t, err)
}

// A dry run prints the curl with its -F parts and the token redacted, and sends nothing.
func TestAPI_MultipartUpload_DryRun(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("JPEGDATA"), 0o600))

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
	t.Cleanup(srv.Close)

	_, stderr, err := run(t, srv, "api", "setMyProfilePhoto",
		"-d", `{"photo":{"type":"static","photo":"attach://pic"}}`, "-F", "pic=@"+pic, "--dry-run")
	require.NoError(t, err)
	assert.Equal(t, 0, hits, "a dry run must not reach the API")
	assert.Contains(t, stderr, "-F pic=@")
	assert.Contains(t, stderr, "curl -sS -X POST")
	assert.Contains(t, stderr, "<redacted>")
	assert.NotContains(t, stderr, "TESTHASHVALUE")
}

func TestAPI_MultipartUpload_BadSpecs(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.jpg")
	require.NoError(t, os.WriteFile(pic, []byte("x"), 0o600))
	srv := newServer(t, routes{"setChatPhoto": `true`})

	cases := []struct {
		name, spec, wantErr string
	}{
		{"no @", "photo=cover.jpg", "must be @<path>"},
		{"no =", "photo@cover.jpg", "want name=@path"},
		{"empty name", "=@cover.jpg", "want name=@path"},
		{"empty path", "photo=@", "empty path after @"},
		{"missing file", "photo=@" + filepath.Join(dir, "nope.jpg"), "file not readable"},
		{"directory", "photo=@" + dir, "is a directory"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, srv, "api", "setChatPhoto", "-F", tc.spec)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}

	t.Run("duplicate part", func(t *testing.T) {
		_, _, err := run(t, srv, "api", "setChatPhoto", "-F", "photo=@"+pic, "-F", "photo=@"+pic)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})
}

// Without -F the call stays JSON, exactly as before.
func TestAPI_NoUploads_StaysJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.True(t, strings.Contains(string(body), "@me"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1}}`))
	}))
	t.Cleanup(srv.Close)

	_, _, err := run(t, srv, "api", "sendMessage", "-q", "chat_id=@me", "-q", "text=hi")
	require.NoError(t, err)
}

// -F photo=@/dev/zero (or a FIFO) must fail locally: the body is buffered before it is sent.
func TestAPI_MultipartUpload_RejectsNonRegularFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIFOs are not a Windows concern")
	}
	dir := t.TempDir()
	fifo := filepath.Join(dir, "pipe")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))
	srv := newServer(t, routes{"setChatPhoto": `true`})

	_, _, err := run(t, srv, "api", "setChatPhoto", "-F", "photo=@"+fifo)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "named pipe")
}
