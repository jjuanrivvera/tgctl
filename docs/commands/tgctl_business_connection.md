## tgctl business connection

Get information about a business connection

### Synopsis

Get information about the connection of the bot with a business account.

```
tgctl business connection [flags]
```

### Examples

```
  tgctl business connection --business-connection-id bc123
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
  -h, --help                            help for connection
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

