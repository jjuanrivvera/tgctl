package commands

func init() {
	registerGroup(group{
		Use:   "bot",
		Short: "Inspect and configure the bot itself",
		Long:  "Read the bot's identity (getMe) and manage its name/description shown in Telegram.",
		Cmds: []methodCmd{
			{
				Use: "info", Method: "getMe", Kind: kindRead,
				Short:   "Show the authenticated bot's identity (getMe)",
				Example: "  tgctl bot info\n  tgctl bot info -o json",
				Columns: []string{"id", "username", "first_name", "can_join_groups"},
			},
			{
				Use: "set-name", Method: "setMyName", Kind: kindWrite,
				Short:   "Set the bot's name",
				Example: `  tgctl bot set-name --name "My Helper Bot"`,
				Flags: []flagSpec{
					{Name: "name", Required: true, Usage: "new bot name (0-64 chars)"},
					{Name: "language-code", Param: "language_code", Usage: "BCP-47 code this name applies to"},
				},
			},
			{
				Use: "get-name", Method: "getMyName", Kind: kindRead,
				Short: "Get the bot's name",
				Flags: []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-description", Method: "setMyDescription", Kind: kindWrite,
				Short: "Set the bot's description (shown in the empty chat)",
				Flags: []flagSpec{
					{Name: "description", Required: true, Usage: "new description (0-512 chars)"},
					{Name: "language-code", Param: "language_code", Usage: "language this description applies to"},
				},
			},
			{
				Use: "get-description", Method: "getMyDescription", Kind: kindRead,
				Short: "Get the bot's description",
				Flags: []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
		},
	})
}
