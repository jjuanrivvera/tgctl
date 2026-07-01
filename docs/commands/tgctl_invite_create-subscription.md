## tgctl invite create-subscription

Create a subscription invite link (paid recurring access in Stars)

```
tgctl invite create-subscription [flags]
```

### Examples

```
  tgctl invite create-subscription --chat @channel --subscription-period 2592000 --subscription-price 100
```

### Options

```
      --chat string               target chat: numeric id or @username
  -h, --help                      help for create-subscription
      --name string               invite link name (0-32 chars)
      --subscription-period int   subscription period in seconds (currently must be 2592000 = 30 days)
      --subscription-price int    price in Telegram Stars per period (1-2500)
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

* [tgctl invite](tgctl_invite.md)	 - Manage chat invite links

