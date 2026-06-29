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
