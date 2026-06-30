## tgctl invite edit

Edit an existing invite link

```
tgctl invite edit [flags]
```

### Examples

```
  tgctl invite edit --chat @group --invite-link https://t.me/+abc --member-limit 10
```

### Options

```
      --chat string            target chat: numeric id or @username
      --creates-join-request   users joining are placed in a join-request queue
      --expire-date int        unix time the link expires
  -h, --help                   help for edit
      --invite-link string     the invite link to act on
      --member-limit int       max users that may join via this link (1-99999)
      --name string            invite link name (0-32 chars)
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

* [tgctl invite](tgctl_invite.md)	 - Manage chat invite links

