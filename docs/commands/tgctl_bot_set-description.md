## tgctl bot set-description

Set the bot's description (shown in the empty chat)

```
tgctl bot set-description [flags]
```

### Options

```
      --description string     new description (0-512 chars)
  -h, --help                   help for set-description
      --language-code string   language this description applies to
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

