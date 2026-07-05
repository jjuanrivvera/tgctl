## tgctl mcp stream

Stream the MCP server over HTTP

### Synopsis

Start HTTP server to expose CLI commands to AI assistants

```
tgctl mcp stream [flags]
```

### Options

```
  -h, --help               help for stream
      --host string        host to listen on
      --log-level string   Log level (debug, info, warn, error)
      --port int           port number to listen on (default 8080)
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

* [tgctl mcp](tgctl_mcp.md)	 - MCP server management

