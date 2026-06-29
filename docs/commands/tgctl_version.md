## tgctl version

Print version, commit, and build date

### Synopsis

Print build metadata. With --check, compare against the latest GitHub release.

```
tgctl version [flags]
```

### Examples

```
  tgctl version
  tgctl version --json
  tgctl version --check
```

### Options

```
      --check   check for a newer release on GitHub
  -h, --help    help for version
      --json    output as JSON
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

