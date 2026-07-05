## tgctl message react

Set (or clear) reactions on a message

### Synopsis

Set reactions on a message. --reaction is a JSON array of ReactionType objects; omit it to remove all reactions.

```
tgctl message react [flags]
```

### Examples

```
  tgctl message react --chat @group --message-id 42 --reaction '[{"type":"emoji","emoji":"👍"}]'
  tgctl message react --chat @group --message-id 42   # clear reactions
```

### Options

```
      --chat string       target chat: numeric id or @username
  -h, --help              help for react
      --is-big            show the reaction with a big animation
      --message-id int    message id
      --reaction string   JSON array of ReactionType objects (omit to clear)
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

