## tgctl log search

Full-text search recorded message/caption text

### Synopsis

Search uses FTS5 MATCH when the linked SQLite build supports it (operators: AND/OR/
NOT, prefix*, "phrases"); otherwise it degrades to a plain substring scan automatically — check
"tgctl doctor" or the store's FTSEnabled to see which mode is active.

```
tgctl log search <query> [flags]
```

### Examples

```
  tgctl log search "deploy failed"
  tgctl log search "deploy* AND staging" --chat 123456789
```

### Options

```
      --chat int       filter by chat id
  -h, --help           help for search
      --kind string    filter by kind: text|photo|document|voice|edit|...
      --limit int      max rows to return (default 50)
      --since string   only messages at/after this time: a Go duration (24h) or RFC3339/YYYY-MM-DD
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

* [tgctl log](tgctl_log.md)	 - Query tgctl's local send/receive history

