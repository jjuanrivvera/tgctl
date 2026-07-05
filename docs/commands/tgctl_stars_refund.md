## tgctl stars refund

Refund a successful Star payment

```
tgctl stars refund [flags]
```

### Examples

```
  tgctl stars refund --user 12345 --charge-id abc123
```

### Options

```
      --charge-id string   the telegram payment charge id to refund
  -h, --help               help for refund
      --user int           target user id
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

