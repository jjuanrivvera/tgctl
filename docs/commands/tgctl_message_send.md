## tgctl message send

Send a text message

```
tgctl message send [flags]
```

### Examples

```
  tgctl message send --chat @me --text "hello"
  tgctl message send --chat -1001234567890 --text "*bold*" --parse-mode MarkdownV2 --silent
```

### Options

```
      --chat string           target chat: numeric id or @username
  -h, --help                  help for send
      --no-preview            disable link previews
      --parse-mode string     text formatting: MarkdownV2 | HTML | Markdown
      --reply-markup string   inline/reply keyboard as JSON
      --reply-to int          message id to reply to
      --silent                send without a notification sound
      --text string           message text
      --thread int            forum topic thread id
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

