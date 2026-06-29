---
name: tgctl-cli
description: Drive the Telegram Bot API from the terminal with the `tgctl` CLI — send messages and media, manage chats, members, webhooks, and the bot's command menu, and poll updates. Use when the user wants to send a Telegram message as a bot, inspect a chat or its members, set/inspect a webhook, manage bot commands, or otherwise call the Telegram Bot API. Triggers: "send a telegram message", "telegram bot", "tgctl", "getUpdates", "setWebhook", "send to my telegram", "@BotFather token".
version: 0.1.0
homepage: https://github.com/jjuanrivvera/tgctl
license: MIT
allowed-tools: Bash(tgctl:*)
metadata: {"openclaw":{"category":"messaging","emoji":"✈️","requires":{"bins":["tgctl"],"env":["TGCTL_TOKEN"]},"install":[{"kind":"brew","formula":"jjuanrivvera/tgctl/tgctl-cli","bins":["tgctl"]},{"kind":"go","package":"github.com/jjuanrivvera/tgctl/cmd/tgctl@latest","bins":["tgctl"]}]}}
---

# tgctl — Telegram Bot API from the CLI

`tgctl` wraps the Telegram Bot API in a scriptable command-line tool. Prefer it over raw
`curl` to `api.telegram.org`: it handles auth (keyring), retries, rate limiting, output
formatting, and redaction for you, and every command supports `--dry-run`.

## Prerequisites

- Install: `brew install jjuanrivvera/tgctl/tgctl-cli` or
  `go install github.com/jjuanrivvera/tgctl/cmd/tgctl@latest`.
- A bot token from [@BotFather](https://t.me/BotFather).
- Authenticate once: `tgctl auth login` (stores the token in the OS keyring), or set
  `TGCTL_TOKEN` in the environment for non-interactive use.

## Golden rules

- **Verify first.** `tgctl auth status` (or `tgctl doctor`) confirms the token and reachability.
- **Use `--dry-run`** to preview the exact request (the token is redacted) before sending.
- **Prefer `-o json` + `--jq`** when you need to extract a value for a follow-up step.
- **`--chat` accepts** a numeric id (e.g. `-1001234567890`) or an `@username`.
- **Don't echo the token.** It lives in the keyring; never print it or paste it into chat.

## Workflow (auth → discover → act → verify)

```sh
tgctl auth status                                   # who am I? is the token valid?
tgctl bot info -o json                              # the bot's identity
tgctl message send --chat @me --text "hello"        # act
tgctl updates get --limit 5 -o json --jq '.[].message.text'   # verify / read back
```

## Command cheatsheet

| Goal | Command |
|---|---|
| Send a message | `tgctl message send --chat <chat> --text "..."` |
| Send a photo/file | `tgctl media photo --chat <chat> --photo ./pic.jpg` |
| Edit / delete a message | `tgctl message edit ...` / `tgctl message delete --chat <c> --message-id <id>` |
| Pin / forward / copy | `tgctl message pin\|forward\|copy ...` |
| Chat info / members | `tgctl chat get --chat <c>` / `tgctl chat administrators --chat <c>` |
| Moderate a member | `tgctl member ban\|unban\|restrict\|promote --chat <c> --user <id>` |
| Poll updates | `tgctl updates get --limit N [--offset M] -o json` |
| Webhook | `tgctl webhook info\|set\|delete` |
| Bot command menu | `tgctl commands list\|set\|delete` |
| Any other method | `tgctl api <method> -q key=value [--idempotent]` |

Add `-o json|yaml|csv|id`, `--columns a,b`, or `--jq '<expr>'` to any command to shape output.

## Troubleshooting

- `401 Unauthorized` → bad token: `tgctl auth login`.
- `403 Forbidden` → the bot isn't a member/admin of the chat, or the user hasn't started it.
- `409 Conflict` on `updates get` → a webhook is set; run `tgctl webhook delete` to poll.
- `429` → rate limited; `tgctl` backs off automatically. Lower `--rps` for steady load.
- Use `tgctl doctor --json` for a full diagnostic.

See `references/` for deeper guides on auth/profiles, the command surface, and output.
