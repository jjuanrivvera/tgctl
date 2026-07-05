## tgctl stars send-paid-media

Send paid media that recipients unlock with Stars

### Synopsis

Send media locked behind a Star paywall. --media is a JSON array of InputPaidMedia objects.

```
tgctl stars send-paid-media [flags]
```

### Examples

```
  tgctl stars send-paid-media --chat @channel --star-count 50 \
    --media '[{"type":"photo","media":"https://e.com/a.jpg"}]'
```

### Options

```
      --business-connection-id string   business connection id on whose behalf to act
      --caption string                  media caption (0-1024 chars)
      --chat string                     target chat: numeric id or @username
  -h, --help                            help for send-paid-media
      --media string                    JSON array of InputPaidMedia objects
      --parse-mode string               text formatting: MarkdownV2 | HTML | Markdown
      --payload string                  bot-defined payload (not shown to users)
      --protect-content                 protect the content from forwarding and saving
      --silent                          send without a notification sound
      --star-count int                  Stars a user must pay to unlock (1-2500)
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

