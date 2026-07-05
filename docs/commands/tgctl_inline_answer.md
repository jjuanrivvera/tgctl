## tgctl inline answer

Answer an inline query with results

### Synopsis

Answer an inline query. --results is a JSON array of InlineQueryResult objects.

```
tgctl inline answer [flags]
```

### Examples

```
  tgctl inline answer --inline-query-id 999 \
    --results '[{"type":"article","id":"1","title":"Hi","input_message_content":{"message_text":"Hi"}}]'
```

### Options

```
      --button string            InlineQueryResultsButton object as JSON
      --cache-time int           seconds the result may be cached server-side
  -h, --help                     help for answer
      --inline-query-id string   id of the inline query to answer
      --is-personal              cache results per-user instead of globally
      --next-offset string       offset a client sends to request the next page
      --results string           JSON array of InlineQueryResult objects (max 50)
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

* [tgctl inline](tgctl_inline.md)	 - Answer inline queries

