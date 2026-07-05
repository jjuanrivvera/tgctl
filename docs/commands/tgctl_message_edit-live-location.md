## tgctl message edit-live-location

Update a live location message

```
tgctl message edit-live-location [flags]
```

### Examples

```
  tgctl message edit-live-location --chat @me --message-id 42 --latitude 3.46 --longitude -76.53
```

### Options

```
      --business-connection-id string   business connection id on whose behalf to act
      --chat string                     target chat: numeric id or @username (with --message-id)
      --heading int                     direction of movement in degrees (1-360)
  -h, --help                            help for edit-live-location
      --horizontal-accuracy float       location uncertainty radius in meters (0-1500)
      --inline-message-id string        inline message id (instead of --chat/--message-id)
      --latitude float                  new latitude
      --live-period int                 new live period in seconds
      --longitude float                 new longitude
      --message-id int                  message id (with --chat)
      --proximity-alert-radius int      max distance for proximity alerts in meters
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
      --no-store          disable local SQLite send/receive history for this invocation (see tgctl log)
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl message](tgctl_message.md)	 - Send and manage messages

