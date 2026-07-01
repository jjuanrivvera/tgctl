## tgctl verify user

Verify a user on behalf of the bot's organization

```
tgctl verify user [flags]
```

### Examples

```
  tgctl verify user --user 12345 --custom-description "Verified staff"
```

### Options

```
      --custom-description string   custom verification description (0-70 chars)
  -h, --help                        help for user
      --user int                    target user id
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

* [tgctl verify](tgctl_verify.md)	 - Verify or unverify chats and users

