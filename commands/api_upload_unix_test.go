//go:build !windows

package commands

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -F photo=@/dev/zero (or a FIFO) must fail locally: the body is buffered before it is sent.
func TestAPI_MultipartUpload_RejectsNonRegularFile(t *testing.T) {
	dir := t.TempDir()
	fifo := filepath.Join(dir, "pipe")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))
	srv := newServer(t, routes{"setChatPhoto": `true`})

	_, _, err := run(t, srv, "api", "setChatPhoto", "-F", "photo=@"+fifo)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "named pipe")
}
