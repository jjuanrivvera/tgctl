## tgctl message unreact-all

Remove every recent reaction one member left in a chat

### Synopsis

Remove up to 10000 recent reactions left in a group by one user or channel
(deleteAllMessageReactions). It sweeps the chat, not a single message, so it takes no
--message-id; the bot needs the can_delete_messages right.

```
tgctl message unreact-all [flags]
```

### Examples

```
  tgctl message unreact-all --chat @group --user 12345
  tgctl message unreact-all --chat @group --actor-chat -1001234567890
```

### Options

```
      --actor-chat int   id of the channel whose reaction to remove (instead of --user)
      --chat string      target chat: numeric id or @username
  -h, --help             help for unreact-all
      --user int         id of the user to act on (instead of --actor-chat)
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

