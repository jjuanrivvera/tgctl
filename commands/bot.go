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
				Short:   "Get the bot's name",
				Example: "  tgctl bot get-name\n  tgctl bot get-name --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-description", Method: "setMyDescription", Kind: kindWrite,
				Short:   "Set the bot's description (shown in the empty chat)",
				Example: `  tgctl bot set-description --description "I help you manage your groups."`,
				Flags: []flagSpec{
					{Name: "description", Required: true, Usage: "new description (0-512 chars)"},
					{Name: "language-code", Param: "language_code", Usage: "language this description applies to"},
				},
			},
			{
				Use: "get-description", Method: "getMyDescription", Kind: kindRead,
				Short:   "Get the bot's description",
				Example: "  tgctl bot get-description\n  tgctl bot get-description --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-short-description", Method: "setMyShortDescription", Kind: kindWrite,
				Short:   "Set the bot's short description (shown on the profile page)",
				Example: `  tgctl bot set-short-description --short-description "Group management, done right."`,
				Flags: []flagSpec{
					{Name: "short-description", Param: "short_description", Usage: "new short description (0-120 chars; empty clears it)"},
					{Name: "language-code", Param: "language_code", Usage: "language this description applies to"},
				},
			},
			{
				Use: "get-short-description", Method: "getMyShortDescription", Kind: kindRead,
				Short:   "Get the bot's short description",
				Example: "  tgctl bot get-short-description\n  tgctl bot get-short-description --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-admin-rights", Method: "setMyDefaultAdministratorRights", Kind: kindWrite,
				Short:   "Set the bot's default administrator rights (requested when added to a group/channel)",
				Example: `  tgctl bot set-admin-rights --rights '{"can_manage_chat":true,"can_delete_messages":true}'`,
				Flags: []flagSpec{
					{Name: "rights", Kind: flagJSON, Usage: "ChatAdministratorRights object as JSON (omit to clear)"},
					{Name: "for-channels", Param: "for_channels", Kind: flagBool, Usage: "apply to channels instead of groups/supergroups"},
				},
			},
			{
				Use: "get-admin-rights", Method: "getMyDefaultAdministratorRights", Kind: kindRead,
				Short:   "Get the bot's default administrator rights",
				Example: "  tgctl bot get-admin-rights\n  tgctl bot get-admin-rights --for-channels -o json",
				Flags: []flagSpec{
					{Name: "for-channels", Param: "for_channels", Kind: flagBool, Usage: "query the channel rights instead of group rights"},
				},
			},
			{
				Use: "close", Method: "close", Kind: kindWrite,
				Short:   "Close the bot instance before moving it to another server",
				Long:    "Close the bot instance (frees server resources). Returns an error for the first 10 minutes after the bot launches.",
				Example: `  tgctl bot close`,
			},
			{
				Use: "logout", Method: "logOut", Kind: kindDestructive,
				Short:   "Log out from the cloud Bot API before running a local Bot API server",
				Long:    "Log the bot out of the cloud Bot API. After this you can use a local Bot API server; you must re-login via api.telegram.org to switch back. Returns an error for the first 10 minutes after launch.",
				Example: `  tgctl bot logout`,
			},
		},
	})
}
