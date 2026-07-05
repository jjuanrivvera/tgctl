## tgctl chat set-description

Change a chat's description

```
tgctl chat set-description [flags]
```

### Examples

```
  tgctl chat set-description --chat @group --description "What this group is about"
  tgctl chat set-description --chat @group --description ""   # clear it
```

### Options

```
      --chat string          target chat: numeric id or @username
      --description string   new chat description (0-255 chars; empty clears it)
  -h, --help                 help for set-description
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

* [tgctl chat](tgctl_chat.md)	 - Inspect chats and their members

