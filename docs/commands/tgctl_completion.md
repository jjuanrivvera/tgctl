## tgctl completion

Generate a shell completion script

### Synopsis

Output a completion script for your shell. See `tgctl completion <shell> --help` for install instructions.

```
tgctl completion [bash|zsh|fish|powershell]
```

### Examples

```
  source <(tgctl completion bash)
  tgctl completion zsh > "${fpath[1]}/_tgctl"
  tgctl completion fish > ~/.config/fish/completions/tgctl.fish
```

### Options

```
  -h, --help   help for completion
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API

