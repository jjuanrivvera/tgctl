## tgctl business read-message

Mark a message from a connected business account as read

### Synopsis

Mark a message from a connected business account as read on behalf of the business account.

```
tgctl business read-message [flags]
```

### Examples

```
  tgctl business read-message --business-connection-id bc123 --chat 12345678 --message-id 42
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --chat string                     target chat: numeric id or @username
  -h, --help                            help for read-message
      --message-id int                  message id
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

