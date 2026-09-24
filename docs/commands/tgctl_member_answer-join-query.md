## tgctl member answer-join-query

Answer a join-request query from a Mini App flow

### Synopsis

Resolve a chat-join-request query (answerChatJoinRequestQuery).

This is the Mini App path, not the plain one: it answers a QUERY the bot received, identified
by its id, where `member approve-join` / `decline-join` act on a chat and a user
directly. --result takes approve, decline, or queue to leave the decision to another admin.

```
tgctl member answer-join-query [flags]
```

### Examples

```
  tgctl member answer-join-query --query-id AAxx... --result approve
  tgctl member answer-join-query --query-id AAxx... --result queue
```

### Options

```
  -h, --help              help for answer-join-query
      --query-id string   id of the join-request query being answered
      --result string     approve | decline | queue
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

* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)

