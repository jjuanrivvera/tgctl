## tgctl suggested decline

Decline a suggested post

### Synopsis

Decline a suggested post message with an optional feedback comment. Requires the can_manage_direct_messages administrator right.

```
tgctl suggested decline [flags]
```

### Examples

```
  tgctl suggested decline --chat -100123456789 --message-id 42
  tgctl suggested decline --chat -100123456789 --message-id 42 --comment "Off-topic"
```

### Options

```
      --chat string      target chat: numeric id or @username
      --comment string   optional feedback comment for the post author (0-128 chars)
  -h, --help             help for decline
      --message-id int   message id
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

