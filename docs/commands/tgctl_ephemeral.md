## tgctl ephemeral

Manage and edit ephemeral messages

### Synopsis

Delete and edit ephemeral messages (visible only to a target user and the bot).

Ephemeral messages are sent in group or channel contexts to a specific user.
Their identifiers (ephemeral_message_id) are obtained from incoming updates or callback queries.
Targeting requires the chat, the recipient user id, and the ephemeral message id.
Delivery of edits and deletions is client-side and not guaranteed if the recipient is offline.

### Options

```
  -h, --help   help for ephemeral
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl ephemeral delete](tgctl_ephemeral_delete.md)	 - Delete an ephemeral message
* [tgctl ephemeral edit-caption](tgctl_ephemeral_edit-caption.md)	 - Edit caption of an ephemeral media message
* [tgctl ephemeral edit-media](tgctl_ephemeral_edit-media.md)	 - Edit media content of an ephemeral message
* [tgctl ephemeral edit-reply-markup](tgctl_ephemeral_edit-reply-markup.md)	 - Edit reply markup of an ephemeral message
* [tgctl ephemeral edit-text](tgctl_ephemeral_edit-text.md)	 - Edit text or rich content of an ephemeral message

