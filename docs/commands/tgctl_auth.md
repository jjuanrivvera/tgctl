## tgctl auth

Manage bot tokens and verify authentication

### Synopsis

Capture, verify, and remove the bot token for a profile. Tokens are stored in your OS keyring, never in the config file.

### Options

```
  -h, --help   help for auth
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
* [tgctl auth login](tgctl_auth_login.md)	 - Store a bot token and verify it
* [tgctl auth logout](tgctl_auth_logout.md)	 - Remove the stored token for the active profile
* [tgctl auth status](tgctl_auth_status.md)	 - Show the active profile, base URL, and token validity

