## tgctl commands set

Set the bot's command menu

```
tgctl commands set [flags]
```

### Examples

```
  tgctl commands set --commands '[{"command":"start","description":"Begin"},{"command":"help","description":"Get help"}]'
```

### Options

```
      --commands string        array of {command,description} objects as JSON
  -h, --help                   help for set
      --language-code string   BCP-47 language code
      --scope string           BotCommandScope object as JSON (default: all private chats)
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

* [tgctl commands](tgctl_commands.md)	 - Manage the bot's command menu

