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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl media animation](tgctl_media_animation.md)	 - Send an animation (GIF or H.264/MPEG-4 without sound)
* [tgctl media audio](tgctl_media_audio.md)	 - Send an audio file (shown in the music player)
* [tgctl media document](tgctl_media_document.md)	 - Send a document/file
* [tgctl media media-group](tgctl_media_media-group.md)	 - Send a group of photos/videos/documents as an album
* [tgctl media photo](tgctl_media_photo.md)	 - Send a photo
* [tgctl media sticker](tgctl_media_sticker.md)	 - Send a sticker (.WEBP, .TGS, or .WEBM)
* [tgctl media video](tgctl_media_video.md)	 - Send a video
* [tgctl media video-note](tgctl_media_video-note.md)	 - Send a video note (round video, up to 1 minute)
* [tgctl media voice](tgctl_media_voice.md)	 - Send a voice message (OGG/OPUS, shown as a waveform)

