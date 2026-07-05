## tgctl mcp

MCP server management

### Synopsis

Manage MCP servers for AI assistants and code editors

### Options

```
  -h, --help   help for mcp
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

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl mcp claude](tgctl_mcp_claude.md)	 - Manage Claude Desktop MCP servers
* [tgctl mcp cursor](tgctl_mcp_cursor.md)	 - Manage Cursor MCP servers
* [tgctl mcp start](tgctl_mcp_start.md)	 - Start the MCP server
* [tgctl mcp stream](tgctl_mcp_stream.md)	 - Stream the MCP server over HTTP
* [tgctl mcp tools](tgctl_mcp_tools.md)	 - Export tools as JSON
* [tgctl mcp vscode](tgctl_mcp_vscode.md)	 - Manage VSCode MCP servers

