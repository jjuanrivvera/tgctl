# tgctl — a command-line tool for the Telegram Bot API

`tgctl` wraps the [Telegram Bot API](https://core.telegram.org/bots/api) in a fast,
scriptable, `gh`-style CLI: send messages, manage chats and members, configure webhooks and
the bot's command menu, and poll updates — with table/JSON/YAML/CSV output, named profiles
for multiple bots, OS-keyring token storage, and an MCP server so AI agents can drive it
safely.

```console
$ tgctl auth login                 # store a @BotFather token in your OS keyring
$ tgctl bot info
ID         USERNAME   FIRST_NAME   CAN_JOIN_GROUPS
123456789  mybot      My Bot       true

$ tgctl message send --chat @me --text "hello from tgctl"
$ tgctl updates get --limit 5 -o json | jq '.[].message.text'
```

## Install

### Homebrew (macOS/Linux)
```sh
brew install jjuanrivvera/tgctl/tgctl-cli
```

### Scoop (Windows)
```sh
scoop bucket add tgctl https://github.com/jjuanrivvera/scoop-tgctl
scoop install tgctl
```

### Go
```sh
go install github.com/jjuanrivvera/tgctl/cmd/tgctl@latest
```

### deb / rpm / apk
Download from the [latest release](https://github.com/jjuanrivvera/tgctl/releases/latest).

## Quickstart

1. Create a bot and get its token from [@BotFather](https://t.me/BotFather).
2. `tgctl auth login` (or `tgctl init` for a guided wizard). The token is verified against
   `getMe` and stored in your OS keyring — never in a plaintext config file.
3. Run any command. Add `--dry-run` to print the equivalent `curl` (token redacted) without
   sending anything.

```sh
tgctl bot info                                   # who am I?
tgctl message send --chat @me --text "hi"        # send a message
tgctl media photo --chat @me --photo ./cat.jpg   # upload a photo
tgctl chat get --chat @telegram -o json          # chat metadata as JSON
tgctl chat administrators --chat @mygroup        # list admins
tgctl webhook info                               # webhook status
tgctl commands set --commands '[{"command":"start","description":"Begin"}]'
tgctl api getMe --idempotent                     # raw escape hatch for any method
```

## Output and filtering

A single renderer serves every command. Pick a format with `-o/--output`:

| Format | Use |
|---|---|
| `table` (default) | aligned, colored on a TTY (honors `NO_COLOR` / `--no-color`) |
| `json` | full record, big ids kept exact |
| `yaml` | readable structured output |
| `csv` | spreadsheet-friendly (formula-injection sanitized) |
| `id` | one id per line, pipeable to `xargs` |

- `--columns a,b,c` selects/orders columns; `--jq '<expr>'` filters with a built-in gojq.
- Notes and warnings go to **stderr** so stdout stays pipe-clean.

## Profiles (multiple bots)

```sh
tgctl auth login --profile staging          # a second bot
tgctl --profile staging bot info
tgctl config use staging                    # make it the default
tgctl config list-profiles
```

Token precedence: `$TGCTL_TOKEN` > `$TELEGRAM_BOT_TOKEN` > the active profile's keyring entry.
Point `--base-url` at a [self-hosted Local Bot API Server](https://github.com/tdlib/telegram-bot-api)
if you run one.

## AI agents (MCP + guard)

`tgctl` ships an MCP server so an agent can call the Bot API through annotated tools:

```sh
tgctl mcp start                  # run as an MCP server (stdio)
tgctl mcp claude                 # install into Claude Desktop
```

Setup/secret commands (`auth`, `config`, the raw `api` hatch) are **not** exposed, and the
token/profile flags never reach the tool schema. Generate host safety rules that hard-block
irreversible operations:

```sh
tgctl agent guard --host claude-code     # deny delete/leave/ban; ask on writes; allow reads
tgctl agent guard --host codex --out ~/.codex/config.toml
```

## Safety & reliability

- **Idempotent-only retries.** Reads retry on 5xx/network; writes retry **only** on 429
  (which means *rejected, not processed*), so a timed-out `sendMessage` is never double-sent.
- **Rate limiting.** Fixed RPS with halve-on-429 and gradual restore (the Bot API exposes no
  quota headers, only `retry_after`, which `tgctl` honors).
- **Actionable errors.** Each failure carries a hint keyed by status (401 → run `auth login`,
  429 → wait N seconds, …).
- **Ctrl-C** cancels in-flight polling and backoff cleanly.

## Documentation

Full command reference: [`docs/commands/`](docs/commands/tgctl.md). Build the docs site with
`make docs-serve` (MkDocs Material).

## Development

```sh
make build        # build to bin/tgctl
make check        # fmt + vet + lint + test
make verify       # the full acceptance gate (check + spec-check + coverage + DoD + judge)
```

See [AGENTS.md](AGENTS.md) for the architecture and house rules, and
[DECISIONS.md](DECISIONS.md) for the pinned design rulings.

## License

MIT © Juan Rivera. See [LICENSE](LICENSE).
