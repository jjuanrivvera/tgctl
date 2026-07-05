## tgctl log

Query tgctl's local send/receive history

### Synopsis

tgctl records every outbound send (and, in polling/webhook mode, inbound updates)
to a local SQLite database — one per bot profile — because the Bot API itself has no history
endpoint. This lets a restarted or compacted session, or any external tool, answer "what did
you send/receive, when, to whom". Disable recording for a single call with --no-store; this
command itself always reads regardless of --no-store (it does not write).

```
tgctl log [flags]
```

### Examples

```
  tgctl log
  tgctl log --chat 123456789 --since 24h
  tgctl log --kind photo --limit 20 -o json
  tgctl log search "deploy failed"
  tgctl log show 42
  tgctl log prune --older-than 2160h
```

### Options

```
      --chat int       filter by chat id
  -h, --help           help for log
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl log prune](tgctl_log_prune.md)	 - Delete recorded messages older than a duration
* [tgctl log search](tgctl_log_search.md)	 - Full-text search recorded message/caption text
* [tgctl log show](tgctl_log_show.md)	 - Show one recorded message, including its full raw API payload

