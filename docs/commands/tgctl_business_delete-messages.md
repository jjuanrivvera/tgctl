## tgctl business delete-messages

Delete messages on behalf of a connected business account

### Synopsis

Delete 1 to 100 messages on behalf of a connected business account.

```
tgctl business delete-messages [flags]
```

### Examples

```
  tgctl business delete-messages --business-connection-id bc123 --message-ids '[10,11,12]'
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
  -h, --help                            help for delete-messages
      --message-ids string              JSON array of 1-100 message identifiers to delete
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

