## tgctl webhook set

Set the webhook URL

```
tgctl webhook set [flags]
```

### Examples

```
  tgctl webhook set --url https://example.com/bot --max-connections 40
  tgctl webhook set --url https://example.com/bot --secret-token s3cr3t --drop-pending
```

### Options

```
      --allowed-updates strings   update types to receive
      --drop-pending              drop queued updates
  -h, --help                      help for set
      --max-connections int       max concurrent connections (1-100)
      --secret-token string       secret echoed in the X-Telegram-Bot-Api-Secret-Token header
      --url string                HTTPS URL to receive updates
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

* [tgctl webhook](tgctl_webhook.md)	 - Manage the bot's webhook

