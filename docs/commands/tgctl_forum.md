## tgctl forum

Manage forum topics in supergroups

### Synopsis

Create, edit, close, reopen, and delete forum topics, plus manage the General topic. The bot must be an admin with can_manage_topics.

### Options

```
  -h, --help   help for forum
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
* [tgctl forum close](tgctl_forum_close.md)	 - Close a forum topic
* [tgctl forum close-general](tgctl_forum_close-general.md)	 - Close the General forum topic
* [tgctl forum create](tgctl_forum_create.md)	 - Create a forum topic
* [tgctl forum delete](tgctl_forum_delete.md)	 - Delete a forum topic and all its messages
* [tgctl forum edit](tgctl_forum_edit.md)	 - Edit a forum topic's name or icon
* [tgctl forum edit-general](tgctl_forum_edit-general.md)	 - Rename the General forum topic
* [tgctl forum hide-general](tgctl_forum_hide-general.md)	 - Hide the General forum topic (also closes it)
* [tgctl forum icon-stickers](tgctl_forum_icon-stickers.md)	 - List the custom emoji stickers usable as topic icons
* [tgctl forum reopen](tgctl_forum_reopen.md)	 - Reopen a closed forum topic
* [tgctl forum reopen-general](tgctl_forum_reopen-general.md)	 - Reopen the General forum topic (also unhides it)
* [tgctl forum unhide-general](tgctl_forum_unhide-general.md)	 - Unhide the General forum topic
* [tgctl forum unpin-all](tgctl_forum_unpin-all.md)	 - Unpin all messages in a forum topic
* [tgctl forum unpin-all-general](tgctl_forum_unpin-all-general.md)	 - Unpin all messages in the General forum topic

