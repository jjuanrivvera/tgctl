## tgctl agent guard

Generate agent-safety config that blocks destructive tgctl operations

### Synopsis

Classify every API command (read / write / irreversible) from the live command
tree and emit host safety config: irreversible operations (delete, leave, ban, revoke,
logout, refund, unpin-all, webhook delete) are hard-blocked, ordinary writes require
approval, and reads are allowed. Cobra alias paths are covered too — "tgctl msg delete"
and "tgctl message delete-many" hit the same rails as "tgctl message delete".

For claude-code the output also includes a PreToolUse hook script
(.claude/hooks/tgctl-guard.sh): it strips quote/backslash obfuscation, matches blocked
subcommand paths at the command position even for path-invoked binaries (./bin/tgctl,
/usr/local/bin/tgctl), and gates the raw "tgctl api <method>" escape hatch — only
read-only get* methods pass, since Bot API method names are case-insensitive
server-side. "tgctl alias set" is denied so an agent cannot mint a new shorthand for a
blocked command.

MCP-only operation is the hard guarantee; the Bash rails are best-effort — the hook
defeats quoting tricks and path prefixes, but not variable indirection
(a=delete; tgctl message $a) or shell aliases. Conservative false positives are
accepted: a line that merely QUOTES a blocked command (echo "tgctl message delete")
is denied.

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
      --no-store          disable local SQLite send/receive history for this invocation (see tgctl log)
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl agent](tgctl_agent.md)	 - AI-agent integration helpers

