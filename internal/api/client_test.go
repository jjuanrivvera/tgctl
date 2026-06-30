package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Call_Success(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/sendMessage"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "hi", body["text"])
		okJSON(w, `{"message_id":42,"text":"hi"}`)
	})

	raw, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.NoError(t, err)

	var got struct {
		MessageID Int    `json:"message_id"`
		Text      string `json:"text"`
	}
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(42), got.MessageID.Int64())
	assert.Equal(t, "hi", got.Text)
}

func TestClient_GetMe(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/getMe"))
		okJSON(w, `{"id":123456,"is_bot":true,"first_name":"Test","username":"testbot"}`)
	})
	me, err := c.GetMe(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "@testbot", me.DisplayName())
	assert.True(t, me.IsBot)
	assert.Equal(t, ID("123456"), me.ID)
}

func TestClient_APIError_CarriesHint(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		errJSON(w, http.StatusUnauthorized, 401, "Unauthorized", "")
	})
	_, err := c.Call(t.Context(), "getMe", nil, true)
	require.Error(t, err)
	var ae *APIError
	require.ErrorAs(t, err, &ae)
	assert.Equal(t, 401, ae.Code)
	assert.Contains(t, err.Error(), "auth login")
}

func TestClient_RetriesIdempotentOn5xx(t *testing.T) {
	var calls atomic.Int32
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) < 3 {
			errJSON(w, http.StatusInternalServerError, 500, "Internal Server Error", "")
			return
		}
		okJSON(w, `{"id":1,"is_bot":true,"first_name":"T"}`)
	})
	_, err := c.Call(t.Context(), "getMe", nil, true) // idempotent
	require.NoError(t, err)
	assert.EqualValues(t, 3, calls.Load())
}

func TestClient_DoesNotRetryWriteOn5xx(t *testing.T) {
	var calls atomic.Int32
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		errJSON(w, http.StatusInternalServerError, 500, "Internal Server Error", "")
	})
	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "x"}, false)
	require.Error(t, err)
	assert.EqualValues(t, 1, calls.Load(), "a write must not auto-retry on an ambiguous 5xx")
}

func TestClient_Retries429EvenForWrites(t *testing.T) {
	var calls atomic.Int32
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) < 2 {
			errJSON(w, http.StatusTooManyRequests, 429, "Too Many Requests", `{"retry_after":0}`)
			return
		}
		okJSON(w, `{"message_id":7}`)
	})
	_, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "x"}, false)
	require.NoError(t, err, "a 429 is safe to retry: the request was rejected, not processed")
	assert.EqualValues(t, 2, calls.Load())
}

func TestClient_NonJSONResponse(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>502 Bad Gateway</html>"))
	})
	_, err := c.Call(t.Context(), "getMe", nil, true)
	require.Error(t, err)
	var ae *APIError
	require.ErrorAs(t, err, &ae)
	assert.Contains(t, ae.Description, "non-JSON")
}

func TestClient_DryRun_PrintsRedactedCurl(t *testing.T) {
	var buf bytes.Buffer
	auth, err := NewBotTokenAuth("999:SECRETHASHVALUE12345")
	require.NoError(t, err)
	c := New(auth, WithDryRun(true), WithDryRunWriter(&buf))

	raw, err := c.Call(t.Context(), "sendMessage", map[string]any{"chat_id": "1", "text": "hi"}, false)
	require.NoError(t, err)
	assert.Nil(t, raw, "dry-run performs no request")

	out := buf.String()
	assert.Contains(t, out, "curl -sS -X POST")
	assert.Contains(t, out, "sendMessage")
	assert.Contains(t, out, "999:<redacted>")
	assert.NotContains(t, out, "SECRETHASHVALUE12345", "the secret must never appear in a dry-run")
}

func TestClient_DryRun_ShowToken(t *testing.T) {
	var buf bytes.Buffer
	auth, _ := NewBotTokenAuth("999:SECRETHASHVALUE12345")
	c := New(auth, WithDryRun(true), WithShowToken(true), WithDryRunWriter(&buf))
	_, err := c.Call(t.Context(), "getMe", nil, true)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "999:SECRETHASHVALUE12345")
}

func TestClient_Upload_Multipart(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "pic.jpg")
	require.NoError(t, os.WriteFile(imgPath, []byte("JPEGDATA"), 0o600))

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "1", r.FormValue("chat_id"))
		f, hdr, err := r.FormFile("photo")
		require.NoError(t, err)
		defer func() { _ = f.Close() }()
		assert.Equal(t, "pic.jpg", hdr.Filename)
		data, _ := io.ReadAll(f)
		assert.Equal(t, "JPEGDATA", string(data))
		okJSON(w, `{"message_id":5}`)
	})

	raw, err := c.Upload(t.Context(), "sendPhoto",
		map[string]any{"chat_id": "1"},
		map[string]string{"photo": imgPath}, false)
	require.NoError(t, err)
	assert.Contains(t, string(raw), "message_id")
}

func TestClient_CallInto_DecodesResult(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		okJSON(w, `{"id":42,"first_name":"X"}`)
	})
	var u User
	require.NoError(t, c.CallInto(t.Context(), "getChat", nil, true, &u))
	assert.Equal(t, ID("42"), u.ID)
}

func TestClient_DownloadFile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/file/bot")
		assert.True(t, strings.HasSuffix(r.URL.Path, "/photos/x.jpg"))
		_, _ = w.Write([]byte("FILEBYTES"))
	})
	var buf bytes.Buffer
	n, err := c.DownloadFile(t.Context(), "photos/x.jpg", &buf)
	require.NoError(t, err)
	assert.EqualValues(t, len("FILEBYTES"), n)
	assert.Equal(t, "FILEBYTES", buf.String())
}

func TestClient_DownloadFile_HTTPError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	var buf bytes.Buffer
	_, err := c.DownloadFile(t.Context(), "photos/missing.jpg", &buf)
	require.Error(t, err)
	var ae *APIError
	require.ErrorAs(t, err, &ae)
	assert.Equal(t, http.StatusNotFound, ae.Code)
}

func TestClient_RedactedFileURL(t *testing.T) {
	auth, _ := NewBotTokenAuth("999:SECRETHASHVALUE12345")
	c := New(auth, WithBaseURL("https://api.telegram.org"))
	url := c.RedactedFileURL("photos/x.jpg")
	assert.Equal(t, "https://api.telegram.org/file/bot999:<redacted>/photos/x.jpg", url)
	assert.NotContains(t, url, "SECRETHASHVALUE12345")
}

func TestClient_ContextCancelled(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		okJSON(w, `{}`)
	})
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // cancel before the call
	_, err := c.Call(ctx, "getMe", nil, true)
	require.Error(t, err)
}
