package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertUploadFile(t *testing.T, r *http.Request, field, filename string, want []byte) {
	t.Helper()
	f, hdr, err := r.FormFile(field)
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = f.Close() }()
	assert.Equal(t, filename, hdr.Filename)
	got, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestClient_Upload_Retry(t *testing.T) {
	for _, tt := range []struct {
		name       string
		status     int
		idempotent bool
		attempts   int32
	}{
		{"write_429", http.StatusTooManyRequests, false, 2},
		{"idempotent_5xx", http.StatusInternalServerError, true, 2},
		{"write_5xx", http.StatusInternalServerError, false, 1},
		{"idempotent_network", 0, true, 2},
		{"write_network", 0, false, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "video.bin")
			data := bytes.Repeat([]byte{0, 1, 0xff, '\r', '\n'}, 40000)
			require.NoError(t, os.WriteFile(path, data, 0o600))
			var calls atomic.Int32
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				attempt := calls.Add(1)
				if !assert.NoError(t, r.ParseMultipartForm(1<<20)) {
					return
				}
				assert.Equal(t, "123", r.FormValue("chat_id"))
				assert.Equal(t, "retry me", r.FormValue("caption"))
				assertUploadFile(t, r, "video", "video.bin", data)
				if attempt == 1 {
					if tt.status == 0 {
						conn, _, err := w.(http.Hijacker).Hijack()
						if assert.NoError(t, err) {
							_ = conn.Close()
						}
						return
					}
					errJSON(w, tt.status, tt.status, "try again", "")
					return
				}
				okJSON(w, `{"message_id":7}`)
			})
			raw, err := c.Upload(t.Context(), "sendVideo",
				map[string]any{"chat_id": "123", "caption": "retry me"},
				map[string]string{"video": path}, tt.idempotent)
			if tt.attempts == 1 {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, `{"message_id":7}`, string(raw))
			}
			assert.Equal(t, tt.attempts, calls.Load())
		})
	}
}

func TestClient_Upload_MissingFile(t *testing.T) {
	var calls atomic.Int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		okJSON(w, `true`)
	})
	path := filepath.Join(t.TempDir(), "missing.bin")
	done := make(chan error, 1)
	go func() {
		_, err := c.Upload(t.Context(), "sendDocument", map[string]any{"chat_id": "1"},
			map[string]string{"document": path}, false)
		done <- err
	}()
	select {
	case err := <-done:
		require.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.ErrorContains(t, err, "open document")
		assert.ErrorContains(t, err, path)
	case <-time.After(5 * time.Second):
		t.Fatal("upload with a missing file hung")
	}
	assert.Zero(t, calls.Load())
}

type uploadTransportFunc func(*http.Request) (*http.Response, error)

func (f uploadTransportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestClient_Upload_ClosesWriterOnEarlyReturn(t *testing.T) {
	for _, outcome := range []string{"network_error", "success", "cancel"} {
		t.Run(outcome, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "document.bin")
			require.NoError(t, os.WriteFile(path, []byte("payload"), 0o600))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			auth, err := NewBotTokenAuth("123456:AAH-fake-test-credential-value")
			require.NoError(t, err)
			var body *multipartBody
			transportErr := errors.New("connection failed before reading body")
			c := New(auth, WithRPS(0), WithHTTPClient(&http.Client{
				Transport: uploadTransportFunc(func(r *http.Request) (*http.Response, error) {
					body = r.Body.(*multipartBody)
					switch outcome {
					case "network_error":
						return nil, transportErr
					case "cancel":
						cancel()
						return nil, r.Context().Err()
					default:
						return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header),
							Body: io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`))}, nil
					}
				}),
			}))
			done := make(chan error, 1)
			go func() {
				_, err := c.Upload(ctx, "sendDocument", map[string]any{"chat_id": "1"},
					map[string]string{"document": path}, false)
				done <- err
			}()
			select {
			case err := <-done:
				switch outcome {
				case "network_error":
					require.ErrorIs(t, err, transportErr)
				case "cancel":
					require.ErrorIs(t, err, context.Canceled)
				default:
					require.NoError(t, err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("upload did not return after the transport stopped reading")
			}
			require.NotNil(t, body)
			select {
			case <-body.done:
			default:
				t.Fatal("multipart writer still running after Upload returned")
			}
		})
	}
}

func TestMultipartBody_WriterError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.bin")
	body, _ := newMultipartBody(map[string]any{"chat_id": "1"}, map[string]string{"document": path})
	defer func() { _ = body.Close() }()
	_, err := io.Copy(io.Discard, body)
	require.ErrorIs(t, err, os.ErrNotExist)
	assert.ErrorContains(t, err, "open document")
}

func TestClient_Upload_CancelDuringCopy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "video.bin")
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, f.Truncate(32<<20))
	require.NoError(t, f.Close())
	started := make(chan struct{})
	release := make(chan struct{})
	c := newTestClient(t, func(_ http.ResponseWriter, r *http.Request) {
		mr, err := r.MultipartReader()
		if !assert.NoError(t, err) {
			return
		}
		part, err := mr.NextPart()
		if !assert.NoError(t, err) {
			return
		}
		_, err = io.CopyN(io.Discard, part, 1024)
		assert.NoError(t, err)
		close(started)
		<-release
	})
	defer close(release)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	t.Cleanup(transport.CloseIdleConnections)
	bodies := make(chan *multipartBody, 1)
	c.http = &http.Client{Transport: uploadTransportFunc(func(r *http.Request) (*http.Response, error) {
		bodies <- r.Body.(*multipartBody)
		return transport.RoundTrip(r)
	})}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := c.Upload(ctx, "sendVideo", nil, map[string]string{"video": path}, false)
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("server did not start receiving the upload")
	}
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("canceled upload hung while the server stopped reading")
	}
	body := <-bodies
	select {
	case <-body.done:
	default:
		t.Fatal("multipart writer still running after cancellation")
	}
}

type failingCloseBody struct {
	io.Reader
	err error
}

func (b failingCloseBody) Close() error { return b.err }

func TestMultipartBody_CloseReportsWriterError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gone.bin")
	body, _ := newMultipartBody(nil, map[string]string{"document": path})
	<-body.(*multipartBody).done // the open fails before anything is written
	err := body.Close()
	require.ErrorIs(t, err, os.ErrNotExist)
}

// A transport may answer without reading the body; an incomplete body must still fail.
func TestClient_Do_IncompleteBodyIsNotSuccess(t *testing.T) {
	auth, err := NewBotTokenAuth("123456:AAH-fake-test-credential-value")
	require.NoError(t, err)
	c := New(auth, WithRPS(0), WithHTTPClient(&http.Client{
		Transport: uploadTransportFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header),
				Body: io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`))}, nil
		}),
	}))
	writerErr := errors.New("open document: file vanished")
	_, err = c.do(t.Context(), &preparedRequest{
		method: "sendDocument",
		newBody: func() (io.ReadCloser, string) {
			return failingCloseBody{Reader: strings.NewReader("partial"), err: writerErr}, "multipart/form-data"
		},
	})
	require.ErrorIs(t, err, writerErr)
}
