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
		},
	})
}

func userFlag() flagSpec {
	return flagSpec{Name: "user", Param: "user_id", Kind: flagInt, Required: true, Usage: "target user id"}
}
