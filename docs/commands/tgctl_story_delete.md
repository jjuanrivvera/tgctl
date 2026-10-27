## tgctl story delete

Delete a story posted on behalf of a business account

### Synopsis

Delete a story previously posted on behalf of a connected business account (deleteStory).
Returns True on success.

```
tgctl story delete [flags]
```

### Examples

```
  tgctl story delete --business-connection-id bc123 --story-id 42
```

### Options

```
      --business-connection-id string   id of the business connection acting on the account
  -h, --help                            help for delete
      --story-id int                    unique identifier of the story
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

* [tgctl story](tgctl_story.md)	 - Manage and publish stories on behalf of business accounts

