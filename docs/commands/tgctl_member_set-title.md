## tgctl member set-title

Set a custom title for an administrator the bot promoted

```
tgctl member set-title [flags]
```

### Examples

```
  tgctl member set-title --chat @group --user 12345 --title "Community Lead"
```

### Options

```
      --chat string    target chat: numeric id or @username
  -h, --help           help for set-title
      --title string   custom admin title (0-16 chars, no emoji)
      --user int       target user id
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

* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)

