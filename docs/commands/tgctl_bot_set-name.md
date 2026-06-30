## tgctl bot set-name

Set the bot's name

```
tgctl bot set-name [flags]
```

### Examples

```
  tgctl bot set-name --name "My Helper Bot"
```

### Options

```
  -h, --help                   help for set-name
      --language-code string   BCP-47 code this name applies to
      --name string            new bot name (0-64 chars)
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

* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself

