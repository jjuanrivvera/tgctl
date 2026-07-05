# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Local SQLite message history (issue #5)**: every outbound send (and, in
  `updates get`/`webhook listen` mode, every inbound update) is now recorded to a
  per-bot-profile SQLite database, since the Bot API itself exposes no history endpoint. New
  `tgctl log` command family:
  - `tgctl log [--chat <id>] [--since 24h|RFC3339|YYYY-MM-DD] [--kind text] [--limit 50]` —
    list recorded messages.
  - `tgctl log search <query>` — full-text search (FTS5 `MATCH` when available, degrading
    automatically to a `LIKE` scan otherwise).
  - `tgctl log show <message_id>` — one message including its full raw API payload.
  - `tgctl log prune --older-than <duration>` — delete rows older than a cutoff.
  - New persistent `--no-store` flag disables recording for a single invocation. The store is
    always best-effort on the write path: a failed/unavailable store never breaks a send.
  - `log`/`log search`/`log show` are exposed to the MCP server (read-only); `log prune` is
    destructive. See DECISIONS.md for the full write-up.

### Fixed
- The message store's SQLite handle is now closed when a command finishes:
  `(*api.Client).Close()` closes an attached recorder if it implements `io.Closer`,
  and every `clientFromCmd` call site defers it. The handle was previously never
  closed, which passed on Unix but broke Windows CI (an open file handle blocks
  deleting/renaming it, so `t.TempDir()` cleanup failed for nearly every command
  test). `--dry-run` also now skips opening the store entirely, since it makes no
  API call and has nothing to record.

## [0.2.0] - 2026-07-02

### Added
- **`agent guard` now generates a PreToolUse enforcement hook** (Bash + MCP
  matchers) instead of only permission rules. It anchors blocked subcommands at
  the command position, matches path-invoked binaries (`./bin/tgctl`,
  `/usr/local/bin/tgctl`) while ignoring a different binary that merely ends in
  `tgctl`, emits every cobra-alias spelling (`msg delete`, `delete-many`,
  `cmds delete`), and — because the `api` escape is RPC-style and Telegram
  method names are case-insensitive — allows only `get*` at the method position.
- Expanded Bot API coverage with new verbs:
  - Media sends: `media audio`, `media voice`, `media animation`, `media video-note`,
    `media sticker`, `media media-group`.
  - Rich message sends: `message location`, `message venue`, `message contact`, `message poll`,
    `message dice`, and `message react` (setMessageReaction).
  - Files: `file info` (getFile) and `file download` (getFile + stream the bytes to disk).
  - Callbacks/inline: `callback answer` (answerCallbackQuery), `inline answer` (answerInlineQuery).
  - Chat admin: `invite create|edit|revoke`, `chat set-title`, `chat set-description`,
    `user photos` (getUserProfilePhotos).

### Changed
- The multi-bot selection flag is now `--bot` (a profile is one bot). `--profile` remains as a
  hidden, still-working alias, and `$TGCTL_BOT` is recognized ahead of the legacy `$TGCTL_PROFILE`.
- `agent guard` promotes `stars refund` and the bulk `unpin-all*` commands to the
  hard-block (irreversible) bucket, matching the guard's contract.

### Fixed
- `agent guard` closes hook bypasses that permission rules alone could not: no
  enforcement hook was generated before, and separators glued to a no-arg verb
  and a no-jq fallback that could fail open are now handled.

## [0.1.0] - 2026-06-29

### Added
- Initial release of `tgctl`, a command-line tool for the Telegram Bot API.
- Resource groups mapping 1:1 to Bot API methods: `bot`, `message`, `media`, `chat`,
  `member`, `updates`, `webhook`, `commands`.
- Meta commands: `auth` (login/logout/status), `config`, `init`, `doctor`, `completion`,
  `alias`, `api` (raw escape hatch), `version`.
- Output formats: table, json, yaml, csv, id — with `--columns` and a built-in `--jq` filter.
- OS-keyring token storage with an encrypted-file fallback; named profiles for multiple bots.
- Resilient client: idempotent-only retries, `retry_after`-aware 429 handling, adaptive rate
  limiting, `--dry-run` (prints the equivalent redacted `curl`), Ctrl-C cancellation.
- MCP server (`mcp`) exposing the API as annotated tools, plus `agent guard` to generate
  host safety config for Claude Code, Codex, and OpenCode.
- `webhook listen` — a local receiver that renders incoming webhook updates.
- Generated command reference, GoReleaser packaging, and CI.

[0.1.0]: https://github.com/jjuanrivvera/tgctl/releases/tag/v0.1.0
