## tgctl mcp vscode

Manage VSCode MCP servers

### Synopsis

Manage MCP server configuration for Visual Studio Code

### Options

```
  -h, --help   help for vscode
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
* [tgctl mcp vscode disable](tgctl_mcp_vscode_disable.md)	 - Remove server from VSCode config
* [tgctl mcp vscode enable](tgctl_mcp_vscode_enable.md)	 - Add server to VSCode config
* [tgctl mcp vscode list](tgctl_mcp_vscode_list.md)	 - Show VSCode MCP servers

