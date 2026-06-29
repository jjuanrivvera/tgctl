## tgctl alias

Manage user-defined command aliases

### Synopsis

Define shorthand commands. Aliases are expanded before parsing and can never shadow a built-in.

### Options

```
  -h, --help   help for alias
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
* [tgctl alias list](tgctl_alias_list.md)	 - List aliases
* [tgctl alias remove](tgctl_alias_remove.md)	 - Remove an alias
* [tgctl alias set](tgctl_alias_set.md)	 - Create or update an alias

