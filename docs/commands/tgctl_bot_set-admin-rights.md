## tgctl bot set-admin-rights

Set the bot's default administrator rights (requested when added to a group/channel)

```
tgctl bot set-admin-rights [flags]
```

### Examples

```
  tgctl bot set-admin-rights --rights '{"can_manage_chat":true,"can_delete_messages":true}'
```

### Options

```
      --for-channels    apply to channels instead of groups/supergroups
  -h, --help            help for set-admin-rights
      --rights string   ChatAdministratorRights object as JSON (omit to clear)
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

* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself

