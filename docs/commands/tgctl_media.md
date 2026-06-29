## tgctl media

Send files: photos, documents, and video

### Synopsis

Send media to a chat. Each --<kind> flag accepts a local file path (uploaded as
multipart/form-data), an http(s) URL (Telegram fetches it), or an existing file_id.

### Options

```
  -h, --help   help for media
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --profile string    profile/instance to use (env TGCTL_PROFILE)
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl media document](tgctl_media_document.md)	 - Send a document/file
* [tgctl media photo](tgctl_media_photo.md)	 - Send a photo
* [tgctl media video](tgctl_media_video.md)	 - Send a video

