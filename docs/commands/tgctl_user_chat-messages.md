## tgctl user chat-messages

Read the last messages on a user's personal chat

### Synopsis

Read the most recent messages from the chat a user has pinned to their profile
(getUserPersonalChatMessages). --limit is required by the API and capped at 20.

```
tgctl user chat-messages [flags]
```

### Examples

```
  tgctl user chat-messages --user 12345 --limit 5
  tgctl user chat-messages --user 12345 --limit 20 -o json
```

### Options

```
  -h, --help        help for chat-messages
      --limit int   how many messages to return (1-20)
      --user int    target user id
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

* [tgctl user](tgctl_user.md)	 - Read user information

