## tgctl config

Inspect and edit tgctl configuration

### Synopsis

View the config file, switch profiles, and set per-profile options. Secrets live in the keyring and are never shown here.

### Options

```
  -h, --help   help for config
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
* [tgctl config list-profiles](tgctl_config_list-profiles.md)	 - List configured profiles
* [tgctl config path](tgctl_config_path.md)	 - Print the config file path
* [tgctl config set](tgctl_config_set.md)	 - Set a per-profile option (key: base_url)
* [tgctl config use](tgctl_config_use.md)	 - Switch the active profile
* [tgctl config view](tgctl_config_view.md)	 - Show the current configuration (secrets redacted)

