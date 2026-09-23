package api

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateUploadPath(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o600))

	require.NoError(t, ValidateUploadPath(file))
	assert.Error(t, ValidateUploadPath(""))
	assert.Error(t, ValidateUploadPath(filepath.Join(dir, "missing")))
	assert.Error(t, ValidateUploadPath(dir), "a directory is not a file")
}

func TestConfineToBase(t *testing.T) {
	base := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(base, "ok.json"), []byte("{}"), 0o600))

	got, err := ConfineToBase(base, "ok.json")
	require.NoError(t, err)
	// The returned path is symlink-resolved (macOS /var → /private/var), so compare against
	// the resolved base rather than the raw temp path.
	resolvedBase, _ := filepath.EvalSymlinks(base)
	assert.Equal(t, filepath.Join(resolvedBase, "ok.json"), got)

	// Escapes must be rejected.
	_, err = ConfineToBase(base, "../escape")
	assert.Error(t, err)
	// An absolute path must be rejected. t.TempDir() yields a platform-appropriate absolute
	// path (a Unix-style "/etc/passwd" is not absolute on Windows, so don't hardcode it).
	_, err = ConfineToBase(base, t.TempDir())
	assert.Error(t, err, "absolute paths from data are rejected")

	// A symlink pointing outside the base must be rejected after resolution.
	outside := t.TempDir()
	secret := filepath.Join(outside, "secret")
	require.NoError(t, os.WriteFile(secret, []byte("s"), 0o600))
	link := filepath.Join(base, "link")
	if err := os.Symlink(secret, link); err == nil {
		_, err = ConfineToBase(base, "link")
		assert.Error(t, err, "a symlink escaping the base must be rejected")
	}
}

// An upload is buffered before it is sent, so a character device would grow that buffer
// without end and a FIFO would block the command. Neither is a file anyone means to upload.
func TestValidateUploadPath_RejectsNonRegularFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIFOs and device files are not a Windows concern")
	}
	dir := t.TempDir()
	fifo := filepath.Join(dir, "pipe")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))

	err := ValidateUploadPath(fifo)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "named pipe")

	if _, statErr := os.Stat("/dev/zero"); statErr == nil {
		err := ValidateUploadPath("/dev/zero")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "device")
	}
}
