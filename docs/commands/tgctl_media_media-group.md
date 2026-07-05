## tgctl media media-group

Send a group of photos/videos/documents as an album

### Synopsis

Send 2-10 items as a single album. --media is a JSON array of InputMedia objects;
each item's "media" is an http(s) URL or an existing file_id (multipart attach:// uploads
are not supported here — upload first with the single-item commands if you need local files).

```
tgctl media media-group [flags]
```

### Examples

```
  tgctl media media-group --chat @me \
    --media '[{"type":"photo","media":"https://e.com/a.jpg"},{"type":"photo","media":"https://e.com/b.jpg"}]'
```

### Options

```
      --chat string    target chat: numeric id or @username
  -h, --help           help for media-group
      --media string   JSON array of InputMedia objects
      --silent         send without a notification sound
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

