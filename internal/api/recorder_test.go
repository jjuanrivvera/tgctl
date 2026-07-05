package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// spyRecorder captures every Record call for assertions, guarded by a mutex since Client
// methods have no documented single-goroutine requirement.
type spyRecorder struct {
	mu    sync.Mutex
	calls []recordedCall
}

type recordedCall struct {
	method string
	params map[string]any
	result json.RawMessage
}

func (s *spyRecorder) Record(_ context.Context, method string, params map[string]any, result json.RawMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, recordedCall{method: method, params: params, result: result})
}

func TestClient_Call_RecordsOnSuccess(t *testing.T) {
	rec := &spyRecorder{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		okJSON(w, `{"message_id":42,"chat":{"id":1}}`)
	})
	c.recorder = rec

	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.NoError(t, err)

	require.Len(t, rec.calls, 1)
	assert.Equal(t, "sendMessage", rec.calls[0].method)
	assert.Equal(t, "hi", rec.calls[0].params["text"])
	assert.Contains(t, string(rec.calls[0].result), "message_id")
}

func TestClient_Call_DoesNotRecordOnError(t *testing.T) {
	rec := &spyRecorder{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		errJSON(w, http.StatusBadRequest, 400, "bad request", "")
	})
	c.recorder = rec

	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.Error(t, err)
	assert.Empty(t, rec.calls, "a failed call must never be recorded")
}

func TestClient_Call_DoesNotRecordOnDryRun(t *testing.T) {
	rec := &spyRecorder{}
	var buf bytes.Buffer
	auth, err := NewBotTokenAuth("123456:TESTHASHVALUE")
	require.NoError(t, err)
	c := New(auth, WithDryRun(true), WithDryRunWriter(&buf), WithRecorder(rec))

	_, err = c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.NoError(t, err)
	assert.Empty(t, rec.calls, "dry-run performs no request and must not be recorded")
}

func TestClient_Call_NilRecorderIsNoop(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		okJSON(w, `{"message_id":1}`)
	})
	// No WithRecorder attached: this must simply not panic or otherwise misbehave.
	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.NoError(t, err)
}

func TestClient_Upload_RecordsOnSuccess(t *testing.T) {
	rec := &spyRecorder{}
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "pic.jpg")
	require.NoError(t, os.WriteFile(imgPath, []byte("JPEGDATA"), 0o600))

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		okJSON(w, `{"message_id":5,"photo":[{"file_id":"ABC"}]}`)
	})
	c.recorder = rec

	_, err := c.Upload(t.Context(), "sendPhoto",
		map[string]any{"chat_id": "1"},
		map[string]string{"photo": imgPath}, false)
	require.NoError(t, err)

	require.Len(t, rec.calls, 1)
	assert.Equal(t, "sendPhoto", rec.calls[0].method)
	assert.Contains(t, string(rec.calls[0].result), "file_id")
}

func TestClient_Upload_DoesNotRecordOnError(t *testing.T) {
	rec := &spyRecorder{}
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "pic.jpg")
	require.NoError(t, os.WriteFile(imgPath, []byte("JPEGDATA"), 0o600))

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		errJSON(w, http.StatusBadRequest, 400, "bad request", "")
	})
	c.recorder = rec

	_, err := c.Upload(t.Context(), "sendPhoto",
		map[string]any{"chat_id": "1"},
		map[string]string{"photo": imgPath}, false)
	require.Error(t, err)
	assert.Empty(t, rec.calls)
}

// closingSpyRecorder is a spyRecorder that also implements io.Closer, so tests can prove
// Client.Close() actually delegates to it rather than merely not panicking.
type closingSpyRecorder struct {
	spyRecorder
	closed bool
	err    error
}

func (c *closingSpyRecorder) Close() error {
	c.closed = true
	return c.err
}

// TestClient_Close_ClosesRecorder pins the fix for a Windows CI regression: clientFromCmd
// attaches a store-backed Recorder to every client, and nothing was closing it, so the SQLite
// file handle stayed open for the life of the process — harmless on Unix, but it blocked
// Windows from deleting/renaming the file (e.g. a test's t.TempDir() cleanup). Close() must
// reach the recorder's own Close so every clientFromCmd caller can defer it uniformly.
func TestClient_Close_ClosesRecorder(t *testing.T) {
	rec := &closingSpyRecorder{}
	auth, err := NewBotTokenAuth("123456:TESTHASHVALUE")
	require.NoError(t, err)
	c := New(auth, WithRecorder(rec))

	require.NoError(t, c.Close())
	assert.True(t, rec.closed, "Client.Close must close a Recorder that implements io.Closer")
}

func TestClient_Close_PropagatesRecorderCloseError(t *testing.T) {
	wantErr := assert.AnError
	rec := &closingSpyRecorder{err: wantErr}
	auth, err := NewBotTokenAuth("123456:TESTHASHVALUE")
	require.NoError(t, err)
	c := New(auth, WithRecorder(rec))

	assert.ErrorIs(t, c.Close(), wantErr)
}

func TestClient_Close_NoRecorderIsNoop(t *testing.T) {
	auth, err := NewBotTokenAuth("123456:TESTHASHVALUE")
	require.NoError(t, err)
	c := New(auth) // no WithRecorder at all
	assert.NoError(t, c.Close())
}

func TestClient_Close_NonCloserRecorderIsNoop(t *testing.T) {
	// spyRecorder implements api.Recorder but not io.Closer — Close must not panic or error.
	auth, err := NewBotTokenAuth("123456:TESTHASHVALUE")
	require.NoError(t, err)
	c := New(auth, WithRecorder(&spyRecorder{}))
	assert.NoError(t, c.Close())
}
