## tgctl stars user-gifts

List the gifts a user owns

### Synopsis

List the gifts owned and hosted by a user (getUserGifts). The --exclude-* flags
narrow the list by gift kind; --sort-by-price orders by value instead of by date.

```
tgctl stars user-gifts [flags]
```

### Examples

```
  tgctl stars user-gifts --user 12345
  tgctl stars user-gifts --user 12345 --exclude-unlimited --sort-by-price
```

### Options

```
      --exclude-from-blockchain          skip gifts held on the blockchain
      --exclude-limited-non-upgradable   skip limited gifts that cannot become unique
      --exclude-limited-upgradable       skip limited gifts that can become unique
      --exclude-unique                   skip unique gifts
      --exclude-unlimited                skip gifts sold in unlimited numbers
  -h, --help                             help for user-gifts
      --limit int                        max gifts to return (1-100)
      --offset string                    opaque cursor from a previous page
      --sort-by-price                    order by value instead of by date
      --user int                         owner of the gifts
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

