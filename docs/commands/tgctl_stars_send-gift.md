## tgctl stars send-gift

Send a gift to a user or channel

```
tgctl stars send-gift [flags]
```

### Examples

```
  tgctl stars send-gift --user 12345 --gift-id 5170233102089322756 --text "Enjoy!"
```

### Options

```
      --chat string              recipient channel chat id or @username (or use --user)
      --gift-id string           id of the gift to send (from stars gifts)
  -h, --help                     help for send-gift
      --pay-for-upgrade          pay for the gift's upgrade to a unique gift
      --text string              text shown with the gift (0-128 chars)
      --text-parse-mode string   parse mode for --text (MarkdownV2 | HTML)
      --user int                 recipient user id (or use --chat)
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

