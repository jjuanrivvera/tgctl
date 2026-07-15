## tgctl update

Update tgctl to the latest GitHub release

### Synopsis

Download the latest tgctl release, verify it against checksums.txt, and replace
the running binary in place. Use 'tgctl update check' to see what's available without
installing.

```
tgctl update [flags]
```

### Examples

```
  tgctl update
  tgctl update check
```

### Options

```
  -h, --help   help for update
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
* [tgctl update check](tgctl_update_check.md)	 - Check for a newer release without installing it

