## tgctl checklist edit

Edit a checklist sent on behalf of a business account

### Synopsis

Edit an interactive checklist message sent on behalf of a connected business account. Returns the edited Message.

```
tgctl checklist edit [flags]
```

### Examples

```
  tgctl checklist edit --business-connection-id bc123 --chat 12345678 --message-id 42 --checklist '{"title":"Updated Tasks","tasks":[{"id":1,"text":"First task"}]}'
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
      --chat string                     target chat: numeric id or @username
      --checklist string                JSON object describing the updated checklist (InputChecklist)
  -h, --help                            help for edit
      --message-id int                  message id
      --reply-markup string             JSON object for an inline keyboard
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

