## tgctl message contact

Send a phone contact

```
tgctl message contact [flags]
```

### Examples

```
  tgctl message contact --chat @me --phone-number "+15551234567" --first-name "Ada"
```

### Options

```
      --chat string           target chat: numeric id or @username
      --first-name string     contact's first name
  -h, --help                  help for contact
      --last-name string      contact's last name
      --phone-number string   contact's phone number
      --silent                send without a notification sound
      --vcard string          additional data about the contact as a vCard (0-2048 bytes)
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

