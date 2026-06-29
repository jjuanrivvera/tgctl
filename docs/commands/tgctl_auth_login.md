## tgctl auth login

Store a bot token and verify it

### Synopsis

Capture a bot token (from @BotFather), verify it against getMe, and save it to the keyring for the active profile.

```
tgctl auth login [flags]
```

### Examples

```
  tgctl auth login                      # prompt for the token (hidden input)
  tgctl auth login --token 123:ABC...   # non-interactive
  TGCTL_PROFILE=staging tgctl auth login
```

### Options

```
  -h, --help           help for login
      --no-verify      skip the getMe verification call
      --token string   bot token (omit to be prompted with hidden input)
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

* [tgctl auth](tgctl_auth.md)	 - Manage bot tokens and verify authentication

