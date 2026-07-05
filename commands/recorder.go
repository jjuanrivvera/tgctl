package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/config"
	"github.com/jjuanrivvera/tgctl/internal/store"
)

// messageBearingMethods maps a Bot API method worth recording to the `kind` label stored
// alongside it. Only methods that produce a message a user would later want to look up are
// here — getMe, getChat, etc. are deliberately excluded. getUpdates is ALSO excluded: its
// result is a batch of inbound Updates, not a sent message, and is recorded separately with
// direction='in' by commands/updates.go and commands/webhook_listen.go. Adding a new send
// command never needs a new recording edit: extend this map instead (DECISIONS.md).
var messageBearingMethods = map[string]string{
	"sendMessage":        "text",
	"sendPhoto":          "photo",
	"sendDocument":       "document",
	"sendAudio":          "audio",
	"sendVoice":          "voice",
	"sendVideo":          "video",
	"sendVideoNote":      "video_note",
	"sendAnimation":      "animation",
	"sendSticker":        "sticker",
	"sendMediaGroup":     "media_group",
	"sendLocation":       "location",
	"sendVenue":          "venue",
	"sendContact":        "contact",
	"sendPoll":           "poll",
	"sendDice":           "dice",
	"copyMessage":        "copy",
	"forwardMessage":     "forward",
	"editMessageText":    "edit",
	"editMessageCaption": "edit_caption",
	"editMessageMedia":   "edit_media",
}

// storeRecorder adapts *store.Store to api.Recorder structurally (internal/api never imports
// internal/store — see internal/api/recorder.go). It is the only place that knows how to pull
// chat_id/message_id/file_id out of Bot API params/results.
type storeRecorder struct {
	st    *store.Store
	quiet bool // suppress the stderr warning on a store write failure, matching --quiet
}

// Record implements api.Recorder. It is a fire-and-forget observer: messageBearingMethods
// decides whether method matters at all, and any store failure is logged (unless --quiet) and
// swallowed — the send this call is observing has already succeeded (DECISIONS.md).
func (r *storeRecorder) Record(ctx context.Context, method string, params map[string]any, result json.RawMessage) {
	kind, ok := messageBearingMethods[method]
	if !ok {
		return
	}
	msg := buildOutboundMessage(kind, params, result)
	if err := r.st.Record(ctx, msg); err != nil && !r.quiet {
		fmt.Fprintf(os.Stderr, "tgctl: warning: message store write failed: %v\n", err)
	}
}

// Close implements io.Closer so api.Client.Close (called via defer at every clientFromCmd call
// site) releases the store's SQLite file handle. Without this the handle stays open for the
// life of the process — harmless on Unix, but it blocks Windows from deleting/renaming the file
// (e.g. a test's t.TempDir() cleanup) until something closes it.
func (r *storeRecorder) Close() error {
	if r.st == nil {
		return nil
	}
	return r.st.Close()
}

// buildOutboundMessage extracts a store.Message from a successful send/edit call. The API's
// own response is preferred over the request params wherever both could carry a field (e.g.
// chat_id): Telegram always resolves and echoes back the numeric chat/message id in the
// result, even when the request targeted a chat by @username, so the result is the more
// reliable source.
func buildOutboundMessage(kind string, params map[string]any, result json.RawMessage) store.Message {
	tm, _ := parseTelegramMessage(result)
	msg := store.Message{
		Direction: "out",
		Kind:      kind,
		ChatID:    firstNonZero(tm.Chat.ID, toInt64(params["chat_id"])),
		MessageID: tm.MessageID,
		Text:      config.FirstNonEmpty(tm.textOrCaption(), paramString(params, "text"), paramString(params, "caption")),
		FileID:    tm.fileID(kind),
		Raw:       result,
	}
	if v, ok := params["reply_to_message_id"]; ok {
		msg.ReplyToMessageID = toInt64(v)
	}
	return msg
}

// recordInboundMessage persists one inbound Message (direction 'in') seen via getUpdates
// (commands/updates.go) or a webhook delivery (commands/webhook_listen.go). Best-effort like
// storeRecorder.Record: a failure is logged (unless --quiet) and swallowed, never surfaced as
// a command error — the update has already been fetched/received by the time this runs.
func recordInboundMessage(cmd *cobra.Command, st *store.Store, tm *telegramMessage) {
	kind := tm.kindFromContent()
	msg := store.Message{
		Direction: "in",
		Kind:      kind,
		ChatID:    tm.Chat.ID,
		MessageID: tm.MessageID,
		Text:      tm.textOrCaption(),
		FileID:    tm.fileID(kind),
	}
	if tm.ReplyToMessage != nil {
		msg.ReplyToMessageID = tm.ReplyToMessage.MessageID
	}
	if raw, err := json.Marshal(tm); err == nil {
		msg.Raw = raw
	}
	quiet, _ := cmd.Flags().GetBool("quiet")
	if err := st.Record(cmd.Context(), msg); err != nil && !quiet {
		fmt.Fprintf(os.Stderr, "tgctl: warning: message store write failed: %v\n", err)
	}
}

