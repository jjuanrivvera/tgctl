## tgctl message poll

Send a native poll

### Synopsis

Send a poll. --options is a JSON array of InputPollOption objects, e.g. '[{"text":"Yes"},{"text":"No"}]'.

```
tgctl message poll [flags]
```

### Examples

```
  tgctl message poll --chat @group --question "Lunch?" \
    --options '[{"text":"Pizza"},{"text":"Sushi"}]'
```

### Options

```
      --allows-multiple-answers   allow multiple answers (regular polls only)
      --anonymous                 make the poll anonymous (Telegram default: true)
      --chat string               target chat: numeric id or @username
      --correct-option-id int     0-based id of the correct option (quiz polls)
  -h, --help                      help for poll
      --options string            JSON array of InputPollOption objects (2-10)
      --question string           poll question (1-300 chars)
      --silent                    send without a notification sound
      --type string               poll type: regular | quiz
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

