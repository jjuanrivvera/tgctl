## tgctl story post

Post a story on behalf of a business account

### Synopsis

Post a story on behalf of a connected business account (postStory).
--content accepts an InputStoryContent JSON object pointing at uploaded files via attach://.
Alternatively, pass --photo or --video with local files. Returns the posted Story.

```
tgctl story post [flags]
```

### Examples

```
  tgctl story post --business-connection-id bc123 --photo ./story.jpg --active-period 86400 --caption "Daily update"
  tgctl story post --business-connection-id bc123 --video ./story.mp4 --active-period 43200
  tgctl story post --business-connection-id bc123 --content '{"type":"photo","photo":"attach://pic"}' --file pic=@./story.jpg --active-period 86400
```

### Options

```
      --active-period int               period after which story moves to archive in seconds (21600, 43200, 86400, or 172800)
      --areas string                    JSON array of StoryArea clickable interactive overlay objects
      --business-connection-id string   id of the business connection acting on the account
      --caption string                  story caption (0-2048 chars)
      --caption-entities string         JSON array of MessageEntity objects
      --content string                  JSON object describing story content (InputStoryContentPhoto or InputStoryContentVideo)
      --cover string                    optional local image file as video cover
  -F, --file stringArray                name=@path to upload as multipart (repeatable)
  -h, --help                            help for post
      --parse-mode string               text formatting: MarkdownV2 | HTML | Markdown
      --photo string                    local image file for a photo story
      --post-to-chat-page               keep story accessible on profile after expiration
      --protect-content                 protect the content from forwarding and saving
      --video string                    local video file for a video story
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

* [tgctl story](tgctl_story.md)	 - Manage and publish stories on behalf of business accounts

