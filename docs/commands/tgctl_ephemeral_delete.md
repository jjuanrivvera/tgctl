## tgctl ephemeral delete

Delete an ephemeral message

### Synopsis

Delete an ephemeral message (deleteEphemeralMessage).
Returns True on success.

```
tgctl ephemeral delete [flags]
```

### Examples

```
  tgctl ephemeral delete --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42
```

### Options

```
      --chat string                target chat: numeric id or @username
      --ephemeral-message-id int   unique identifier of the ephemeral message
  -h, --help                       help for delete
      --receiver int               user id who received the ephemeral message
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

* [tgctl ephemeral](tgctl_ephemeral.md)	 - Manage and edit ephemeral messages

