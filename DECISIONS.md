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
| Profile flag name | User-facing flag is **`--bot`** (env `TGCTL_BOT`); `--profile`/`TGCTL_PROFILE` kept as hidden, still-working aliases | For Telegram a profile *is* a bot, so `--bot` reads truer; the alias keeps existing scripts working. Both names are excluded from the MCP tool schema. |
| Float params | Generic builder gained a `flagFloat` kind | `sendLocation`/`sendVenue` take `latitude`/`longitude` floats; string-encoding them risks API rejection. |
| `base_url` override | `--base-url` / config (default `https://api.telegram.org`) | Supports self-hosted [Local Bot API Server](https://github.com/tdlib/telegram-bot-api). |
| CSV output | Kept | Most list results (updates, administrators, commands) are tabular. |
| "id" rendering | `ID` flexible type (string-or-number) | chat_id / user_id are large int64; render as string to avoid >2^53 precision loss. |

## Beyond-the-API value-adds (GOAL.md §3c)
- **`webhook listen`** — a local HTTP receiver that renders incoming webhook updates. It is
  not a single Bot API method, so it's a hand-written `Extra` command on the webhook group
  (the generic builder gained an `Extra []func() *cobra.Command` hook — extend, don't fork).
  Excluded from the MCP surface (a blocking server must never be an agent tool). Not in
  `api-manifest.json` because that manifest tracks the pure API surface; spec-check only
  enforces a minimum, so value-adds beyond it are allowed.
- **`file download`** — resolves a `file_id` (getFile) and streams the file's bytes from the
  Bot API's `/file/` endpoint to a local path or stdout. Two steps, not one method, so it's an
  `Extra` on the `file` group (same pattern as `webhook listen`); `file info` is the pure
  getFile wrap. Not in `api-manifest.json` for the same reason. The token is embedded by the
  authenticator (`FileURL`) and never logged.

## Resource set (derived from the Bot API method surface; see api-manifest.json)
Grouped by noun; verbs map 1:1 to Bot API methods. Read-only verbs annotated read-only for
MCP/agent-guard; destructive verbs (delete/leave/ban/unpin) annotated destructive.

## API completeness (cliwright GOAL.md §0/§11)
The full Bot API method set is **enumerated from a source, not recalled**: the ark0f/tg-bot-api
community machine spec (`api_method_source` in api-manifest.json) yields **135 methods** (Bot API
8.3). `scripts/spec-completeness.sh` reads `api_method_total` and fails if the manifest covers
below its threshold without a recorded waiver.

`tgctl` wraps **109 / 135 methods (80%)**: all messaging, media, chat administration, member
moderation, forum-topic management, invite links (incl. subscription links), bot configuration,
Telegram Stars, chat/user verification, webhooks, updates, files, callbacks, and inline queries.

**coverage-waiver: 80% (109/135). The 26 uncovered methods are five genuinely-niche families
deferred deliberately, not overlooked** — each is a self-contained sub-API most bots never touch:
- **stickers-set-management (15)** — createNewStickerSet, addStickerToSet, replaceStickerInSet,
  deleteStickerFromSet, deleteStickerSet, setSticker*, getStickerSet, getCustomEmojiStickers,
  uploadStickerFile, setCustomEmojiStickerSetThumbnail. A full sticker-authoring workflow.
- **payments/invoices (4)** — sendInvoice, createInvoiceLink, answerShippingQuery,
  answerPreCheckoutQuery. Requires a payment-provider token and a checkout callback loop.
- **games (3)** — sendGame, setGameScore, getGameHighScores. The HTML5 Games platform.
- **telegram-passport (1)** — setPassportDataErrors. Encrypted identity documents.
- **business-connection (3)** — getBusinessConnection, answerWebAppQuery,
  savePreparedInlineMessage. Telegram Business / Web-App-specific.

They stay recorded here (not silently dropped) so the completeness gate sees the decision on
every pass. Adding any family later is the same declarative pattern (manifest verb → group
methodCmd → mocked test); the waiver line is removed when coverage clears the threshold.

## v0.2 surface expansion (52 → 109 verbs)
Added following the generic method-command builder (extend, don't fork):
- **message** — sendChatAction (`action`), the edit* family (edit-caption/edit-media/
  edit-reply-markup/edit-live-location/stop-live-location), stopPoll, and the bulk
  copy/forward/delete-batch methods. edit* targets are `--chat/--message-id` **or**
  `--inline-message-id`, so those flags are optional (`optChatFlag`/`optMessageIDFlag`).
- **chat** — admin setters: set/delete-photo, set-permissions, set/delete-sticker-set,
  get/set-menu-button, unpin-all, and getUserChatBoosts (`boosts`).
- **member** — setChatAdministratorCustomTitle (`set-title`), approve/decline-join,
  ban/unban-sender (channel-as-sender bans).
- **forum** (new) — full topic lifecycle + the General-topic variants.
- **verify** (new) — verify/unverify chats and users (org-verified bots only).
- **stars** (new) — Star transactions, gifts, refunds, subscription edit, emoji status,
  and paid media.
- **bot** — short-description get/set, default-admin-rights get/set, close, logout.
- **invite** — exportChatInviteLink and the subscription-invite-link pair.
