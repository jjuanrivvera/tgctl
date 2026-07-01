## tgctl bot

Inspect and configure the bot itself

### Synopsis

Read the bot's identity (getMe) and manage its name/description shown in Telegram.

### Options

```
  -h, --help   help for bot
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
* [tgctl bot close](tgctl_bot_close.md)	 - Close the bot instance before moving it to another server
* [tgctl bot get-admin-rights](tgctl_bot_get-admin-rights.md)	 - Get the bot's default administrator rights
* [tgctl bot get-description](tgctl_bot_get-description.md)	 - Get the bot's description
* [tgctl bot get-name](tgctl_bot_get-name.md)	 - Get the bot's name
* [tgctl bot get-short-description](tgctl_bot_get-short-description.md)	 - Get the bot's short description
* [tgctl bot info](tgctl_bot_info.md)	 - Show the authenticated bot's identity (getMe)
* [tgctl bot logout](tgctl_bot_logout.md)	 - Log out from the cloud Bot API before running a local Bot API server
* [tgctl bot set-admin-rights](tgctl_bot_set-admin-rights.md)	 - Set the bot's default administrator rights (requested when added to a group/channel)
* [tgctl bot set-description](tgctl_bot_set-description.md)	 - Set the bot's description (shown in the empty chat)
* [tgctl bot set-name](tgctl_bot_set-name.md)	 - Set the bot's name
* [tgctl bot set-short-description](tgctl_bot_set-short-description.md)	 - Set the bot's short description (shown on the profile page)

