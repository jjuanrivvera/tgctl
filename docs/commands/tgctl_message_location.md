## tgctl message location

Send a point on the map

```
tgctl message location [flags]
```

### Examples

```
  tgctl message location --chat @me --latitude 3.4516 --longitude -76.532
```

### Options

```
      --chat string       target chat: numeric id or @username
  -h, --help              help for location
      --latitude float    latitude of the location
      --live-period int   seconds the location is updated live (60-86400)
      --longitude float   longitude of the location
      --silent            send without a notification sound
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

