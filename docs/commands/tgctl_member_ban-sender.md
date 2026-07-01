## tgctl member ban-sender

Ban a channel from posting as itself in a chat

```
tgctl member ban-sender [flags]
```

### Examples

```
  tgctl member ban-sender --chat @group --sender-chat -1001234567890
```

### Options

```
      --chat string       target chat: numeric id or @username
  -h, --help              help for ban-sender
      --sender-chat int   id of the channel/chat to ban as a sender
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)

