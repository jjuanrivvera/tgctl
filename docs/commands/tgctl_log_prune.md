## tgctl log prune

Delete recorded messages older than a duration

```
tgctl log prune [flags]
```

### Examples

```
  tgctl log prune --older-than 2160h   # 90 days
  tgctl log prune --older-than 720h    # 30 days
```

### Options

```
  -h, --help                help for prune
      --older-than string   delete messages recorded before now minus this Go duration (required)
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

