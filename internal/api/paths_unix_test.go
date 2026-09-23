//go:build !windows

package api

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An upload is buffered before it is sent, so a character device would grow that buffer
// without end and a FIFO would block the command. Neither is a file anyone means to upload.
func TestValidateUploadPath_RejectsNonRegularFiles(t *testing.T) {
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
