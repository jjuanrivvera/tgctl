# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Reaction moderation and member tags (issue #29, first batch).** Three methods the Bot API
  gained after the manifest's old baseline, all small additions to groups that already existed:
  - `tgctl message unreact --chat … --message-id … --user …` (`deleteMessageReaction`) removes a
    reaction somebody *else* left. It is moderation, not the inverse of `message react` — that
    one clears the bot's own reactions by omitting `--reaction` — so it needs the
    `can_delete_messages` right. `--actor-chat` names a channel that reacted on its own behalf.
  - `tgctl message unreact-all --chat … --user …` (`deleteAllMessageReactions`) sweeps up to
    10000 of one actor's recent reactions across a chat. It deliberately takes no
    `--message-id`: it is about an actor, not a message.
  - `tgctl member set-tag --chat … --user … [--tag …]` (`setChatMemberTag`) sets the tag shown
    beside a regular member's name; omit `--tag` to clear it. Needs `can_manage_tags`.

  Both reaction verbs are classified **destructive**, so the agent guard and the MCP server
  gate them like any other delete.
- **Three more from the same batch**: `tgctl user audios` (`getUserProfileAudios`, the twin of
  `user photos`), `tgctl user chat-messages` (`getUserPersonalChatMessages` — the last messages
  on the chat a user pinned to their profile; `--limit` is required and capped at 20 by the
  API), and `tgctl member answer-join-query` (`answerChatJoinRequestQuery`, the Mini App path
  for a join request, as opposed to `approve-join`/`decline-join` which act on a chat and user
  directly). `--result` accepts only `approve`, `decline` or `queue`, checked before the
  request: in a join flow somebody is waiting on the other side of a typo.
- **And two more**: `tgctl message draft` (`sendMessageDraft`) streams a partial message while
  the real one is still being written — the "typing out an answer" effect for a private chat.
  The draft is ephemeral, so the help says outright that you still have to `message send` the
  final text for it to exist in the chat; reusing `--draft-id` animates the change and an empty
  `--text` shows the "Thinking…" placeholder. `tgctl inline answer-guest` (`answerGuestQuery`)
  replies to a guest message with a **single** result object, as opposed to the array
  `inline answer` takes.

## [0.4.0] - 2026-09-23

### Added
- **`tgctl api` can upload files (issue #24).** The escape hatch only spoke JSON, so every
  method that *requires* multipart — `setMyProfilePhoto`, `setChatPhoto`, `uploadStickerFile`,
  `editMessageMedia` with `attach://`, … — was unreachable, and the only way out was to read the
  token out of the keyring and fall back to curl, which is exactly what tgctl exists to avoid.
  New repeatable `-F, --file name=@path`: one or more of them switches the request to
  `multipart/form-data`, sending each file as the part `<name>` while `-d`/`-q` ride along as
  fields, so a JSON body can point at an upload with `"attach://<name>"`:

  ```
  tgctl api setMyProfilePhoto -d '{"photo":{"type":"static","photo":"attach://pic"}}' -F pic=@logo.jpg
  ```

  `--dry-run` prints the equivalent curl with its `-F` parts and the token redacted, and a dry
  run no longer streams the file into memory just to describe it.
- **`tgctl bot set-photo` / `remove-photo` / `photo` (issue #23).** The bot's own profile photo
  was the one piece of its identity tgctl could not touch: `setMyProfilePhoto` and
  `removeMyProfilePhoto` were unwrapped, so the only routes were BotFather's `/setuserpic` or
  pulling the token out of the keyring for curl.
  - `tgctl bot set-photo --photo logo.jpg` — static .JPG.
  - `tgctl bot set-photo --animated intro.mp4 [--main-frame-timestamp 1.5]` — MPEG4, choosing
    the frame Telegram shows as the still image.
  - `tgctl bot remove-photo` (classified destructive) and `tgctl bot photo`, which asks `getMe`
    for the bot's own id and reports whether it has a photo.

  A profile photo cannot be reused, so Telegram accepts neither a URL nor a `file_id`: the
  command says so up front instead of forwarding the value and relaying a 400. The wire shape
  (an `InputProfilePhoto` object pointing at the upload with `attach://`, and the animated
  variant naming its field `animation`) is pinned in DECISIONS.md.

## [0.3.5] - 2026-09-22

### Security
- **The bot token no longer leaks when a request fails at the transport level (issue #21).**
  The Bot API carries its credential in the URL path, and Go's `*url.Error` — which
  `http.Client.Do` returns for every dropped connection, DNS failure or read timeout — prints
  that URL verbatim. The secret therefore reached stderr, any log capturing it, and the tool
  result an MCP client stores in its transcript, bypassing the keyring that exists precisely
  so the token never touches disk in clear text. Redaction now happens in two places:
  - `internal/api` masks the credential in every error leaving the network path (method calls
    and the `/file/bot<token>/…` download URL alike), keeping the non-secret bot id so an
    error still says which bot failed, and preserving the cause so `errors.Is`/`errors.As`
    keep working.
  - `cmd/tgctl` masks it again on the way to stderr — the last boundary before text leaves
    the process, and the one an MCP client reads.

  `--dry-run` already redacted (and still honors the explicit `--show-token`); `--verbose`
  logs responses, never the request URL.

## [0.2.2] - 2026-07-12

### Security
- Bump the Go toolchain to **go1.25.12**, clearing the reachable standard-library advisories
  (crypto/tls GO-2026-5856, crypto/x509, net/http, net/textproto) that govulncheck flagged.

## [0.2.1] - 2026-07-12

### Fixed
- The hidden secret prompt now reads in **raw mode** instead of `term.ReadPassword`'s canonical
  mode (capped at MAX_CANON, 1024 bytes on macOS), so a long pasted token no longer hangs the
  prompt until Ctrl-C. Bracketed-paste markers are still stripped as a defensive guard.

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
