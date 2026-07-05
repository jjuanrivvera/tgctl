## tgctl stars

Telegram Stars: transactions, gifts, and paid media

### Synopsis

Inspect the bot's Star balance transactions, send gifts, refund Star payments, manage Star subscriptions, and send paid media.

### Options

```
  -h, --help   help for stars
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl stars edit-subscription](tgctl_stars_edit-subscription.md)	 - Cancel or re-enable a user's Star subscription
* [tgctl stars gifts](tgctl_stars_gifts.md)	 - List the gifts the bot can send
* [tgctl stars refund](tgctl_stars_refund.md)	 - Refund a successful Star payment
* [tgctl stars send-gift](tgctl_stars_send-gift.md)	 - Send a gift to a user or channel
* [tgctl stars send-paid-media](tgctl_stars_send-paid-media.md)	 - Send paid media that recipients unlock with Stars
* [tgctl stars set-emoji-status](tgctl_stars_set-emoji-status.md)	 - Set a user's emoji status (requires the user's prior consent)
* [tgctl stars transactions](tgctl_stars_transactions.md)	 - List the bot's Star transactions

