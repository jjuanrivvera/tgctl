package commands

func init() {
	registerGroup(group{
		Use:   "chat",
		Short: "Inspect chats and their members",
		Long:  "Read chat metadata, member counts, administrators, and individual members; leave a chat.",
		Cmds: []methodCmd{
			{
				Use: "get", Method: "getChat", Kind: kindRead,
				Short: "Show a chat's metadata",
				Example: `  tgctl chat get --chat @telegram
  tgctl chat get --chat -1001234567890 -o json`,
				Flags:   []flagSpec{chatFlag()},
				Columns: []string{"id", "type", "title", "username"},
			},
			{
				Use: "members-count", Method: "getChatMemberCount", Kind: kindRead,
				Short:   "Show the number of members in a chat",
				Example: `  tgctl chat members-count --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "administrators", Aliases: []string{"admins"}, Method: "getChatAdministrators", Kind: kindRead,
				Short:   "List a chat's administrators",
				Example: `  tgctl chat administrators --chat @group`,
				Flags:   []flagSpec{chatFlag()},
				Columns: []string{"status", "user.id", "user.username"},
			},
			{
				Use: "member", Method: "getChatMember", Kind: kindRead,
				Short:   "Show one member's status in a chat",
				Example: `  tgctl chat member --chat @group --user 12345`,
				Flags:   []flagSpec{chatFlag(), userFlag()},
				Columns: []string{"status", "user.id", "user.username"},
			},
			{
				Use: "leave", Method: "leaveChat", Kind: kindDestructive,
				Short:   "Make the bot leave a chat",
				Example: `  tgctl chat leave --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "set-title", Method: "setChatTitle", Kind: kindWrite,
				Short:   "Change a chat's title",
				Example: `  tgctl chat set-title --chat @group --title "New title"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "title", Required: true, Usage: "new chat title (1-128 chars)"},
				},
			},
			{
				Use: "set-description", Method: "setChatDescription", Kind: kindWrite,
				Short: "Change a chat's description",
				Example: `  tgctl chat set-description --chat @group --description "What this group is about"
  tgctl chat set-description --chat @group --description ""   # clear it`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "description", Usage: "new chat description (0-255 chars; empty clears it)"},
				},
			},
			{
				Use: "set-photo", Method: "setChatPhoto", Kind: kindWrite,
				Short:   "Set a chat's photo",
				Example: `  tgctl chat set-photo --chat @group --photo ./logo.png`,
				Flags:   []flagSpec{chatFlag()},
				Files:   []fileSpec{{Name: "photo", Required: true, Usage: "local image path (URLs/file_id not accepted by Telegram here)"}},
			},
			{
				Use: "delete-photo", Method: "deleteChatPhoto", Kind: kindDestructive,
				Short:   "Delete a chat's photo",
				Example: `  tgctl chat delete-photo --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "set-permissions", Method: "setChatPermissions", Kind: kindWrite,
				Short: "Set the default permissions for all members",
				Example: `  tgctl chat set-permissions --chat @group \
    --permissions '{"can_send_messages":true,"can_send_polls":false}'`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "permissions", Kind: flagJSON, Required: true, Usage: "ChatPermissions object as JSON"},
					{Name: "use-independent-chat-permissions", Param: "use_independent_chat_permissions", Kind: flagBool, Usage: "treat each permission independently"},
				},
			},
			{
				Use: "set-sticker-set", Method: "setChatStickerSet", Kind: kindWrite,
				Short:   "Set the group sticker set for a supergroup",
				Example: `  tgctl chat set-sticker-set --chat @group --sticker-set-name MyPack`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "sticker-set-name", Param: "sticker_set_name", Required: true, Usage: "name of the sticker set"},
				},
			},
			{
				Use: "delete-sticker-set", Method: "deleteChatStickerSet", Kind: kindDestructive,
				Short:   "Remove the group sticker set from a supergroup",
				Example: `  tgctl chat delete-sticker-set --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "menu-button", Method: "getChatMenuButton", Kind: kindRead,
				Short:   "Show the chat's menu button (default: the bot's global button)",
				Example: `  tgctl chat menu-button --chat 12345`,
				Flags: []flagSpec{
					{Name: "chat", Param: "chat_id", Kind: flagInt, Usage: "private chat id (omit for the bot's default button)"},
				},
			},
			{
				Use: "set-menu-button", Method: "setChatMenuButton", Kind: kindWrite,
				Short:   "Set the chat's menu button",
				Example: `  tgctl chat set-menu-button --chat 12345 --menu-button '{"type":"commands"}'`,
				Flags: []flagSpec{
					{Name: "chat", Param: "chat_id", Kind: flagInt, Usage: "private chat id (omit to set the bot's default button)"},
					{Name: "menu-button", Param: "menu_button", Kind: flagJSON, Usage: "MenuButton object as JSON (omit to reset to default)"},
				},
			},
			{
				Use: "unpin-all", Method: "unpinAllChatMessages", Kind: kindWrite,
				Short:   "Unpin every pinned message in a chat",
				Example: `  tgctl chat unpin-all --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "boosts", Method: "getUserChatBoosts", Kind: kindRead,
				Short:   "List the boosts a user added to a chat",
				Example: `  tgctl chat boosts --chat @channel --user 12345`,
				Flags:   []flagSpec{chatFlag(), userFlag()},
			},
		},
	})
}

func userFlag() flagSpec {
	return flagSpec{Name: "user", Param: "user_id", Kind: flagInt, Required: true, Usage: "target user id"}
}
