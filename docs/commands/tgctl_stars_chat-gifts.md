## tgctl stars chat-gifts

List the gifts a chat owns

### Synopsis

List the gifts owned by a channel (getChatGifts). Whether unsaved gifts are visible
depends on the bot holding can_post_messages in that channel.

```
tgctl stars chat-gifts [flags]
```

### Examples

```
  tgctl stars chat-gifts --chat @mychannel
```

### Options

```
      --chat string                      target chat: numeric id or @username
      --exclude-from-blockchain          skip gifts held on the blockchain
      --exclude-limited-non-upgradable   skip limited gifts that cannot become unique
      --exclude-limited-upgradable       skip limited gifts that can become unique
      --exclude-saved                    skip gifts saved to the profile page
      --exclude-unique                   skip unique gifts
      --exclude-unlimited                skip gifts sold in unlimited numbers
      --exclude-unsaved                  skip gifts not saved to the profile page
  -h, --help                             help for chat-gifts
      --limit int                        max gifts to return (1-100)
      --offset string                    opaque cursor from a previous page
      --sort-by-price                    order by value instead of by date
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

* [tgctl stars](tgctl_stars.md)	 - Telegram Stars: transactions, gifts, and paid media

