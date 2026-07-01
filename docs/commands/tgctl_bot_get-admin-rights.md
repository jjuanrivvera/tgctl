## tgctl bot get-admin-rights

Get the bot's default administrator rights

```
tgctl bot get-admin-rights [flags]
```

### Examples

```
  tgctl bot get-admin-rights
  tgctl bot get-admin-rights --for-channels -o json
```

### Options

```
      --for-channels   query the channel rights instead of group rights
  -h, --help           help for get-admin-rights
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

