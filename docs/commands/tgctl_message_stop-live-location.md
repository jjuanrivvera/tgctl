## tgctl message stop-live-location

Stop updating a live location before its period expires

```
tgctl message stop-live-location [flags]
```

### Examples

```
  tgctl message stop-live-location --chat @me --message-id 42
```

### Options

```
      --business-connection-id string   business connection id on whose behalf to act
      --chat string                     target chat: numeric id or @username (with --message-id)
  -h, --help                            help for stop-live-location
      --inline-message-id string        inline message id (instead of --chat/--message-id)
      --message-id int                  message id (with --chat)
      --reply-markup string             inline/reply keyboard as JSON
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

