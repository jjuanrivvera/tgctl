## tgctl forum reopen-general

Reopen the General forum topic (also unhides it)

```
tgctl forum reopen-general [flags]
```

### Examples

```
  tgctl forum reopen-general --chat @group
```

### Options

```
      --chat string   target chat: numeric id or @username
  -h, --help          help for reopen-general
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl forum](tgctl_forum.md)	 - Manage forum topics in supergroups

