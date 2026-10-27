## tgctl checklist send

Send a checklist on behalf of a business account

### Synopsis

Send an interactive checklist on behalf of a connected business account. Returns the sent Message.

```
tgctl checklist send [flags]
```

### Examples

```
  tgctl checklist send --business-connection-id bc123 --chat 12345678 --checklist '{"title":"Tasks","tasks":[{"id":1,"text":"First task"}]}'
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --chat string                     target chat: numeric id or @username
      --checklist string                JSON object describing the checklist to send (InputChecklist)
  -h, --help                            help for send
      --message-effect-id string        unique identifier of the message effect
      --protect-content                 protect the content from forwarding and saving
      --reply-markup string             JSON object for an inline keyboard
      --reply-parameters string         JSON object describing message to reply to
      --silent                          send without a notification sound
      --thread int                      forum topic thread id
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

* [tgctl checklist](tgctl_checklist.md)	 - Send and edit interactive checklists

