## tgctl message unreact

Remove someone's reaction from a message

### Synopsis

Remove a reaction another member left on a message (deleteMessageReaction).

This is moderation, not the inverse of `message react`: it removes a reaction somebody
else added, so the bot needs the can_delete_messages right in the group. To clear the bot's
own reactions, call `message react` without --reaction.

Name whose reaction to remove with --user, or --actor-chat when it was left on behalf of a
channel.

```
tgctl message unreact [flags]
```

### Examples

```
  tgctl message unreact --chat @group --message-id 42 --user 12345
  tgctl message unreact --chat @group --message-id 42 --actor-chat -1001234567890
```

### Options

```
      --actor-chat int   id of the channel whose reaction to remove (instead of --user)
      --chat string      target chat: numeric id or @username
  -h, --help             help for unreact
      --message-id int   message id
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

