## tgctl business set-photo

Set profile photo of a connected business account

### Synopsis

Set the profile photo of a connected business account (setBusinessAccountProfilePhoto).
The photo must be a local file. Pass --photo for a static image, or --animated for an animation (with --main-frame-timestamp).

```
tgctl business set-photo [flags]
```

### Examples

```
  tgctl business set-photo --business-connection-id bc123 --photo logo.jpg
  tgctl business set-photo --business-connection-id bc123 --animated intro.mp4 --public
```

### Options

```
      --animated string                 local MPEG4 to use as an animated profile photo
      --business-connection-id string   id of the business connection acting on the account
  -h, --help                            help for set-photo
      --main-frame-timestamp float      seconds into --animated for the still frame (default 0.0)
      --photo string                    local .JPG to use as the static profile photo
      --public                          set the public photo visible even if main photo is hidden by privacy settings
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

* [tgctl business](tgctl_business.md)	 - Manage Telegram Business connections and accounts

