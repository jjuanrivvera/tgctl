## tgctl

Command-line tool for the Telegram Bot API

### Synopsis

tgctl is a fast, scriptable command-line tool for the Telegram Bot API.

It wraps the Bot API methods (sendMessage, getChat, getUpdates, ...) behind ergonomic
commands with table/json/yaml/csv output, named profiles for multiple bots, OS-keyring
token storage, and an MCP server so AI agents can drive it safely.

Get a bot token from @BotFather, then:

  tgctl auth login                       # store the token in your OS keyring
  tgctl bot info                         # who am I?
  tgctl message send --chat @me --text "hello from tgctl"
  tgctl updates get --limit 5 -o json    # poll recent updates as JSON

Every command honors --dry-run (prints the equivalent curl), -o/--output, and --jq.

### Options

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
  -h, --help              help for tgctl
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
* [tgctl alias](tgctl_alias.md)	 - Manage user-defined command aliases
* [tgctl api](tgctl_api.md)	 - Call any Bot API method directly (raw escape hatch)
* [tgctl auth](tgctl_auth.md)	 - Manage bot tokens and verify authentication
* [tgctl bot](tgctl_bot.md)	 - Inspect and configure the bot itself
* [tgctl callback](tgctl_callback.md)	 - Answer callback queries from inline keyboards
* [tgctl chat](tgctl_chat.md)	 - Inspect chats and their members
* [tgctl commands](tgctl_commands.md)	 - Manage the bot's command menu
* [tgctl completion](tgctl_completion.md)	 - Generate a shell completion script
* [tgctl config](tgctl_config.md)	 - Inspect and edit tgctl configuration
* [tgctl doctor](tgctl_doctor.md)	 - Diagnose configuration, credentials, and connectivity
* [tgctl file](tgctl_file.md)	 - Inspect and download files
* [tgctl init](tgctl_init.md)	 - First-run wizard: pick a base URL, capture a token, and smoke-test
* [tgctl inline](tgctl_inline.md)	 - Answer inline queries
* [tgctl invite](tgctl_invite.md)	 - Manage chat invite links
* [tgctl mcp](tgctl_mcp.md)	 - MCP server management
* [tgctl media](tgctl_media.md)	 - Send files: photos, documents, and video
* [tgctl member](tgctl_member.md)	 - Moderate chat members (ban, restrict, promote)
* [tgctl message](tgctl_message.md)	 - Send and manage messages
* [tgctl updates](tgctl_updates.md)	 - Fetch incoming updates (long polling)
* [tgctl user](tgctl_user.md)	 - Read user information
* [tgctl version](tgctl_version.md)	 - Print version, commit, and build date
* [tgctl webhook](tgctl_webhook.md)	 - Manage the bot's webhook

