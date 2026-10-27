## tgctl inline prepare-button

Prepare a keyboard button for a Mini App user

### Synopsis

Store a KeyboardButton with exactly one of request_users, request_chat or request_managed_bot. Returns its identifier.

```
tgctl inline prepare-button [flags]
```

### Options

```
      --button string   KeyboardButton object as JSON
  -h, --help            help for prepare-button
      --user int        target user id
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

