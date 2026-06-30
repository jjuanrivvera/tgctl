## tgctl media audio

Send an audio file (shown in the music player)

```
tgctl media audio [flags]
```

### Examples

```
  tgctl media audio --chat @me --audio ./song.mp3 --performer "Artist" --title "Track"
```

### Options

```
      --audio string        local path, URL, or file_id
      --caption string      audio caption
      --chat string         target chat: numeric id or @username
      --duration int        duration in seconds
  -h, --help                help for audio
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --performer string    performer name
      --silent              send without a notification sound
      --title string        track title
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

