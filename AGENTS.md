# AGENTS.md — working in the tgctl repo

`tgctl` is a command-line tool for the **Telegram Bot API**, built to the cliwright
standard (Go + Cobra + GoReleaser). This file orients an AI agent (or human) contributing.

## The one rule that matters
**`make verify` is the gate. It runs `make check` (fmt, vet, golangci-lint including gosec,
and tests), `spec-check`, `spec-completeness`, `cover-check` (≥80% total), and `dod-check.sh`.
Run `make build` first so `spec-check` inspects the current binary. Run `govulncheck ./...`
separately; `make judge` is a separate subjective acceptance step, not part of `verify`.
Run the full gate for changes to the command surface or documented behavior.

## Architecture (where things live)
- `internal/api/` — the generic client core (auth, retry, rate limit, dry-run curl, flexible
  JSON types, the typed `Call`/`CallInto`/`Upload`). Written once; never copy-paste per method.
- `commands/` — thin, declarative command groups. Adding a Bot API method is a few lines in a
  group file via `registerGroup` — **zero edits to shared code**. The generic builder stamps
  MCP read-only/write/destructive annotations from each command's `Kind`.
- `internal/{config,auth,output,version}` — profiles + manual precedence (no Viper), keyring
  token storage, the table/json/yaml/csv renderer, build metadata.
- `cmd/tgctl/main.go` — entry point: `signal.NotifyContext` (Ctrl-C cancels in-flight work)
  + alias expansion before cobra parses.

## House rules
- Comments explain **WHY**, not WHAT.
- Thread `cmd.Context()` everywhere; never `context.Background()` (it breaks Ctrl-C). Tests use
  `t.Context()`.
- Secrets live in the OS keyring — never in config-in-repo, code, or commit messages.
- Pin every ambiguous API assumption in `DECISIONS.md`; read it back, never silently re-decide.
- The resource set is derived from the Bot API surface (`api-manifest.json`), not hand-picked.
