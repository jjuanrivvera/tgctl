## tgctl stars upgrade-gift

Upgrade a regular gift to a unique one

### Synopsis

Upgrade an owned regular gift to a unique gift (upgradeGift). Irreversible, and --star-count pays for it when the upgrade is not free. Needs the can_transfer_and_upgrade_gifts business right.

```
tgctl stars upgrade-gift [flags]
```

### Examples

```
  tgctl stars upgrade-gift --business-connection-id BQ... --gift-id abc123 --keep-original-details
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --gift-id string                  id of the owned gift (from stars user-gifts)
  -h, --help                            help for upgrade-gift
      --keep-original-details           keep the original text, sender and receiver
      --star-count int                  stars to pay when the upgrade is paid
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

