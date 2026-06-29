## tgctl media document

Send a document/file

```
tgctl media document [flags]
```

### Examples

```
  tgctl media document --chat @me --document ./report.pdf --caption "Q2 report"
```

### Options

```
      --caption string      document caption
      --chat string         target chat: numeric id or @username
      --document string     local path, URL, or file_id
  -h, --help                help for document
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
      --silent              send without a notification sound
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --profile string    profile/instance to use (env TGCTL_PROFILE)
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl media](tgctl_media.md)	 - Send files: photos, documents, and video

