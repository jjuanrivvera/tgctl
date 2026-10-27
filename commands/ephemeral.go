package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "ephemeral",
		Short: "Manage and edit ephemeral messages",
		Long: `Delete and edit ephemeral messages (visible only to a target user and the bot).

Ephemeral messages are sent in group or channel contexts to a specific user.
Their identifiers (ephemeral_message_id) are obtained from incoming updates or callback queries.
Targeting requires the chat, the recipient user id, and the ephemeral message id.
Delivery of edits and deletions is client-side and not guaranteed if the recipient is offline.`,
		Cmds: []methodCmd{
			{
				Use: "delete", Method: "deleteEphemeralMessage", Kind: kindDestructive,
				Short: "Delete an ephemeral message",
				Long: `Delete an ephemeral message (deleteEphemeralMessage).
Returns True on success.`,
				Example: `  tgctl ephemeral delete --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42`,
				Flags:   ephemeralTripleFlags(),
			},
			{
				Use: "edit-text", Method: "editEphemeralMessageText", Kind: kindWrite,
				Short: "Edit text or rich content of an ephemeral message",
				Long: `Edit the text or structured rich message of an ephemeral message (editEphemeralMessageText).
Either --text or --rich-message must be provided. Returns True on success.`,
				Example: `  tgctl ephemeral edit-text --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --text "Updated text"
  tgctl ephemeral edit-text --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --rich-message '{"blocks":[{"type":"paragraph","text":"New content"}]}'`,
				Flags: append(ephemeralTripleFlags(),
					flagSpec{Name: "text", Usage: "new text of the ephemeral message"},
					flagSpec{Name: "rich-message", Param: "rich_message", Kind: flagJSON, Usage: "JSON object describing updated rich message (InputRichMessage)"},
					parseModeFlag(),
					flagSpec{Name: "entities", Kind: flagJSON, Usage: "JSON array of MessageEntity objects"},
					flagSpec{Name: "link-preview-options", Param: "link_preview_options", Kind: flagJSON, Usage: "JSON object for link preview options"},
					replyMarkupFlag(),
				),
				PreCall: validateEphemeralEditText,
			},
			{
				Use: "edit-caption", Method: "editEphemeralMessageCaption", Kind: kindWrite,
				Short: "Edit caption of an ephemeral media message",
				Long: `Edit the caption of an ephemeral media message (editEphemeralMessageCaption).
Returns True on success.`,
				Example: `  tgctl ephemeral edit-caption --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --caption "New caption"`,
				Flags: append(ephemeralTripleFlags(),
					flagSpec{Name: "caption", Usage: "new caption of the media message"},
					parseModeFlag(),
					flagSpec{Name: "caption-entities", Param: "caption_entities", Kind: flagJSON, Usage: "JSON array of MessageEntity objects"},
					flagSpec{Name: "show-caption-above-media", Param: "show_caption_above_media", Kind: flagBool, Usage: "display caption above media"},
					replyMarkupFlag(),
				),
			},
			{
				Use: "edit-media", Method: "editEphemeralMessageMedia", Kind: kindWrite,
				Short: "Edit media content of an ephemeral message",
				Long: `Edit animation, audio, document, photo, or video of an ephemeral message (editEphemeralMessageMedia).
--media accepts an InputMedia JSON object. New file uploads can be referenced via attach:// and passed with --file. Returns True on success.`,
				Example: `  tgctl ephemeral edit-media --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --media '{"type":"photo","media":"https://example.com/pic.jpg"}'`,
				Flags: append(ephemeralTripleFlags(),
					flagSpec{Name: "media", Kind: flagJSON, Required: true, Usage: "JSON object describing new media content (InputMedia)"},
					replyMarkupFlag(),
				),
				AllowFiles: true,
				PreCall:    validateEphemeralEditMedia,
			},
			{
				Use: "edit-reply-markup", Method: "editEphemeralMessageReplyMarkup", Kind: kindWrite,
				Short: "Edit reply markup of an ephemeral message",
				Long: `Edit only the reply markup (inline keyboard) of an ephemeral message (editEphemeralMessageReplyMarkup).
Returns True on success.`,
				Example: `  tgctl ephemeral edit-reply-markup --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --reply-markup '{"inline_keyboard":[[{"text":"Refresh","callback_data":"ref"}]]}'`,
				Flags: append(ephemeralTripleFlags(),
					replyMarkupFlag(),
				),
			},
		},
	})
}

func ephemeralReceiverFlag() flagSpec {
	return flagSpec{
		Name:     "receiver",
		Param:    "receiver_user_id",
		Kind:     flagInt,
		Required: true,
		Usage:    "user id who received the ephemeral message",
	}
}

func ephemeralIDFlag() flagSpec {
	return flagSpec{
		Name:     "ephemeral-message-id",
		Param:    "ephemeral_message_id",
		Kind:     flagInt,
		Required: true,
		Usage:    "unique identifier of the ephemeral message",
	}
}

func ephemeralTripleFlags() []flagSpec {
	return []flagSpec{
		chatFlag(),
		ephemeralReceiverFlag(),
		ephemeralIDFlag(),
	}
}

func validateEphemeralEditText(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	_, hasText := params["text"]
	v, hasRich := params["rich_message"]
	if !hasText && !hasRich {
		return fmt.Errorf("pass either --text or --rich-message to edit text")
	}
	if hasRich {
		if _, ok := v.(map[string]any); !ok {
			return fmt.Errorf("--rich-message must be a JSON object")
		}
	}
	return nil
}

func validateEphemeralEditMedia(_ *cobra.Command, params map[string]any, files map[string]string) error {
	v, ok := params["media"]
	if !ok {
		return nil
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("--media must be a JSON object")
	}
	for _, ref := range findAttachReferences(obj) {
		if _, hasFile := files[ref]; !hasFile {
			return fmt.Errorf("missing upload for attach reference %q (pass --file %s=@<path>)", ref, ref)
		}
	}
	return nil
}
