## tgctl member promote

Promote or demote an administrator

```
tgctl member promote [flags]
```

### Examples

```
  tgctl member promote --chat @group --user 12345 --can-delete-messages --can-pin-messages
```

### Options

```
      --can-change-info        can change chat title/photo
      --can-delete-messages    can delete others' messages
      --can-invite-users       can invite new users
      --can-manage-chat        can access the admin log, etc.
      --can-pin-messages       can pin messages
      --can-promote-members    can add new admins
      --can-restrict-members   can restrict/ban members
      --chat string            target chat: numeric id or @username
  -h, --help                   help for promote
      --user int               target user id
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

* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)

