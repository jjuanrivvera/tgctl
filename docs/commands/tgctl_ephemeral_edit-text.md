## tgctl ephemeral edit-text

Edit text or rich content of an ephemeral message

### Synopsis

Edit the text or structured rich message of an ephemeral message (editEphemeralMessageText).
Either --text or --rich-message must be provided. Returns True on success.

```
tgctl ephemeral edit-text [flags]
```

### Examples

```
  tgctl ephemeral edit-text --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --text "Updated text"
  tgctl ephemeral edit-text --chat -1001234567890 --receiver 12345 --ephemeral-message-id 42 --rich-message '{"blocks":[{"type":"paragraph","text":"New content"}]}'
```

### Options

```
      --chat string                   target chat: numeric id or @username
      --entities string               JSON array of MessageEntity objects
      --ephemeral-message-id int      unique identifier of the ephemeral message
  -h, --help                          help for edit-text
      --link-preview-options string   JSON object for link preview options
      --parse-mode string             text formatting: MarkdownV2 | HTML | Markdown
      --receiver int                  user id who received the ephemeral message
      --reply-markup string           inline/reply keyboard as JSON
      --rich-message string           JSON object describing updated rich message (InputRichMessage)
      --text string                   new text of the ephemeral message
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

