## tgctl callback answer

Answer a callback query (toast, alert, or URL)

```
tgctl callback answer [flags]
```

### Examples

```
  tgctl callback answer --callback-query-id 12345 --text "Saved!"
  tgctl callback answer --callback-query-id 12345 --text "Not allowed" --show-alert
```

### Options

```
      --cache-time int             seconds the result may be cached client-side
      --callback-query-id string   id of the callback query to answer
  -h, --help                       help for answer
      --show-alert                 show an alert dialog instead of a toast
      --text string                notification text shown to the user (0-200 chars)
      --url string                 URL the client opens (game / t.me deep link)
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

* [tgctl callback](tgctl_callback.md)	 - Answer callback queries from inline keyboards

