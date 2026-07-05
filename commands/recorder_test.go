package commands

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/tgctl/internal/store"
)

func TestBuildOutboundMessage(t *testing.T) {
	tests := []struct {
		name       string
		kind       string
		params     map[string]any
		result     string
		wantChat   int64
		wantMsgID  int64
		wantText   string
		wantFileID string
	}{
		{
			name:      "text send: chat_id and message_id come from the result",
			kind:      "text",
			params:    map[string]any{"chat_id": "@me", "text": "hi"},
			result:    `{"message_id":42,"chat":{"id":7},"text":"hi"}`,
			wantChat:  7,
			wantMsgID: 42,
			wantText:  "hi",
		},
		{
			name:      "chat_id falls back to params when the result omits it",
			kind:      "text",
			params:    map[string]any{"chat_id": "123", "text": "hi"},
			result:    `true`,
			wantChat:  123,
			wantMsgID: 0,
			wantText:  "hi", // no result text; falls back to params["text"]
		},
		{
			name:       "photo: file_id comes from the largest (last) size",
			kind:       "photo",
			params:     map[string]any{"chat_id": "7", "caption": "a pic"},
			result:     `{"message_id":1,"chat":{"id":7},"caption":"a pic","photo":[{"file_id":"small"},{"file_id":"BIG"}]}`,
			wantChat:   7,
			wantMsgID:  1,
			wantText:   "a pic",
			wantFileID: "BIG",
		},
		{
			name:       "document: file_id from the document object",
			kind:       "document",
			params:     map[string]any{"chat_id": "7"},
			result:     `{"message_id":2,"chat":{"id":7},"document":{"file_id":"DOC1"}}`,
			wantChat:   7,
			wantMsgID:  2,
			wantFileID: "DOC1",
		},
		{
			name:      "sendMediaGroup: first element of the array is used",
			kind:      "media_group",
			params:    map[string]any{"chat_id": "7"},
			result:    `[{"message_id":10,"chat":{"id":7}},{"message_id":11,"chat":{"id":7}}]`,
			wantChat:  7,
			wantMsgID: 10,
		},
		{
			name:      "non-numeric chat_id param with no usable result falls back to 0",
			kind:      "text",
			params:    map[string]any{"chat_id": "@channel", "text": "x"},
			result:    `not-a-message-shape`,
			wantChat:  0,
			wantMsgID: 0,
			wantText:  "x",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := buildOutboundMessage(tc.kind, tc.params, json.RawMessage(tc.result))
			assert.Equal(t, "out", msg.Direction)
			assert.Equal(t, tc.kind, msg.Kind)
			assert.Equal(t, tc.wantChat, msg.ChatID)
			assert.Equal(t, tc.wantMsgID, msg.MessageID)
			assert.Equal(t, tc.wantText, msg.Text)
			assert.Equal(t, tc.wantFileID, msg.FileID)
			assert.Equal(t, tc.result, string(msg.Raw), "the full raw result is always preserved for forensics")
		})
	}
}

func TestBuildOutboundMessage_ReplyToMessageID(t *testing.T) {
	msg := buildOutboundMessage("text",
		map[string]any{"chat_id": "7", "text": "hi", "reply_to_message_id": int64(99)},
		json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))
	assert.EqualValues(t, 99, msg.ReplyToMessageID)
}

func TestStoreRecorder_SkipsNonMessageBearingMethod(t *testing.T) {
	// getMe is not in messageBearingMethods; Record must return before ever touching r.st, so
	// a nil store here proves the method filter runs first (touching r.st would panic).
	r := &storeRecorder{st: nil}
	assert.NotPanics(t, func() {
		r.Record(t.Context(), "getMe", nil, json.RawMessage(`{"id":1}`))
	})
}

