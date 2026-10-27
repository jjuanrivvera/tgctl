## tgctl inline answer-web-app

Answer a Web App query with a message

### Synopsis

Send one InlineQueryResult on behalf of the user to the originating chat. Returns a SentWebAppMessage.

```
tgctl inline answer-web-app [flags]
```

### Options

```
  -h, --help              help for answer-web-app
      --query-id string   Web App query identifier
      --result string     single InlineQueryResult object as JSON
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

