## tgctl file download

Download a file by file_id to a local path

### Synopsis

Resolve a file_id with getFile and download its bytes. The destination defaults to the
file's base name in the current directory; pass --dest to choose a path, or --dest - to write
to stdout. Honors --dry-run (it prints the getFile request without downloading).

```
tgctl file download [flags]
```

### Examples

```
  tgctl file download --file-id BAADBAADrwAD...
  tgctl file download --file-id BAADBAADrwAD... --dest ./photo.jpg
  tgctl file download --file-id BAADBAADrwAD... --dest - > out.bin
```

### Options

```
      --dest string      destination path (default: the file's base name; '-' for stdout)
      --file-id string   file_id to download (from a message's document/photo/etc.)
  -h, --help             help for download
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

* [tgctl file](tgctl_file.md)	 - Inspect and download files

