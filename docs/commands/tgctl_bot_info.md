## tgctl bot info

Show the authenticated bot's identity (getMe)

```
tgctl bot info [flags]
```

### Examples

```
  tgctl bot info
  tgctl bot info -o json
```

### Options

```
  -h, --help   help for info
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --profile string    profile/instance to use (env TGCTL_PROFILE)
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself

