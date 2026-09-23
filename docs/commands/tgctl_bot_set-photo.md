## tgctl bot set-photo

Set the bot's profile photo

### Synopsis

Set the bot's own profile photo (setMyProfilePhoto).

The photo must be a local file: Telegram does not accept a URL or an existing file_id
here, because a profile photo cannot be reused. Pass --photo for a static .JPG, or
--animated for an MPEG4 animation (with --main-frame-timestamp to choose the frame
Telegram shows as the still image).

```
tgctl bot set-photo [flags]
```

### Examples

```
  tgctl bot set-photo --photo logo.jpg
  tgctl bot set-photo --animated intro.mp4 --main-frame-timestamp 1.5
```

### Options

```
      --animated string              local MPEG4 to use as an animated profile photo
  -h, --help                         help for set-photo
      --main-frame-timestamp float   seconds into --animated for the still frame (default 0.0)
      --photo string                 local .JPG to use as the static profile photo
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

* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself

