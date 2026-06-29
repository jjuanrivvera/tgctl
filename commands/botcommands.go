package commands

func init() {
	registerGroup(group{
		Use:     "commands",
		Aliases: []string{"cmds"},
		Short:   "Manage the bot's command menu",
		Long:    "List, set, and delete the slash-command menu Telegram shows users (getMyCommands/setMyCommands).",
		Cmds: []methodCmd{
			{
				Use: "list", Method: "getMyCommands", Kind: kindRead,
				Short: "List the bot's commands",
				Example: `  tgctl commands list
  tgctl commands list --language-code es`,
				Flags:   []flagSpec{scopeFlag(), languageFlag()},
				Columns: []string{"command", "description"},
			},
			{
				Use: "set", Method: "setMyCommands", Kind: kindWrite,
				Short:   "Set the bot's command menu",
				Example: `  tgctl commands set --commands '[{"command":"start","description":"Begin"},{"command":"help","description":"Get help"}]'`,
				Flags: []flagSpec{
					{Name: "commands", Kind: flagJSON, Required: true, Usage: "array of {command,description} objects as JSON"},
					scopeFlag(), languageFlag(),
				},
			},
			{
				Use: "delete", Method: "deleteMyCommands", Kind: kindDestructive,
				Short: "Delete the bot's command menu",
				Example: `  tgctl commands delete
  tgctl commands delete --language-code es`,
				Flags: []flagSpec{scopeFlag(), languageFlag()},
			},
		},
	})
}

func scopeFlag() flagSpec {
	return flagSpec{Name: "scope", Kind: flagJSON, Usage: "BotCommandScope object as JSON (default: all private chats)"}
}

func languageFlag() flagSpec {
	return flagSpec{Name: "language-code", Param: "language_code", Usage: "BCP-47 language code"}
}
