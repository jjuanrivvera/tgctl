## tgctl mcp claude

Manage Claude Desktop MCP servers

### Synopsis

Manage MCP server configuration for Claude Desktop

### Options

```
  -h, --help   help for claude
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

* [tgctl mcp](tgctl_mcp.md)	 - MCP server management
* [tgctl mcp claude disable](tgctl_mcp_claude_disable.md)	 - Remove server from Claude config
* [tgctl mcp claude enable](tgctl_mcp_claude_enable.md)	 - Add server to Claude config
* [tgctl mcp claude list](tgctl_mcp_claude_list.md)	 - Show Claude MCP servers

