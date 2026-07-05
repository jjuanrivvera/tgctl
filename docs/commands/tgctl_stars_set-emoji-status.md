## tgctl stars set-emoji-status

Set a user's emoji status (requires the user's prior consent)

```
tgctl stars set-emoji-status [flags]
```

### Examples

```
  tgctl stars set-emoji-status --user 12345 --emoji-status-custom-emoji-id 5170233102089322756
```

### Options

```
      --emoji-status-custom-emoji-id string   custom emoji id for the status (omit to remove)
      --emoji-status-expiration-date int      unix time the status expires
  -h, --help                                  help for set-emoji-status
      --user int                              target user id
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
      --no-store          disable local SQLite send/receive history for this invocation (see tgctl log)
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl stars](tgctl_stars.md)	 - Telegram Stars: transactions, gifts, and paid media

