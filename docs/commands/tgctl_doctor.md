## tgctl doctor

Diagnose configuration, credentials, and connectivity

### Synopsis

Run a series of checks (config, token, API reachability, clock) and exit non-zero if any fail.

```
tgctl doctor [flags]
```

### Examples

```
  tgctl doctor
  tgctl doctor --json
```

### Options

```
  -h, --help   help for doctor
      --json   output checks as JSON
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API

