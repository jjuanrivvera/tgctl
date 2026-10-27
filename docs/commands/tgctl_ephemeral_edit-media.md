## tgctl ephemeral edit-media

Edit media content of an ephemeral message

### Synopsis

Edit animation, audio, document, photo, or video of an ephemeral message (editEphemeralMessageMedia).
--media accepts an InputMedia JSON object. New file uploads can be referenced via attach:// and passed with --file. Returns True on success.

```
tgctl ephemeral edit-media [flags]
```

### Examples

```
  tgctl ephemeral edit-media --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --media '{"type":"photo","media":"https://example.com/pic.jpg"}'
```

### Options

```
      --chat string                target chat: numeric id or @username
      --ephemeral-message-id int   unique identifier of the ephemeral message
  -F, --file stringArray           name=@path to upload as multipart (repeatable)
  -h, --help                       help for edit-media
      --media string               JSON object describing new media content (InputMedia)
      --receiver int               user id who received the ephemeral message
      --reply-markup string        inline/reply keyboard as JSON
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

