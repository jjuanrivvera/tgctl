package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateUploadPath checks a user-supplied upload path (a `--photo`/`--document` flag the
// user typed explicitly). The user choosing an absolute path is legitimate, so we do NOT
// confine it — we only confirm it exists and is a REGULAR file, to fail fast with a clear
// message instead of streaming a directory or a missing file to the API.
//
// "Regular" is the load-bearing word: an upload is buffered before it is sent, so a character
// device (`-F photo=@/dev/zero`) would grow that buffer until the process dies, and a FIFO
// would block the command before it ever made a request. Neither is a file anyone means to
// upload.
func ValidateUploadPath(p string) error {
	if p == "" {
		return fmt.Errorf("empty file path")
	}
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("file not readable: %w", err)
	}
	return checkRegular(p, info)
}

// checkRegular names what the path actually is, so the message says why it was refused.
func checkRegular(p string, info os.FileInfo) error {
	if info.Mode().IsRegular() {
		return nil
	}
	switch {
	case info.IsDir():
		return fmt.Errorf("%s is a directory, not a file", p)
	case info.Mode()&os.ModeNamedPipe != 0:
		return fmt.Errorf("%s is a named pipe, not a file (an upload needs a file it can size and re-read)", p)
	case info.Mode()&os.ModeDevice != 0:
		return fmt.Errorf("%s is a device, not a file", p)
	case info.Mode()&os.ModeSocket != 0:
		return fmt.Errorf("%s is a socket, not a file", p)
	}
	return fmt.Errorf("%s is not a regular file (mode %s)", p, info.Mode())
}

// ConfineToBase resolves a path that originated from DATA (e.g. a path referenced inside a
// JSON record or config), confining it to baseDir. It rejects absolute paths and `..`
// escapes AND resolves symlinks before returning, so a crafted input file cannot turn a
// read into an arbitrary-local-file-read or exfiltration primitive (GOAL.md §1).
func ConfineToBase(baseDir, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("absolute paths are not allowed in data references: %q", rel)
	}
	base, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	// Resolve the base's own symlinks first (e.g. macOS /var → /private/var) so the
	// containment comparison below is against the real path, not an alias.
	if rb, err := filepath.EvalSymlinks(base); err == nil {
		base = rb
	}
	joined := filepath.Join(base, rel)
	// Resolve symlinks on the part that exists, then re-check containment.
	resolved := joined
	if r, err := filepath.EvalSymlinks(joined); err == nil {
		resolved = r
	}
	rabs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	if rabs != base && !strings.HasPrefix(rabs, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes base directory: %q", rel)
	}
	return rabs, nil
}
