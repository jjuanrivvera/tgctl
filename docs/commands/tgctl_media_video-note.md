## tgctl media video-note

Send a video note (round video, up to 1 minute)

```
tgctl media video-note [flags]
```

### Examples

```
  tgctl media video-note --chat @me --video-note ./round.mp4 --length 240
```

### Options

```
      --chat string         target chat: numeric id or @username
      --duration int        duration in seconds
  -h, --help                help for video-note
      --length int          video width and height (it is square)
      --silent              send without a notification sound
      --video-note string   local path or file_id (URLs not supported by Telegram)
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

