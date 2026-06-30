## tgctl file

Inspect and download files

### Synopsis

Resolve a file_id to its metadata (getFile) and download the file's bytes to disk.

### Options

```
  -h, --help   help for file
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API
* [tgctl file download](tgctl_file_download.md)	 - Download a file by file_id to a local path
* [tgctl file info](tgctl_file_info.md)	 - Resolve a file_id to its path and size (getFile)

