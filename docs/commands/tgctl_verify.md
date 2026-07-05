## tgctl verify

Verify or unverify chats and users

### Synopsis

Add or remove the verification mark on a chat or user. Only available to bots acting on behalf of an organization that Telegram authorized to verify.

### Options

```
  -h, --help   help for verify
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
* [tgctl verify chat](tgctl_verify_chat.md)	 - Verify a chat on behalf of the bot's organization
* [tgctl verify remove-chat](tgctl_verify_remove-chat.md)	 - Remove verification from a chat
* [tgctl verify remove-user](tgctl_verify_remove-user.md)	 - Remove verification from a user
* [tgctl verify user](tgctl_verify_user.md)	 - Verify a user on behalf of the bot's organization

