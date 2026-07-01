## tgctl chat set-menu-button

Set the chat's menu button

```
tgctl chat set-menu-button [flags]
```

### Examples

```
  tgctl chat set-menu-button --chat 12345 --menu-button '{"type":"commands"}'
```

### Options

```
      --chat int             private chat id (omit to set the bot's default button)
  -h, --help                 help for set-menu-button
      --menu-button string   MenuButton object as JSON (omit to reset to default)
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

* [tgctl chat](tgctl_chat.md)	 - Inspect chats and their members

