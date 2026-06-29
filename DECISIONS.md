# DECISIONS — tgctl (Telegram Bot API CLI)

Pinned assumptions and design rulings (cliwright GOAL.md §11). Read on every iteration;
never silently re-decide.

## Target API
- **Telegram Bot API** — `https://api.telegram.org/bot<token>/<METHOD>`. This is Telegram's
  HTTP/JSON API. (The MTProto *client* API is a binary protocol requiring api_id/api_hash +
  phone auth; it is not an HTTP/REST API and is out of scope for an HTTP-wrapping CLI.)
- Response envelope: `{"ok":true,"result":...}` on success;
  `{"ok":false,"error_code":N,"description":"...","parameters":{...}}` on error.
- `parameters` may carry `retry_after` (seconds, on 429) and `migrate_to_chat_id`.
- Methods accept GET or POST; params as query, `application/json`, or `multipart/form-data`
  (file uploads). We send `application/json` for non-file calls, `multipart/form-data` for
  file uploads.

## Decision log
| Question | Decision | Why |
|---|---|---|
| Resource pattern (§11) | **Pattern B (service-layer / method-command)** | The Bot API is RPC-method-oriented (`sendMessage`, `getChat`, `getUpdates`); no resource exposes a uniform list/get/create/update/delete. A generic *method-command builder* keeps command files thin without faking CRUD. |
| HTTP verb for calls | **POST** with JSON body (multipart for files) | Uniform, avoids URL-length limits, matches the docs' file-upload requirement. |
| Retry of non-idempotent methods | Retry **only 429** (honoring `retry_after`) for writes; retry 429+5xx+network for reads | A 429 means the call was rejected (not processed) → safe to retry. A 5xx/network error after a `sendMessage` is ambiguous and could double-send → never auto-retry writes on those. Each method declares `Idempotent`. |
| Rate limiting | Fixed RPS (default 25/s) with **halve-on-429 + gradual restore** | The Bot API exposes no quota headers; only 429 + `retry_after`. |
| Auth methods | Single method: **bot token** (`<bot_id>:<hash>`) in the URL path | The Bot API has exactly one credential. Modeled behind the same `Authenticator` interface so it scales to the simple case. |
| Token env var | `TGCTL_TOKEN` (primary), `TELEGRAM_BOT_TOKEN` (recognized alias) | Namespaced primary per house rule; the conventional name accepted for convenience. |
| Profiles | **Yes** (multi-bot) | Operators commonly run several bots; profile records bot id + base URL, secret in keyring. |
| `base_url` override | `--base-url` / config (default `https://api.telegram.org`) | Supports self-hosted [Local Bot API Server](https://github.com/tdlib/telegram-bot-api). |
| CSV output | Kept | Most list results (updates, administrators, commands) are tabular. |
| "id" rendering | `ID` flexible type (string-or-number) | chat_id / user_id are large int64; render as string to avoid >2^53 precision loss. |

## Resource set (derived from the Bot API method surface; see api-manifest.json)
Grouped by noun; verbs map 1:1 to Bot API methods. Read-only verbs annotated read-only for
MCP/agent-guard; destructive verbs (delete/leave/ban/unpin) annotated destructive.
