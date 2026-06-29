package api

import (
	"os"
	"path/filepath"
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
	_, err = ConfineToBase(base, "/etc/passwd")
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
