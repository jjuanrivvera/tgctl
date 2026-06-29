# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
