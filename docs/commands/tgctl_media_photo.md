## tgctl media photo

Send a photo

```
tgctl media photo [flags]
```

### Examples

```
  tgctl media photo --chat @me --photo ./cat.jpg --caption "my cat"
  tgctl media photo --chat @me --photo https://example.com/pic.png
```

### Options

```
      --caption string      photo caption
      --chat string         target chat: numeric id or @username
  -h, --help                help for photo
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --photo string        local path, URL, or file_id
      --silent              send without a notification sound
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

* [tgctl media](tgctl_media.md)	 - Send files: photos, documents, and video

