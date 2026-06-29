# tgctl

A fast, scriptable command-line tool for the [Telegram Bot API](https://core.telegram.org/bots/api).

- Send messages and media, manage chats/members, configure webhooks and the command menu.
- table / json / yaml / csv output, `--columns`, and a built-in `--jq` filter.
- OS-keyring token storage, named profiles for multiple bots.
- An MCP server and an `agent guard` so AI agents can drive it safely.

## Get started

```sh
go install github.com/jjuanrivvera/tgctl/cmd/tgctl@latest
tgctl auth login          # paste a @BotFather token
tgctl bot info
tgctl message send --chat @me --text "hello from tgctl"
```

See the full [command reference](commands/tgctl.md), or the
[README](https://github.com/jjuanrivvera/tgctl#readme) for installation options.
