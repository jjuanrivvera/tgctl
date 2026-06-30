## tgctl updates get

Get pending updates

```
tgctl updates get [flags]
```

### Examples

```
  tgctl updates get --limit 5
  tgctl updates get --offset 123456789 --timeout 30 -o json
  tgctl updates get --allowed-updates message,callback_query
```

### Options

```
      --allowed-updates strings   update types to receive
  -h, --help                      help for get
      --limit int                 max updates to return (1-100)
      --offset int                first update id to return (ack earlier ones)
      --timeout int               long-poll seconds (0 = short poll)
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

* [tgctl updates](tgctl_updates.md)	 - Fetch incoming updates (long polling)

