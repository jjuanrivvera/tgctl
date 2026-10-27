## tgctl story

Manage and publish stories on behalf of business accounts

### Synopsis

Post, edit, repost, and delete stories on behalf of connected business accounts.

Stories require the can_manage_stories business bot right.
New story content requires uploading a local photo or video (Telegram accepts no URLs or existing file_ids here).
Stories move to the archive after the specified active period (6h, 12h, 24h, or 48h).

### Options

```
  -h, --help   help for story
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl story delete](tgctl_story_delete.md)	 - Delete a story posted on behalf of a business account
* [tgctl story edit](tgctl_story_edit.md)	 - Edit a story posted on behalf of a business account
* [tgctl story post](tgctl_story_post.md)	 - Post a story on behalf of a business account
* [tgctl story repost](tgctl_story_repost.md)	 - Repost a story across business accounts

