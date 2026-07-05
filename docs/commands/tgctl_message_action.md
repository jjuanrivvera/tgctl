## tgctl message action

Show a chat action (typing, uploading, …) for a few seconds

### Synopsis

Tell the user something is happening on the bot's side. The status is cleared when a message arrives or after ~5 seconds.

```
tgctl message action [flags]
```

### Examples

```
  tgctl message action --chat @me --action typing
  tgctl message action --chat @me --action upload_photo
```

### Options

```
      --action string                   typing | upload_photo | record_video | upload_video | record_voice | upload_voice | upload_document | choose_sticker | find_location | record_video_note | upload_video_note
      --business-connection-id string   business connection id on whose behalf to act
      --chat string                     target chat: numeric id or @username
  -h, --help                            help for action
      --thread int                      forum topic thread id
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

