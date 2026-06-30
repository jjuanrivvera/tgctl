## tgctl agent guard

Generate agent-safety config that blocks destructive tgctl operations

### Synopsis

Classify every API command (read / write / irreversible) from the live command
tree and emit host safety config: irreversible verbs (delete, leave, ban, webhook delete)
are hard-blocked, ordinary writes require approval, and reads are allowed.

MCP-only operation is the hard guarantee; the Bash patterns are best-effort (they defeat
quoting tricks, not variable indirection).

```
tgctl agent guard --host <claude-code|codex|opencode> [flags]
```

### Examples

```
  tgctl agent guard --host claude-code
  tgctl agent guard --host codex --out ~/.codex/config.toml
  tgctl agent guard --host opencode --all-writes
```

### Options

```
      --all-writes    also hard-block ordinary writes, not just irreversible ops
  -h, --help          help for guard
      --host string   target agent host: claude-code|codex|opencode (required)
      --out string    write to this file instead of stdout
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

* [tgctl agent](tgctl_agent.md)	 - AI-agent integration helpers

