package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "media",
		Short: "Send files: photos, documents, and video",
		Long: `Send media to a chat. Each --<kind> flag accepts a local file path (uploaded as
multipart/form-data), an http(s) URL (Telegram fetches it), or an existing file_id.`,
		Cmds: []methodCmd{
			{
				Use: "photo", Method: "sendPhoto", Kind: kindWrite,
				Short: "Send a photo",
				Example: `  tgctl media photo --chat @me --photo ./cat.jpg --caption "my cat"
  tgctl media photo --chat @me --photo https://example.com/pic.png`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "photo caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "photo", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "document", Method: "sendDocument", Kind: kindWrite,
				Short:   "Send a document/file",
				Example: `  tgctl media document --chat @me --document ./report.pdf --caption "Q2 report"`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "document caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "document", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "video", Method: "sendVideo", Kind: kindWrite,
				Short:   "Send a video",
				Example: `  tgctl media video --chat @me --video ./clip.mp4`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "video caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "video", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "audio", Method: "sendAudio", Kind: kindWrite,
				Short:   "Send an audio file (shown in the music player)",
				Example: `  tgctl media audio --chat @me --audio ./song.mp3 --performer "Artist" --title "Track"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "caption", Usage: "audio caption"}, parseModeFlag(),
					{Name: "duration", Kind: flagInt, Usage: "duration in seconds"},
					{Name: "performer", Usage: "performer name"},
					{Name: "title", Usage: "track title"},
					silentFlag(),
				},
				Files:   []fileSpec{{Name: "audio", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "voice", Method: "sendVoice", Kind: kindWrite,
				Short:   "Send a voice message (OGG/OPUS, shown as a waveform)",
				Example: `  tgctl media voice --chat @me --voice ./note.ogg --duration 7`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "caption", Usage: "voice caption"}, parseModeFlag(),
					{Name: "duration", Kind: flagInt, Usage: "duration in seconds"},
					silentFlag(),
				},
				Files:   []fileSpec{{Name: "voice", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "animation", Method: "sendAnimation", Kind: kindWrite,
				Short:   "Send an animation (GIF or H.264/MPEG-4 without sound)",
				Example: `  tgctl media animation --chat @me --animation ./loop.gif --caption "nice"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "caption", Usage: "animation caption"}, parseModeFlag(),
					{Name: "duration", Kind: flagInt, Usage: "duration in seconds"},
					silentFlag(),
				},
				Files:   []fileSpec{{Name: "animation", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "video-note", Method: "sendVideoNote", Kind: kindWrite,
				Short:   "Send a video note (round video, up to 1 minute)",
				Example: `  tgctl media video-note --chat @me --video-note ./round.mp4 --length 240`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "duration", Kind: flagInt, Usage: "duration in seconds"},
					{Name: "length", Kind: flagInt, Usage: "video width and height (it is square)"},
					silentFlag(),
				},
				Files:   []fileSpec{{Name: "video-note", Param: "video_note", Required: true, Usage: "local path or file_id (URLs not supported by Telegram)"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "sticker", Method: "sendSticker", Kind: kindWrite,
				Short:   "Send a sticker (.WEBP, .TGS, or .WEBM)",
				Example: `  tgctl media sticker --chat @me --sticker CAACAgIAAxkBA...`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "emoji", Usage: "emoji associated with the uploaded sticker"},
					silentFlag(),
				},
				Files:   []fileSpec{{Name: "sticker", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "media-group", Aliases: []string{"album"}, Method: "sendMediaGroup", Kind: kindWrite,
				Short: "Send a group of photos/videos/documents as an album",
				Long: `Send 2-10 items as a single album. --media is a JSON array of InputMedia objects;
each item's "media" is an http(s) URL or an existing file_id (multipart attach:// uploads
are not supported here — upload first with the single-item commands if you need local files).`,
				Example: `  tgctl media media-group --chat @me \
    --media '[{"type":"photo","media":"https://e.com/a.jpg"},{"type":"photo","media":"https://e.com/b.jpg"}]'`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "media", Kind: flagJSON, Required: true, Usage: "JSON array of InputMedia objects"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "live-photo", Method: "sendLivePhoto", Kind: kindWrite,
				Short: "Send a live photo",
				Long: `Send a live photo (a still photo accompanied by a short video).
--photo and --live-photo both accept a local file path or an existing file_id.
HTTP(S) URLs are not supported by the Bot API for this method.`,
				Example: `  tgctl media live-photo --chat @me --photo ./still.jpg --live-photo ./motion.mp4 --caption "look at this"
  tgctl media live-photo --chat @me --photo AgACAgI... --live-photo BAACAgI...`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "business-connection-id", Param: "business_connection_id", Usage: "unique identifier of the business connection"},
					threadFlag(),
					{Name: "caption", Usage: "live photo caption"},
					parseModeFlag(),
					{Name: "caption-entities", Param: "caption_entities", Kind: flagJSON, Usage: "JSON array of special entities that appear in the caption"},
					{Name: "show-caption-above-media", Param: "show_caption_above_media", Kind: flagBool, Usage: "display caption above media"},
					{Name: "has-spoiler", Param: "has_spoiler", Kind: flagBool, Usage: "cover the live photo with a spoiler animation"},
					silentFlag(),
					protectContentFlag(),
					{Name: "allow-paid-broadcast", Param: "allow_paid_broadcast", Kind: flagBool, Usage: "allow up to 1000 messages per second by paying Telegram Stars"},
					{Name: "message-effect-id", Param: "message_effect_id", Usage: "unique identifier of the message effect"},
					{Name: "reply-parameters", Param: "reply_parameters", Kind: flagJSON, Usage: "JSON object describing message to reply to"},
					replyMarkupFlag(),
				},
				Files: []fileSpec{
					{Name: "live-photo", Param: "live_photo", Required: true, Usage: "live photo video part (local path or file_id; URLs not supported)"},
					{Name: "photo", Param: "photo", Required: true, Usage: "still photo part (local path or file_id; URLs not supported)"},
				},
				PreCall: requireNoURLLivePhoto,
				Columns: []string{"message_id", "chat.id"},
			},
		},
	})
}

// requireNoURLLivePhoto rejects HTTP(S) URLs locally for sendLivePhoto because the Bot API
// specification requires a local file upload or an existing Telegram file_id for both parts.
func requireNoURLLivePhoto(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	for _, field := range []string{"live_photo", "photo"} {
		if v, ok := params[field].(string); ok {
			if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
				flagName := strings.ReplaceAll(field, "_", "-")
				return fmt.Errorf("--%s does not support HTTP(S) URLs: pass a local file path or an existing Telegram file_id", flagName)
			}
		}
	}
	return nil
}

func silentFlag() flagSpec {
	return flagSpec{Name: "silent", Param: "disable_notification", Kind: flagBool, Usage: "send without a notification sound"}
}
