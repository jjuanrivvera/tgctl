## tgctl stars transfer-gift

Transfer a unique gift to another owner

### Synopsis

Transfer an owned unique gift to another chat (transferGift). The gift leaves this account. The receiving chat must have been active in the last 24 hours.

```
tgctl stars transfer-gift [flags]
```

### Examples

```
  tgctl stars transfer-gift --business-connection-id BQ... --gift-id abc123 --new-owner 12345
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --gift-id string                  id of the owned gift (from stars user-gifts)
  -h, --help                            help for transfer-gift
      --new-owner int                   chat id that will own the gift
      --star-count int                  stars to pay when the transfer is paid
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

