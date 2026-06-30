# Auth & configuration

## Token storage
- `tgctl auth login` verifies the token against `getMe` and stores it in the OS keyring.
- Non-interactive: `tgctl auth login --token "<bot_id>:<hash>"` or set `TGCTL_TOKEN`.
- Token precedence: `$TGCTL_TOKEN` > `$TELEGRAM_BOT_TOKEN` > the active profile's keyring entry.
- The config file (`~/.tgctl-cli/config.yaml`, or `$XDG_CONFIG_HOME/tgctl/config.yaml`) stores
  only non-secret data (base URL, bot id, auth method). Permissions are `0600`.

## Bots / profiles (`--bot`)
A profile is one bot, so select it with `--bot`. The old `--profile` flag still works as a
hidden alias, so existing scripts don't break.
```sh
tgctl auth login --bot staging          # add a second bot
tgctl --bot staging bot info            # one-off use
tgctl config use staging                # set the default
tgctl config list-profiles              # see them all
tgctl config set base_url http://localhost:8081   # self-hosted Local Bot API Server
```
Selection precedence: `--bot` > `$TGCTL_BOT` > `$TGCTL_PROFILE` (legacy) > the active bot >
`default`.

## Self-hosted Bot API server
Point `--base-url` (or the profile's `base_url`) at a
[Local Bot API Server](https://github.com/tdlib/telegram-bot-api). Cleartext `http://` is
rejected for non-loopback hosts to avoid leaking the token.

## Headless hosts
With no OS keyring, tokens fall back to an AES-256-GCM encrypted file. Set
`TGCTL_KEYRING_PASSWORD` to derive the key from a passphrase instead of the host-bound default.
