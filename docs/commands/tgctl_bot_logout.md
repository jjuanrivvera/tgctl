## tgctl bot logout

Log out from the cloud Bot API before running a local Bot API server

### Synopsis

Log the bot out of the cloud Bot API. After this you can use a local Bot API server; you must re-login via api.telegram.org to switch back. Returns an error for the first 10 minutes after launch.

```
tgctl bot logout [flags]
```

### Examples

```
  tgctl bot logout
```

### Options

```
  -h, --help   help for logout
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

* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself

