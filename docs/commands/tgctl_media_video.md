## tgctl media video

Send a video

```
tgctl media video [flags]
```

### Examples

```
  tgctl media video --chat @me --video ./clip.mp4
```

### Options

```
      --caption string      video caption
      --chat string         target chat: numeric id or @username
  -h, --help                help for video
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --silent              send without a notification sound
      --video string        local path, URL, or file_id
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

* [tgctl media](tgctl_media.md)	 - Send files: photos, documents, and video

