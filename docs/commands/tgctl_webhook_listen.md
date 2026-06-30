## tgctl webhook listen

Run a local server that receives webhook updates and prints them

### Synopsis

Start an HTTP server that receives Telegram webhook updates (the Update objects
Telegram POSTs to your webhook URL) and renders each one with the usual -o/--output.

Telegram only delivers to a public HTTPS endpoint, so put this behind a tunnel
(cloudflared, ngrok) or run it on a public host. Point the webhook at that URL either
yourself (tgctl webhook set --url ...) or with --set-url here. When --secret-token is set,
requests whose X-Telegram-Bot-Api-Secret-Token header doesn't match are rejected. Ctrl-C
stops the server (and deletes the webhook if --delete-on-exit).

```
tgctl webhook listen [flags]
```

### Examples

```
  tgctl webhook listen --port 8080 -o json
  tgctl webhook listen --port 8080 --secret-token s3cr3t \
      --set-url https://my-tunnel.example.com/bot --delete-on-exit
  # local test (no bot needed):
  curl -XPOST localhost:8080 -d '{"update_id":1,"message":{"message_id":7,"text":"hi"}}'
```

### Options

```
      --delete-on-exit        delete the webhook when the server stops (requires --set-url)
  -h, --help                  help for listen
      --path string           URL path to receive updates on (default "/")
      --port int              port to listen on (default 8080)
      --secret-token string   require this X-Telegram-Bot-Api-Secret-Token header
      --set-url string        register the webhook at this public HTTPS URL before listening
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

* [tgctl webhook](tgctl_webhook.md)	 - Manage the bot's webhook

