## tgctl invite

Manage chat invite links

### Synopsis

Create, edit, and revoke additional invite links for a chat (the bot must be an admin with can_invite_users).

### Options

```
  -h, --help   help for invite
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
* [tgctl invite create](tgctl_invite_create.md)	 - Create a new invite link
* [tgctl invite edit](tgctl_invite_edit.md)	 - Edit an existing invite link
* [tgctl invite revoke](tgctl_invite_revoke.md)	 - Revoke an invite link (a new one is generated automatically)