func TestParseTelegramMessage(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		tm, ok := parseTelegramMessage(json.RawMessage(`{"message_id":5,"chat":{"id":9},"text":"hi"}`))
		require.True(t, ok)
		assert.EqualValues(t, 5, tm.MessageID)
		assert.EqualValues(t, 9, tm.Chat.ID)
	})
	t.Run("array takes first element", func(t *testing.T) {
		tm, ok := parseTelegramMessage(json.RawMessage(`[{"message_id":1},{"message_id":2}]`))
		require.True(t, ok)
		assert.EqualValues(t, 1, tm.MessageID)
	})
	t.Run("bare bool is not a message", func(t *testing.T) {
		_, ok := parseTelegramMessage(json.RawMessage(`true`))
		assert.False(t, ok)
	})
	t.Run("empty array is not a message", func(t *testing.T) {
		_, ok := parseTelegramMessage(json.RawMessage(`[]`))
		assert.False(t, ok)
	})
}

func TestTelegramMessage_KindFromContent(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"text", `{"text":"hi"}`, "text"},
		{"photo", `{"photo":[{"file_id":"a"}]}`, "photo"},
		{"document", `{"document":{"file_id":"a"}}`, "document"},
		{"audio", `{"audio":{"file_id":"a"}}`, "audio"},
		{"voice", `{"voice":{"file_id":"a"}}`, "voice"},
		{"video", `{"video":{"file_id":"a"}}`, "video"},
		{"video_note", `{"video_note":{"file_id":"a"}}`, "video_note"},
		{"animation", `{"animation":{"file_id":"a"}}`, "animation"},
		{"sticker", `{"sticker":{"file_id":"a"}}`, "sticker"},
		{"other", `{"contact":{"phone_number":"1"}}`, "other"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var tm telegramMessage
			require.NoError(t, json.Unmarshal([]byte(tc.json), &tm))
			assert.Equal(t, tc.want, tm.kindFromContent())
		})
	}
}

func TestTelegramMessage_FileID(t *testing.T) {
	var tm telegramMessage
	require.NoError(t, json.Unmarshal([]byte(`{
		"photo": [{"file_id":"small"},{"file_id":"BIG"}],
		"document": {"file_id":"DOC"},
		"video": {"file_id":"VID"}
	}`), &tm))
	assert.Equal(t, "BIG", tm.fileID("photo"), "the largest (last) photo size wins")
	assert.Equal(t, "DOC", tm.fileID("document"))
	assert.Equal(t, "VID", tm.fileID("video"))
	assert.Empty(t, tm.fileID("text"), "text is not a media kind")
	assert.Empty(t, tm.fileID("audio"), "the audio field is absent from this message")
}

func TestToInt64(t *testing.T) {
	assert.EqualValues(t, 7, toInt64(int64(7)))
	assert.EqualValues(t, 7, toInt64(7))
	assert.EqualValues(t, 7, toInt64(float64(7)))
	assert.EqualValues(t, 7, toInt64(json.Number("7")))
	assert.EqualValues(t, 7, toInt64("7"))
	assert.EqualValues(t, 0, toInt64("@username"))
	assert.EqualValues(t, 0, toInt64(nil))
	assert.EqualValues(t, 0, toInt64(true))
}

func TestFirstNonZero(t *testing.T) {
	assert.EqualValues(t, 5, firstNonZero(0, 5, 9))
	assert.EqualValues(t, 0, firstNonZero(0, 0))
	assert.EqualValues(t, 0, firstNonZero())
}

// TestStoreRecorder_Close pins the Windows CI regression fix: storeRecorder must implement
// io.Closer and actually close the underlying *store.Store, not just satisfy the interface.
// Record()ing after Close() failing proves the real sql.DB handle was released, not merely
// that Close() returned nil without doing anything.
func TestStoreRecorder_Close(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "default.db"))
	require.NoError(t, err)
	r := &storeRecorder{st: st}

	require.NoError(t, r.Close())
	err = st.Record(t.Context(), store.Message{Direction: "out", ChatID: 1, Kind: "text", Text: "after close"})
	assert.Error(t, err, "the store's sql.DB must actually be closed, not just no-op'd")
}

func TestStoreRecorder_Close_NilStoreIsNoop(t *testing.T) {
	r := &storeRecorder{}
	assert.NoError(t, r.Close())
}
