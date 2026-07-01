## tgctl member

Moderate chat members (ban, restrict, promote)

### Synopsis

Administrative actions on members. The bot must be an admin with the relevant rights.

### Options

```
  -h, --help   help for member
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl member approve-join](tgctl_member_approve-join.md)	 - Approve a chat join request
* [tgctl member ban](tgctl_member_ban.md)	 - Ban a user from a chat
* [tgctl member ban-sender](tgctl_member_ban-sender.md)	 - Ban a channel from posting as itself in a chat
* [tgctl member decline-join](tgctl_member_decline-join.md)	 - Decline a chat join request
* [tgctl member promote](tgctl_member_promote.md)	 - Promote or demote an administrator
* [tgctl member restrict](tgctl_member_restrict.md)	 - Restrict what a member can do
* [tgctl member set-title](tgctl_member_set-title.md)	 - Set a custom title for an administrator the bot promoted
* [tgctl member unban](tgctl_member_unban.md)	 - Unban a previously banned user
* [tgctl member unban-sender](tgctl_member_unban-sender.md)	 - Unban a channel that was banned as a sender

