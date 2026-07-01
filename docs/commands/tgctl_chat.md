## tgctl chat

Inspect chats and their members

### Synopsis

Read chat metadata, member counts, administrators, and individual members; leave a chat.

### Options

```
  -h, --help   help for chat
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
* [tgctl chat administrators](tgctl_chat_administrators.md)	 - List a chat's administrators
* [tgctl chat boosts](tgctl_chat_boosts.md)	 - List the boosts a user added to a chat
* [tgctl chat delete-photo](tgctl_chat_delete-photo.md)	 - Delete a chat's photo
* [tgctl chat delete-sticker-set](tgctl_chat_delete-sticker-set.md)	 - Remove the group sticker set from a supergroup
* [tgctl chat get](tgctl_chat_get.md)	 - Show a chat's metadata
* [tgctl chat leave](tgctl_chat_leave.md)	 - Make the bot leave a chat
* [tgctl chat member](tgctl_chat_member.md)	 - Show one member's status in a chat
* [tgctl chat members-count](tgctl_chat_members-count.md)	 - Show the number of members in a chat
* [tgctl chat menu-button](tgctl_chat_menu-button.md)	 - Show the chat's menu button (default: the bot's global button)
* [tgctl chat set-description](tgctl_chat_set-description.md)	 - Change a chat's description
* [tgctl chat set-menu-button](tgctl_chat_set-menu-button.md)	 - Set the chat's menu button
* [tgctl chat set-permissions](tgctl_chat_set-permissions.md)	 - Set the default permissions for all members
* [tgctl chat set-photo](tgctl_chat_set-photo.md)	 - Set a chat's photo
* [tgctl chat set-sticker-set](tgctl_chat_set-sticker-set.md)	 - Set the group sticker set for a supergroup
* [tgctl chat set-title](tgctl_chat_set-title.md)	 - Change a chat's title
* [tgctl chat unpin-all](tgctl_chat_unpin-all.md)	 - Unpin every pinned message in a chat

