## tgctl mcp claude list

Show Claude MCP servers

### Synopsis

Show all MCP servers configured in Claude Desktop

```
tgctl mcp claude list [flags]
```

### Options

```
      --config-path string   Path to Claude config file
  -h, --help                 help for list
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

* [tgctl mcp claude](tgctl_mcp_claude.md)	 - Manage Claude Desktop MCP servers

