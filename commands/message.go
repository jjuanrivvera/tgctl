package commands

func init() {
	registerGroup(group{
		Use:     "message",
		Aliases: []string{"msg"},
		Short:   "Send and manage messages",
		Long:    "Send, edit, delete, forward, copy, and pin messages. --chat accepts a numeric id or @username.",
		Cmds: []methodCmd{
			{
				Use: "send", Method: "sendMessage", Kind: kindWrite,
				Short: "Send a text message",
				Example: `  tgctl message send --chat @me --text "hello"
  tgctl message send --chat -1001234567890 --text "*bold*" --parse-mode MarkdownV2 --silent`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "text", Required: true, Usage: "message text"},
					parseModeFlag(),
					{Name: "reply-to", Param: "reply_to_message_id", Kind: flagInt, Usage: "message id to reply to"},
					{Name: "thread", Param: "message_thread_id", Kind: flagInt, Usage: "forum topic thread id"},
					{Name: "silent", Param: "disable_notification", Kind: flagBool, Usage: "send without a notification sound"},
					{Name: "no-preview", Param: "disable_web_page_preview", Kind: flagBool, Usage: "disable link previews"},
					replyMarkupFlag(),
				},
				Columns: []string{"message_id", "chat.id", "date", "text"},
			},
			{
				Use: "edit", Method: "editMessageText", Kind: kindWrite,
				Short:   "Edit a message's text",
				Example: `  tgctl message edit --chat @me --message-id 42 --text "updated"`,
				Flags: []flagSpec{
					chatFlag(),
					messageIDFlag(),
					{Name: "text", Required: true, Usage: "new message text"},
					parseModeFlag(),
					replyMarkupFlag(),
				},
				Columns: []string{"message_id", "chat.id", "edit_date", "text"},
			},
			{
				Use: "delete", Method: "deleteMessage", Kind: kindDestructive,
				Short:   "Delete a message",
				Example: `  tgctl message delete --chat @me --message-id 42`,
				Flags:   []flagSpec{chatFlag(), messageIDFlag()},
			},
			{
				Use: "forward", Method: "forwardMessage", Kind: kindWrite,
				Short:   "Forward a message to another chat",
				Example: `  tgctl message forward --chat @dest --from-chat @src --message-id 42`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "from-chat", Param: "from_chat_id", Required: true, Usage: "source chat id or @username"},
					messageIDFlag(),
					{Name: "silent", Param: "disable_notification", Kind: flagBool, Usage: "forward without a notification"},
				},
				Columns: []string{"message_id", "chat.id", "forward_date"},
			},
			{
				Use: "copy", Method: "copyMessage", Kind: kindWrite,
				Short:   "Copy a message (without a 'forwarded from' header)",
				Example: `  tgctl message copy --chat @dest --from-chat @src --message-id 42 --caption "fyi"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "from-chat", Param: "from_chat_id", Required: true, Usage: "source chat id or @username"},
					messageIDFlag(),
					{Name: "caption", Usage: "new caption for the copied message"},
					parseModeFlag(),
				},
				Columns: []string{"message_id"},
			},
			{
				Use: "pin", Method: "pinChatMessage", Kind: kindWrite,
				Short:   "Pin a message in a chat",
				Example: `  tgctl message pin --chat @group --message-id 42 --silent`,
				Flags: []flagSpec{
					chatFlag(),
					messageIDFlag(),
					{Name: "silent", Param: "disable_notification", Kind: flagBool, Usage: "pin without notifying members"},
				},
			},
			{
				Use: "unpin", Method: "unpinChatMessage", Kind: kindWrite,
				Short: "Unpin a message (or the most recent pin) in a chat",
				Example: `  tgctl message unpin --chat @group --message-id 42
  tgctl message unpin --chat @group   # unpins the most recent`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "message-id", Param: "message_id", Kind: flagInt, Usage: "message id to unpin (omit for most recent)"},
				},
			},
			{
				Use: "react", Method: "setMessageReaction", Kind: kindWrite,
				Short: "Set (or clear) reactions on a message",
				Long:  "Set reactions on a message. --reaction is a JSON array of ReactionType objects; omit it to remove all reactions.",
				Example: `  tgctl message react --chat @group --message-id 42 --reaction '[{"type":"emoji","emoji":"👍"}]'
  tgctl message react --chat @group --message-id 42   # clear reactions`,
				Flags: []flagSpec{
					chatFlag(),
					messageIDFlag(),
					{Name: "reaction", Kind: flagJSON, Usage: "JSON array of ReactionType objects (omit to clear)"},
					{Name: "is-big", Param: "is_big", Kind: flagBool, Usage: "show the reaction with a big animation"},
				},
			},
			{
				Use: "location", Method: "sendLocation", Kind: kindWrite,
				Short:   "Send a point on the map",
				Example: `  tgctl message location --chat @me --latitude 3.4516 --longitude -76.532`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "latitude", Kind: flagFloat, Required: true, Usage: "latitude of the location"},
					{Name: "longitude", Kind: flagFloat, Required: true, Usage: "longitude of the location"},
					{Name: "live-period", Param: "live_period", Kind: flagInt, Usage: "seconds the location is updated live (60-86400)"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "venue", Method: "sendVenue", Kind: kindWrite,
				Short:   "Send information about a venue",
				Example: `  tgctl message venue --chat @me --latitude 3.45 --longitude -76.53 --title "Office" --address "Av. 1 #2-3"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "latitude", Kind: flagFloat, Required: true, Usage: "latitude of the venue"},
					{Name: "longitude", Kind: flagFloat, Required: true, Usage: "longitude of the venue"},
					{Name: "title", Required: true, Usage: "name of the venue"},
					{Name: "address", Required: true, Usage: "address of the venue"},
					{Name: "foursquare-id", Param: "foursquare_id", Usage: "Foursquare identifier of the venue"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "contact", Method: "sendContact", Kind: kindWrite,
				Short:   "Send a phone contact",
				Example: `  tgctl message contact --chat @me --phone-number "+15551234567" --first-name "Ada"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "phone-number", Param: "phone_number", Required: true, Usage: "contact's phone number"},
					{Name: "first-name", Param: "first_name", Required: true, Usage: "contact's first name"},
					{Name: "last-name", Param: "last_name", Usage: "contact's last name"},
					{Name: "vcard", Usage: "additional data about the contact as a vCard (0-2048 bytes)"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "poll", Method: "sendPoll", Kind: kindWrite,
				Short: "Send a native poll",
				Long:  "Send a poll. --options is a JSON array of InputPollOption objects, e.g. '[{\"text\":\"Yes\"},{\"text\":\"No\"}]'.",
				Example: `  tgctl message poll --chat @group --question "Lunch?" \
    --options '[{"text":"Pizza"},{"text":"Sushi"}]'`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "question", Required: true, Usage: "poll question (1-300 chars)"},
					{Name: "options", Kind: flagJSON, Required: true, Usage: "JSON array of InputPollOption objects (2-10)"},
					{Name: "anonymous", Param: "is_anonymous", Kind: flagBool, Usage: "make the poll anonymous (Telegram default: true)"},
					{Name: "type", Usage: "poll type: regular | quiz"},
					{Name: "allows-multiple-answers", Param: "allows_multiple_answers", Kind: flagBool, Usage: "allow multiple answers (regular polls only)"},
					{Name: "correct-option-id", Param: "correct_option_id", Kind: flagInt, Usage: "0-based id of the correct option (quiz polls)"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "dice", Method: "sendDice", Kind: kindWrite,
				Short:   "Send an animated emoji with a random value (dice, dart, etc.)",
				Example: `  tgctl message dice --chat @me --emoji 🎯`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "emoji", Usage: "one of 🎲 🎯 🏀 ⚽ 🎳 🎰 (default 🎲)"},
					silentFlag(),
				},
				Columns: []string{"message_id", "chat.id", "dice.emoji", "dice.value"},
			},
		},
	})
}

// Shared flag builders keep the common Bot API params consistent across commands.
func chatFlag() flagSpec {
	return flagSpec{Name: "chat", Param: "chat_id", Required: true, Usage: "target chat: numeric id or @username"}
}

func messageIDFlag() flagSpec {
	return flagSpec{Name: "message-id", Param: "message_id", Kind: flagInt, Required: true, Usage: "message id"}
}

func parseModeFlag() flagSpec {
	return flagSpec{Name: "parse-mode", Param: "parse_mode", Usage: "text formatting: MarkdownV2 | HTML | Markdown"}
}

func replyMarkupFlag() flagSpec {
	return flagSpec{Name: "reply-markup", Param: "reply_markup", Kind: flagJSON, Usage: "inline/reply keyboard as JSON"}
}
