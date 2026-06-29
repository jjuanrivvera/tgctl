# Command surface

Groups map 1:1 to Bot API methods. Run `tgctl <group> <verb> --help` for flags and examples.

## bot
- `bot info` (getMe), `bot set-name`/`get-name`, `bot set-description`/`get-description`.

## message
- `send` (sendMessage): `--chat --text [--parse-mode MarkdownV2|HTML] [--reply-to] [--silent] [--no-preview] [--reply-markup JSON]`.
- `edit` (editMessageText), `delete` (deleteMessage), `forward`, `copy`, `pin`, `unpin`.

## media
- `photo` / `document` / `video`: `--chat --<kind> <path|URL|file_id> [--caption] [--parse-mode] [--silent]`.
  A local path is uploaded as multipart; an http(s) URL or a file_id is sent as a string.

## chat
- `get` (getChat), `members-count`, `administrators`, `member --user <id>`, `leave`.

## member (moderation; the bot must be an admin)
- `ban` (`--until`, `--revoke-messages`), `unban` (`--only-if-banned`),
  `restrict --permissions '<ChatPermissions JSON>'`, `promote --can-...` boolean rights.

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
