## tgctl story edit

Edit a story posted on behalf of a business account

### Synopsis

Edit a story previously posted on behalf of a connected business account (editStory).
--content accepts an InputStoryContent JSON object pointing at uploaded files via attach://.
Alternatively, pass --photo or --video with local files. Returns the edited Story.

```
tgctl story edit [flags]
```

### Examples

```
  tgctl story edit --business-connection-id bc123 --story-id 42 --caption "Updated caption"
  tgctl story edit --business-connection-id bc123 --story-id 42 --photo ./new.jpg
```

### Options

```
      --areas string                    new JSON array of StoryArea objects
      --business-connection-id string   id of the business connection acting on the account
      --caption string                  new story caption
      --caption-entities string         JSON array of MessageEntity objects
      --content string                  JSON object describing new story content
      --cover string                    optional local image file as video cover
  -F, --file stringArray                name=@path to upload as multipart (repeatable)
  -h, --help                            help for edit
      --parse-mode string               text formatting: MarkdownV2 | HTML | Markdown
      --photo string                    local image file for a photo story
      --story-id int                    unique identifier of the story
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

