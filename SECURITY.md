# Security Policy

## Supported versions

The latest released minor version receives security fixes. Older versions are best-effort.

| Version | Supported |
|---|---|
| latest | ✅ |
| older  | ⚠️ best-effort |

## Reporting a vulnerability

Please report security issues privately via GitHub Security Advisories
("Report a vulnerability" on the repository's **Security** tab), not in a public issue.
Include reproduction steps and the affected version. We aim to acknowledge within a few days.

## Token handling

`tgctl` is built so a bot token never lands in plaintext where it shouldn't:

- Tokens are stored in the **OS keyring** (macOS Keychain, Linux Secret Service, Windows
  Credential Manager). When no keyring is available, the fallback is an **AES-256-GCM
  encrypted file** (`credentials.enc`, `0600`); set `TGCTL_KEYRING_PASSWORD` to derive the
  key from a passphrase rather than the host-bound default.
- The config file (`~/.tgctl-cli/config.yaml`, `0600`) stores only non-secret data (base URL,
  bot id, chosen auth method) — **never** the token.
- `--dry-run` and `config view` **redact** the token by default. Use `--show-token` only when
  you explicitly want the secret in the output.
- The MCP tool surface excludes the token/profile/base-url flags and the raw `api` command,
  so an AI agent driving `tgctl` cannot read the credential or switch instances.
- CSV output is sanitized against spreadsheet formula injection (CWE-1236).

If you find a case where a token can leak (logs, error messages, generated config, a dry-run
that isn't redacted), please treat it as a security issue and report it privately.
