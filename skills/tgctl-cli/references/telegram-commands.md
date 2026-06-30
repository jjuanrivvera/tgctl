# Command surface

Groups map 1:1 to Bot API methods. Run `tgctl <group> <verb> --help` for flags and examples.

## bot
- `bot info` (getMe), `bot set-name`/`get-name`, `bot set-description`/`get-description`.

## message
- `send` (sendMessage): `--chat --text [--parse-mode MarkdownV2|HTML] [--reply-to] [--silent] [--no-preview] [--reply-markup JSON]`.
- `edit` (editMessageText), `delete` (deleteMessage), `forward`, `copy`, `pin`, `unpin`.
- `react` (setMessageReaction): `--chat --message-id [--reaction '<ReactionType[] JSON>'] [--is-big]` (omit `--reaction` to clear).
- `location` (sendLocation): `--chat --latitude --longitude [--live-period]`.
- `venue` (sendVenue): `--chat --latitude --longitude --title --address [--foursquare-id]`.
- `contact` (sendContact): `--chat --phone-number --first-name [--last-name] [--vcard]`.
- `poll` (sendPoll): `--chat --question --options '[{"text":"A"},{"text":"B"}]' [--type regular|quiz] [--anonymous] [--allows-multiple-answers] [--correct-option-id]`.
- `dice` (sendDice): `--chat [--emoji 🎲|🎯|🏀|⚽|🎳|🎰]`.

## media
- `photo` / `document` / `video` / `audio` / `voice` / `animation` / `video-note` / `sticker`:
  `--chat --<kind> <path|URL|file_id> [--caption] [--parse-mode] [--silent]` (audio adds
  `--performer --title --duration`). A local path is uploaded as multipart; an http(s) URL or a
  file_id is sent as a string.
- `media-group` (sendMediaGroup): `--chat --media '[{"type":"photo","media":"<url|file_id>"}]'` — 2-10 items as an album.

## file
- `info` (getFile): `--file-id <id>` → resolves to `file_path` and `file_size`.
- `download` (value-add): `--file-id <id> [--dest <path>|-]` — getFile + stream the bytes to disk (`-` for stdout).

## chat
- `get` (getChat), `members-count`, `administrators`, `member --user <id>`, `leave`.
- `set-title --title`, `set-description [--description]` (empty clears it).

## member (moderation; the bot must be an admin)
- `ban` (`--until`, `--revoke-messages`), `unban` (`--only-if-banned`),
  `restrict --permissions '<ChatPermissions JSON>'`, `promote --can-...` boolean rights.

## invite (invite links; the bot needs can_invite_users)
- `create` (createChatInviteLink): `--chat [--name] [--expire-date] [--member-limit] [--creates-join-request]`.
- `edit` (editChatInviteLink): `--chat --invite-link [--name] [--expire-date] [--member-limit]`.
- `revoke` (revokeChatInviteLink): `--chat --invite-link`.

## user
- `photos` (getUserProfilePhotos): `--user <id> [--offset] [--limit]`.

## callback / inline
- `callback answer` (answerCallbackQuery): `--callback-query-id [--text] [--show-alert] [--url] [--cache-time]`.
- `inline answer` (answerInlineQuery): `--inline-query-id --results '<InlineQueryResult[] JSON>' [--cache-time] [--is-personal] [--next-offset] [--button JSON]`.

## updates
- `get` (getUpdates): `--offset --limit --timeout --allowed-updates a,b`. Conflicts with a set
  webhook — delete it first to poll.

## webhook
- `info` (getWebhookInfo), `set --url [--secret-token] [--max-connections] [--drop-pending]`,
  `delete [--drop-pending]`.

## commands (the slash-command menu)
- `list` (getMyCommands), `set --commands '[{"command":"start","description":"Begin"}]'`,
  `delete`. All take optional `--scope '<BotCommandScope JSON>'` and `--language-code`.

## api (escape hatch)
- `tgctl api <method> -q key=value [-d '<JSON body>'] [--idempotent]` calls any Bot API method.
  Mark read-only methods `--idempotent` so transient failures retry safely.
