# Output & filtering

One renderer serves every command. Select with `-o/--output`:

| Format | Notes |
|---|---|
| `table` | default; aligned, colored on a TTY (honors `NO_COLOR` and `--no-color`) |
| `json` | full record; large ids stay exact (no float rounding) |
| `yaml` | readable structured output |
| `csv` | spreadsheet-friendly; cells are sanitized against formula injection |
| `id` | one id per line — pipe to `xargs` |

## Selecting and filtering
- `--columns message_id,chat.id,text` selects and orders columns (nested fields use dotted keys).
- `--jq '<expr>'` runs a built-in gojq program over the result before rendering:
  ```sh
  tgctl updates get -o json --jq '.[].message | {id: .message_id, text}'
  tgctl chat administrators --chat @g -o json --jq '.[].user.username'
  ```
- Notes/warnings go to **stderr**, so stdout stays clean for pipes.

## Scripting tips
- `--dry-run` prints the equivalent `curl` (token redacted) — great for debugging.
- `-o id` + `xargs` for bulk follow-ups.
- `tgctl doctor --json` for machine-readable health checks in CI.
- Exit codes are non-zero on failure (e.g. `auth status`, `doctor`), so `tgctl auth status && ...`
  works in scripts.
