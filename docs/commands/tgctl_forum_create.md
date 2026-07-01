## tgctl forum create

Create a forum topic

```
tgctl forum create [flags]
```

### Examples

```
  tgctl forum create --chat @group --name "Announcements" --icon-color 7322096
```

### Options

```
      --chat string                   target chat: numeric id or @username
  -h, --help                          help for create
      --icon-color int                RGB color of the topic icon (one of the allowed palette values)
      --icon-custom-emoji-id string   custom emoji id shown as the topic icon
      --name string                   topic name (1-128 chars)
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

* [tgctl forum](tgctl_forum.md)	 - Manage forum topics in supergroups

