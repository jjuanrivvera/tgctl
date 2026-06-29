## tgctl member restrict

Restrict what a member can do

```
tgctl member restrict [flags]
```

### Examples

```
  tgctl member restrict --chat @group --user 12345 \
    --permissions '{"can_send_messages":false}'
```

### Options

```
      --chat string          target chat: numeric id or @username
  -h, --help                 help for restrict
      --permissions string   ChatPermissions object as JSON
      --until int            unix time the restriction lifts
      --user int             target user id
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

* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)

