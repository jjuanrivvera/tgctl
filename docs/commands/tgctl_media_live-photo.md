## tgctl media live-photo

Send a live photo

### Synopsis

Send a live photo (a still photo accompanied by a short video).
--photo and --live-photo both accept a local file path or an existing file_id.
HTTP(S) URLs are not supported by the Bot API for this method.

```
tgctl media live-photo [flags]
```

### Examples

```
  tgctl media live-photo --chat @me --photo ./still.jpg --live-photo ./motion.mp4 --caption "look at this"
  tgctl media live-photo --chat @me --photo AgACAgI... --live-photo BAACAgI...
```

### Options

```
      --allow-paid-broadcast            allow up to 1000 messages per second by paying Telegram Stars
      --business-connection-id string   unique identifier of the business connection
      --caption string                  live photo caption
      --caption-entities string         JSON array of special entities that appear in the caption
      --chat string                     target chat: numeric id or @username
      --has-spoiler                     cover the live photo with a spoiler animation
  -h, --help                            help for live-photo
      --live-photo string               live photo video part (local path or file_id; URLs not supported)
      --message-effect-id string        unique identifier of the message effect
      --parse-mode string               text formatting: MarkdownV2 | HTML | Markdown
      --photo string                    still photo part (local path or file_id; URLs not supported)
      --protect-content                 protect the content from forwarding and saving
      --reply-markup string             inline/reply keyboard as JSON
      --reply-parameters string         JSON object describing message to reply to
      --show-caption-above-media        display caption above media
      --silent                          send without a notification sound
      --thread int                      forum topic thread id
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

