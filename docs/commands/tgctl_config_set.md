## tgctl config set

Set a per-profile option (key: base_url)

### Synopsis

Set a non-secret option on the active profile. Supported keys: base_url.

```
tgctl config set <key> <value> [flags]
```

### Examples

```
  tgctl config set base_url https://api.telegram.org
  tgctl --profile staging config set base_url http://localhost:8081
```

### Options

```
  -h, --help   help for set
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

* [tgctl config](tgctl_config.md)	 - Inspect and edit tgctl configuration

