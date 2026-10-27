## tgctl business

Manage Telegram Business connections and accounts

### Synopsis

Inspect business connections, manage business account profiles, read/delete messages, and manage Stars and gifts.

### Options

```
  -h, --help   help for business
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
* [tgctl business connection](tgctl_business_connection.md)	 - Get information about a business connection
* [tgctl business delete-messages](tgctl_business_delete-messages.md)	 - Delete messages on behalf of a connected business account
* [tgctl business gifts](tgctl_business_gifts.md)	 - Get gifts received by a business account
* [tgctl business read-message](tgctl_business_read-message.md)	 - Mark a message from a connected business account as read
* [tgctl business remove-photo](tgctl_business_remove-photo.md)	 - Remove profile photo of a connected business account
* [tgctl business set-bio](tgctl_business_set-bio.md)	 - Set bio of a connected business account
* [tgctl business set-gift-settings](tgctl_business_set-gift-settings.md)	 - Set gift settings of a connected business account
* [tgctl business set-name](tgctl_business_set-name.md)	 - Set name of a connected business account
* [tgctl business set-username](tgctl_business_set-username.md)	 - Set username of a connected business account
* [tgctl business star-balance](tgctl_business_star-balance.md)	 - Get Star balance of a business account
* [tgctl business transfer-stars](tgctl_business_transfer-stars.md)	 - Transfer Stars from a business account to the bot

