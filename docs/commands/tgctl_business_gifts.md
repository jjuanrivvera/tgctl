## tgctl business gifts

Get gifts received by a business account

### Synopsis

Get the list of gifts received by the business account.

```
tgctl business gifts [flags]
```

### Examples

```
  tgctl business gifts --business-connection-id bc123
  tgctl business gifts --business-connection-id bc123 --exclude-unlimited --limit 50
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --exclude-limited                 exclude limited gifts
      --exclude-saved                   exclude gifts that are saved to the profile
      --exclude-unique                  exclude unique gifts
      --exclude-unlimited               exclude unlimited gifts
      --exclude-unsaved                 exclude gifts that are not saved to the profile
  -h, --help                            help for gifts
      --limit int                       number of gifts to return (1-100)
      --offset string                   offset of the first gift to return from previous request
      --sort-by-price                   sort results by gift price instead of send date
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

* [tgctl business](tgctl_business.md)	 - Manage Telegram Business connections and accounts

