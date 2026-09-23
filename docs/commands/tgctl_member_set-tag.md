## tgctl member set-tag

Set (or clear) a regular member's tag

### Synopsis

Set the tag shown next to a regular member's name in a group (setChatMemberTag).
The bot must be an administrator with the can_manage_tags right. Omit --tag to clear it.

```
tgctl member set-tag [flags]
```

### Examples

```
  tgctl member set-tag --chat @group --user 12345 --tag "Moderator"
  tgctl member set-tag --chat @group --user 12345   # clear the tag
```

### Options

```
      --chat string   target chat: numeric id or @username
  -h, --help          help for set-tag
      --tag string    new tag (0-16 chars, no emoji; omit to clear)
      --user int      target user id
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