// telegramMessage is the subset of the Bot API's Message object
// (https://core.telegram.org/bots/api#message) both the outbound recorder and the inbound
// update recorder need. Sharing one struct avoids declaring the same field list twice.
type telegramMessage struct {
	MessageID int64 `json:"message_id"`
	Chat      struct {
		ID int64 `json:"id"`
	} `json:"chat"`
	Text           string           `json:"text"`
	Caption        string           `json:"caption"`
	Photo          []fileIDField    `json:"photo"` // largest size is last
	Document       *fileIDField     `json:"document"`
	Audio          *fileIDField     `json:"audio"`
	Voice          *fileIDField     `json:"voice"`
	Video          *fileIDField     `json:"video"`
	VideoNote      *fileIDField     `json:"video_note"`
	Animation      *fileIDField     `json:"animation"`
	Sticker        *fileIDField     `json:"sticker"`
	ReplyToMessage *telegramMessage `json:"reply_to_message"`
}

type fileIDField struct {
	FileID string `json:"file_id"`
}

func (tm telegramMessage) textOrCaption() string {
	if tm.Text != "" {
		return tm.Text
	}
	return tm.Caption
}

// fileID returns the file_id for the media field matching kind, or "" when kind isn't a media
// kind (e.g. "text") or the field is absent (e.g. an edit that only touched text).
func (tm telegramMessage) fileID(kind string) string {
	switch kind {
	case "photo":
		if len(tm.Photo) > 0 {
			return tm.Photo[len(tm.Photo)-1].FileID
		}
	case "document":
		return holderFileID(tm.Document)
	case "audio":
		return holderFileID(tm.Audio)
	case "voice":
		return holderFileID(tm.Voice)
	case "video":
		return holderFileID(tm.Video)
	case "video_note":
		return holderFileID(tm.VideoNote)
	case "animation":
		return holderFileID(tm.Animation)
	case "sticker":
		return holderFileID(tm.Sticker)
	}
	return ""
}

// kindFromContent infers a `kind` label for an INBOUND update, which (unlike an outbound send)
// carries no method name to key off of — only the Message object itself.
func (tm telegramMessage) kindFromContent() string {
	switch {
	case tm.Text != "":
		return "text"
	case len(tm.Photo) > 0:
		return "photo"
	case tm.Document != nil:
		return "document"
	case tm.Audio != nil:
		return "audio"
	case tm.Voice != nil:
		return "voice"
	case tm.Video != nil:
		return "video"
	case tm.VideoNote != nil:
		return "video_note"
	case tm.Animation != nil:
		return "animation"
	case tm.Sticker != nil:
		return "sticker"
	default:
		return "other"
	}
}

func holderFileID(h *fileIDField) string {
	if h == nil {
		return ""
	}
	return h.FileID
}

// parseTelegramMessage decodes result as a single Message object. sendMediaGroup returns an
// ARRAY of Message objects; the first element is used (DECISIONS.md: one recorded row per
// call, not per message, keeps the write path a single INSERT). A result that isn't shaped
// like a Message at all (e.g. `true` from an inline-message edit) yields a zero value and
// ok=false — callers already treat every field of a zero telegramMessage as "unknown".
func parseTelegramMessage(result json.RawMessage) (telegramMessage, bool) {
	var tm telegramMessage
	if err := json.Unmarshal(result, &tm); err == nil && (tm.MessageID != 0 || tm.Chat.ID != 0) {
		return tm, true
	}
	var arr []telegramMessage
	if err := json.Unmarshal(result, &arr); err == nil && len(arr) > 0 {
		return arr[0], true
	}
	return telegramMessage{}, false
}

func paramString(params map[string]any, key string) string {
	s, _ := params[key].(string)
	return s
}

// toInt64 coerces a Bot API param/result value to int64. Params flow in as whatever the flag
// layer produced (string for chat_id/@username, int64 for message_id) or whatever `tgctl api
// -d`'s raw JSON decoded to (float64 for a JSON number); a non-numeric string (a bare
// @username, with no result available to resolve it) safely yields 0 rather than an error,
// since ChatID is a best-effort convenience field, not the row's only chat reference (Raw
// keeps the full request/response for forensics).
func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case float64:
		return int64(t)
	case json.Number:
		n, _ := t.Int64()
		return n
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

func firstNonZero(vals ...int64) int64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}
