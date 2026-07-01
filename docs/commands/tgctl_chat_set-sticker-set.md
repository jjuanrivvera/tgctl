## tgctl chat set-sticker-set

Set the group sticker set for a supergroup

```
tgctl chat set-sticker-set [flags]
```

### Examples

```
  tgctl chat set-sticker-set --chat @group --sticker-set-name MyPack
```

### Options

```
      --chat string               target chat: numeric id or @username
  -h, --help                      help for set-sticker-set
      --sticker-set-name string   name of the sticker set
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

* [tgctl chat](tgctl_chat.md)	 - Inspect chats and their members

