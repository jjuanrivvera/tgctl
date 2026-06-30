## tgctl message venue

Send information about a venue

```
tgctl message venue [flags]
```

### Examples

```
  tgctl message venue --chat @me --latitude 3.45 --longitude -76.53 --title "Office" --address "Av. 1 #2-3"
```

### Options

```
      --address string         address of the venue
      --chat string            target chat: numeric id or @username
      --foursquare-id string   Foursquare identifier of the venue
  -h, --help                   help for venue
      --latitude float         latitude of the venue
      --longitude float        longitude of the venue
      --silent                 send without a notification sound
      --title string           name of the venue
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

