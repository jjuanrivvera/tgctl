## tgctl story repost

Repost a story across business accounts

### Synopsis

Repost a story across managed business accounts (repostStory).
Does not require uploading new media content. Returns the reposted Story.

```
tgctl story repost [flags]
```

### Examples

```
  tgctl story repost --business-connection-id dest_bc --from-business-connection-id src_bc --story-id 42 --active-period 86400
```

### Options

```
      --active-period int                    period after which story moves to archive in seconds (21600, 43200, 86400, or 172800)
      --areas string                         JSON array of StoryArea objects
      --business-connection-id string        id of the business connection acting on the account
      --caption string                       optional custom caption for the reposted story
      --caption-entities string              JSON array of MessageEntity objects
      --from-business-connection-id string   business connection id that originally posted the story
  -h, --help                                 help for repost
      --parse-mode string                    text formatting: MarkdownV2 | HTML | Markdown
      --post-to-chat-page                    keep story accessible on profile after expiration
      --protect-content                      protect the content from forwarding and saving
      --story-id int                         unique identifier of the story
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

