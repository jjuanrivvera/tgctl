## tgctl api

Call any Bot API method directly (raw escape hatch)

### Synopsis

Invoke an arbitrary Bot API method with a JSON body and/or key=value parameters.

This is the documented escape hatch for methods tgctl does not wrap as first-class
commands. It honors --dry-run and -o/--output like every other command. By default a
raw call is treated as a write (not auto-retried); pass --idempotent for read-only
methods (getX) so transient failures retry safely.

Uploading files: some methods refuse anything but an uploaded file (setMyProfilePhoto,
setChatPhoto, uploadStickerFile, ...). Pass -F name=@path — repeatable — and the request
switches to multipart/form-data, sending each file as the part <name> while -d/-q travel
alongside as fields. That is what lets a JSON body reference an upload as "attach://<name>":

  tgctl api setMyProfilePhoto \
    -d '{"photo":{"type":"static","photo":"attach://pic"}}' -F pic=@logo.jpg

```
tgctl api <method> [-d body] [-q key=value ...] [-F name=@file ...] [flags]
```

### Examples

```
  tgctl api getMe
  tgctl api sendMessage -q chat_id=@me -q text="hi from tgctl"
  tgctl api getChat -q chat_id=@telegram --idempotent
  tgctl api sendMessage -d '{"chat_id":"@me","text":"json body"}'
  tgctl api setChatPhoto -q chat_id=@mygroup -F photo=@cover.jpg
  tgctl api setMyProfilePhoto -d '{"photo":{"type":"static","photo":"attach://pic"}}' -F pic=@logo.jpg
```

### Options

```
  -d, --data string         raw JSON request body
  -F, --file stringArray    name=@path to upload as multipart (repeatable)
  -h, --help                help for api
      --idempotent          treat as read-only (safe to auto-retry)
  -q, --query stringArray   key=value parameter (repeatable)
```

### Options inherited from parent commands

```
      --base-url string   Bot API base URL (default https://api.telegram.org)
      --bot string        bot to use: a named profile/credential (env TGCTL_BOT)
      --columns strings   explicit, ordered table/csv columns
      --dry-run           print the equivalent curl and make no request
      --jq string         gojq expression applied to the result before rendering
      --no-color          disable colored output
      --no-store          disable local SQLite send/receive history for this invocation (see tgctl log)
  -o, --output string     output format: table|json|yaml|csv|id (default "table")
      --quiet             suppress notes on stderr
      --rps float         client-side requests-per-second cap (0 = default)
      --show-token        do not redact the bot token in --dry-run output
  -v, --verbose           log raw API responses to stderr
```

### SEE ALSO

* [tgctl](tgctl.md)	 - Command-line tool for the Telegram Bot API

