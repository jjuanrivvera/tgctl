## tgctl webhook

Manage the bot's webhook

### Synopsis

Inspect, set, and delete the webhook used to receive updates over HTTPS (instead of polling).

### Options

```
  -h, --help   help for webhook
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
* [tgctl webhook delete](tgctl_webhook_delete.md)	 - Delete the webhook (switch back to polling)
* [tgctl webhook info](tgctl_webhook_info.md)	 - Show the current webhook status
* [tgctl webhook listen](tgctl_webhook_listen.md)	 - Run a local server that receives webhook updates and prints them
* [tgctl webhook set](tgctl_webhook_set.md)	 - Set the webhook URL

