## tgctl media voice

Send a voice message (OGG/OPUS, shown as a waveform)

```
tgctl media voice [flags]
```

### Examples

```
  tgctl media voice --chat @me --voice ./note.ogg --duration 7
```

### Options

```
      --caption string      voice caption
      --chat string         target chat: numeric id or @username
      --duration int        duration in seconds
  -h, --help                help for voice
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --silent              send without a notification sound
      --voice string        local path, URL, or file_id
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

* [tgctl media](tgctl_media.md)	 - Send files: photos, documents, and video

