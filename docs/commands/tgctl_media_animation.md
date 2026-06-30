## tgctl media animation

Send an animation (GIF or H.264/MPEG-4 without sound)

```
tgctl media animation [flags]
```

### Examples

```
  tgctl media animation --chat @me --animation ./loop.gif --caption "nice"
```

### Options

```
      --animation string    local path, URL, or file_id
      --caption string      animation caption
      --chat string         target chat: numeric id or @username
      --duration int        duration in seconds
  -h, --help                help for animation
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --silent              send without a notification sound
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

