## tgctl inline answer-guest

Reply to a guest message with a single result

### Synopsis

Reply to a guest message (answerGuestQuery).

Unlike `inline answer`, which returns a LIST the user picks from, this sends ONE
result as the reply: --result is a single InlineQueryResult object, not an array.

```
tgctl inline answer-guest [flags]
```

### Examples

```
  tgctl inline answer-guest --query-id AAxx     --result '{"type":"article","id":"1","title":"Hi","input_message_content":{"message_text":"Hi"}}'
```

### Options

```
  -h, --help              help for answer-guest
      --query-id string   id of the guest query to answer
      --result string     a single InlineQueryResult object as JSON
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

* [tgctl inline](tgctl_inline.md)	 - Answer inline and guest queries

