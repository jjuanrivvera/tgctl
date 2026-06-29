## tgctl message copy

Copy a message (without a 'forwarded from' header)

```
tgctl message copy [flags]
```

### Examples

```
  tgctl message copy --chat @dest --from-chat @src --message-id 42 --caption "fyi"
```

### Options

```
      --caption string      new caption for the copied message
      --chat string         target chat: numeric id or @username
      --from-chat string    source chat id or @username
  -h, --help                help for copy
      --message-id int      message id
      --parse-mode string   text formatting: MarkdownV2 | HTML | Markdown
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

