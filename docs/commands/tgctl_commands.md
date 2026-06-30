## tgctl commands

Manage the bot's command menu

### Synopsis

List, set, and delete the slash-command menu Telegram shows users (getMyCommands/setMyCommands).

### Options

```
  -h, --help   help for commands
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl commands delete](tgctl_commands_delete.md)	 - Delete the bot's command menu
* [tgctl commands list](tgctl_commands_list.md)	 - List the bot's commands
* [tgctl commands set](tgctl_commands_set.md)	 - Set the bot's command menu

