## tgctl member join-web-app

Show a Mini App for a join request query

### Synopsis

Open an HTTPS Mini App for a received join request query; use member answer-join-query to resolve it afterwards.

```
tgctl member join-web-app [flags]
```

### Options

```
  -h, --help                 help for join-web-app
      --query-id string      join request query identifier
      --web-app-url string   HTTPS URL of the Mini App
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

