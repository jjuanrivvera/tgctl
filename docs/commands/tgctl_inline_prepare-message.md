## tgctl inline prepare-message

Prepare a message for a Mini App user to send

### Synopsis

Store one InlineQueryResult for the user. Returns an identifier and expiration date.

```
tgctl inline prepare-message [flags]
```

### Options

```
      --allow-bot-chats       allow private chats with bots
      --allow-channel-chats   allow channel chats
      --allow-group-chats     allow group and supergroup chats
      --allow-user-chats      allow private chats with users
  -h, --help                  help for prepare-message
      --result string         single InlineQueryResult object as JSON
      --user int              target user id
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

