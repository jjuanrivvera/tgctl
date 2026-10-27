## tgctl business set-username

Set username of a connected business account

### Synopsis

Change the username of the business account. An empty username removes it (5-32 characters).

```
tgctl business set-username [flags]
```

### Examples

```
  tgctl business set-username --business-connection-id bc123 --username acmeshop
  tgctl business set-username --business-connection-id bc123 --username ""
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
  -h, --help                            help for set-username
      --username string                 new business account username (5-32 characters); empty to remove
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

