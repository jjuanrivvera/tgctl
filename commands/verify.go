package commands

func init() {
	registerGroup(group{
		Use:   "verify",
		Short: "Verify or unverify chats and users",
		Long:  "Add or remove the verification mark on a chat or user. Only available to bots acting on behalf of an organization that Telegram authorized to verify.",
		Cmds: []methodCmd{
			{
				Use: "chat", Method: "verifyChat", Kind: kindWrite,
				Short:   "Verify a chat on behalf of the bot's organization",
				Example: `  tgctl verify chat --chat @group --custom-description "Official channel"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "custom-description", Param: "custom_description", Usage: "custom verification description (0-70 chars)"},
				},
			},
			{
				Use: "user", Method: "verifyUser", Kind: kindWrite,
				Short:   "Verify a user on behalf of the bot's organization",
				Example: `  tgctl verify user --user 12345 --custom-description "Verified staff"`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "custom-description", Param: "custom_description", Usage: "custom verification description (0-70 chars)"},
				},
			},
			{
				Use: "remove-chat", Method: "removeChatVerification", Kind: kindDestructive,
				Short:   "Remove verification from a chat",
				Example: `  tgctl verify remove-chat --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "remove-user", Method: "removeUserVerification", Kind: kindDestructive,
				Short:   "Remove verification from a user",
				Example: `  tgctl verify remove-user --user 12345`,
				Flags:   []flagSpec{userFlag()},
			},
		},
	})
}
