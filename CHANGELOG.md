# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
