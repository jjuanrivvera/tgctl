<div align="center">

# tgctl

[![CI](https://github.com/jjuanrivvera/tgctl/actions/workflows/ci.yml/badge.svg)](https://github.com/jjuanrivvera/tgctl/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/jjuanrivvera/tgctl)](https://github.com/jjuanrivvera/tgctl/releases/latest)
[![Coverage](https://img.shields.io/badge/coverage-%E2%89%A580%25-brightgreen)](https://github.com/jjuanrivvera/tgctl/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jjuanrivvera/tgctl.svg)](https://pkg.go.dev/github.com/jjuanrivvera/tgctl)
[![Go version](https://img.shields.io/github/go-mod/go-version/jjuanrivvera/tgctl)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/jjuanrivvera/tgctl)
[![Built with cliwright](https://img.shields.io/badge/built_with-cliwright-1f6feb)](https://cliwright.jjuanrivvera.com)

**A gh-style CLI for the Telegram Bot API — messages, chats, webhooks, and local message history.**

[Documentation](https://jjuanrivvera.github.io/tgctl/) · [Commands](https://jjuanrivvera.github.io/tgctl/commands/)

</div>

---

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

### curl | sh (macOS/Linux)
Downloads the release archive for your OS/arch and verifies its SHA-256 against the release
`checksums.txt` before installing:
```sh
curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/tgctl/main/install.sh | sh
```

### Homebrew (macOS/Linux)
```sh
brew install jjuanrivvera/tgctl/tgctl-cli
```

### Scoop (Windows)
```sh
scoop bucket add tgctl https://github.com/jjuanrivvera/scoop-tgctl
scoop install tgctl
```

### Docker
```sh
docker run --rm -e TGCTL_TOKEN=123:ABC ghcr.io/jjuanrivvera/tgctl bot info
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
tgctl media audio --chat @me --audio ./song.mp3  # audio, voice, animation, sticker, ...
tgctl message poll --chat @g --question "Lunch?" --options '[{"text":"A"},{"text":"B"}]'
tgctl message location --chat @me --latitude 3.45 --longitude -76.53
tgctl file download --file-id <id> --dest ./out.jpg   # getFile + download the bytes
tgctl chat get --chat @telegram -o json          # chat metadata as JSON
tgctl chat administrators --chat @mygroup        # list admins
tgctl invite create --chat @mygroup --member-limit 100   # one-off invite link
tgctl callback answer --callback-query-id <id> --text "Saved!"
tgctl forum create --chat @mygroup --name "Announcements"  # forum topic management
tgctl member set-title --chat @mygroup --user 123 --title "Community Lead"
tgctl stars transactions --limit 20              # Telegram Stars balance ledger
tgctl webhook listen --port 8080 -o json         # receive + print webhook updates locally
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

## Multiple bots (`--bot`)

A profile is one bot, so select it with `--bot` (the old `--profile` flag still works as a
hidden alias):

```sh
tgctl auth login --bot staging              # a second bot
tgctl --bot staging bot info
tgctl config use staging                    # make it the default
tgctl config list-profiles
```

Bot selection precedence: `--bot` > `$TGCTL_BOT` > `$TGCTL_PROFILE` (legacy) > the active bot.
Token precedence: `$TGCTL_TOKEN` > `$TELEGRAM_BOT_TOKEN` > the active bot's keyring entry.
Point `--base-url` at a [self-hosted Local Bot API Server](https://github.com/tdlib/telegram-bot-api)
if you run one.

## AI agents (MCP + guard)

`tgctl` ships an MCP server so an agent can call the Bot API through annotated tools:

```sh
tgctl mcp start                  # run as an MCP server (stdio)
tgctl mcp claude                 # install into Claude Desktop
```

Setup/secret commands (`auth`, `config`, the raw `api` hatch) are **not** exposed, and the
token and bot-selection flags (`--bot`/`--profile`) never reach the tool schema. Generate host
safety rules that hard-block irreversible operations:

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
make verify       # deterministic gate: check + spec-check + spec-completeness + coverage + DoD
make judge        # the LLM-scored subjective gate (needs an agent; build-acceptance only)
make accept       # verify + judge — the full build-acceptance gate
```

`make verify` is what CI runs; it is deterministic and spends no tokens. `make judge` is kept
out of the routine gate because it needs an agent (`claude`/`codex`) and scores subjective
Definition-of-Done items an LLM can judge but a grep can't.

See [AGENTS.md](AGENTS.md) for the architecture and house rules, and
[DECISIONS.md](DECISIONS.md) for the pinned design rulings.

## Roadmap / Pending

`tgctl` wraps **109 of the 135** methods in the Telegram Bot API (v8.3, enumerated from the
[ark0f/tg-bot-api](https://ark0f.github.io/tg-bot-api/) machine spec) — **80%** coverage,
enforced by the `spec-completeness` gate. Done since the first cut:

- **Bot API coverage 52 → 109 verbs** — full member/admin management, forum topics (incl. the
  General topic), chat-admin setters (photo, permissions, sticker set, menu button, unpin-all),
  the `edit*` message family + bulk copy/forward/delete, chat/user verification, Telegram Stars
  (transactions, gifts, refunds, subscriptions, paid media), invite export + subscription links,
  and bot short-description / default-admin-rights / close / logout.
- **Packaging** — multi-stage `Dockerfile` + a distroless GHCR image, and a checksum-verifying
  `install.sh` (`curl | sh`) one-liner.

Deliberately deferred (recorded as a `coverage-waiver` in [DECISIONS.md](DECISIONS.md) so the
completeness gate accounts for them, not a silent gap): five self-contained niche families —
**sticker-set authoring, payments/invoices, games, Telegram Passport, and business-connection**.
Each is the same declarative pattern to add (manifest verb → group `methodCmd` → mocked test)
when a real need shows up; the waiver line is removed once coverage clears the threshold.

## License

MIT © Juan Rivera. See [LICENSE](LICENSE).
