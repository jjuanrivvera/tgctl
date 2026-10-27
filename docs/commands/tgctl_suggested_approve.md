## tgctl suggested approve

Approve a suggested post

### Synopsis

Approve a suggested post message for publication in the corresponding channel. Requires the can_post_messages administrator right.

```
tgctl suggested approve [flags]
```

### Examples

```
  tgctl suggested approve --chat -100123456789 --message-id 42
  tgctl suggested approve --chat @channel --message-id 42 --send-date 1780000000
```

### Options

```
      --chat string      target chat: numeric id or @username
  -h, --help             help for approve
      --message-id int   message id
      --send-date int    point in time (Unix timestamp) when the post is expected to be published
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

* [tgctl suggested](tgctl_suggested.md)	 - Manage suggested posts in channel direct messages

