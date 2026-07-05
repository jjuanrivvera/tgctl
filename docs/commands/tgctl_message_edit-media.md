## tgctl message edit-media

Replace the media of a message

### Synopsis

Replace a message's photo/video/animation/audio/document. --media is an InputMedia object as JSON.

```
tgctl message edit-media [flags]
```

### Examples

```
  tgctl message edit-media --chat @me --message-id 42 \
    --media '{"type":"photo","media":"https://e.com/new.jpg"}'
```

### Options

```
      --business-connection-id string   business connection id on whose behalf to act
      --chat string                     target chat: numeric id or @username (with --message-id)
  -h, --help                            help for edit-media
      --inline-message-id string        inline message id (instead of --chat/--message-id)
      --media string                    InputMedia object as JSON
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
      --no-store          disable local SQLite send/receive history for this invocation (see tgctl log)
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl message](tgctl_message.md)	 - Send and manage messages

