## tgctl business set-name

Set name of a connected business account

### Synopsis

Change the first and last name of the business account.

```
tgctl business set-name [flags]
```

### Examples

```
  tgctl business set-name --business-connection-id bc123 --first-name "Acme Corp"
  tgctl business set-name --business-connection-id bc123 --first-name "Acme" --last-name "HQ"
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --first-name string               new first name of the business account (1-64 characters)
  -h, --help                            help for set-name
      --last-name string                new last name of the business account (0-64 characters)
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

