## tgctl user photos

List a user's profile photos

```
tgctl user photos [flags]
```

### Examples

```
  tgctl user photos --user 12345
  tgctl user photos --user 12345 --limit 1 -o json
```

### Options

```
  -h, --help         help for photos
      --limit int    max photos to return (1-100)
      --offset int   number of photos to skip
      --user int     target user id
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

* [tgctl user](tgctl_user.md)	 - Read user information

