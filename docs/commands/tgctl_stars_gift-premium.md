## tgctl stars gift-premium

Gift a user a Telegram Premium subscription

### Synopsis

Gift Telegram Premium to a user (giftPremiumSubscription). It SPENDS the bot's
Stars and cannot be undone. The API fixes the price: 3 months costs 1000 stars, 6 months 1500,
and 12 months 2500 — any other pairing is rejected before the request leaves.

```
tgctl stars gift-premium [flags]
```

### Examples

```
  tgctl stars gift-premium --user 12345 --months 3 --stars 1000 --text "Happy birthday"
```

### Options

```
  -h, --help                     help for gift-premium
      --months int               3, 6 or 12
      --stars int                1000 for 3 months, 1500 for 6, 2500 for 12
      --text string              text shown with the service message
      --text-entities string     JSON array of MessageEntity (instead of --text-parse-mode)
      --text-parse-mode string   parse mode for --text
      --user int                 who receives the subscription
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

