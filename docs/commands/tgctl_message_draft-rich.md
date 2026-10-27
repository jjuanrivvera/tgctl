## tgctl message draft-rich

Stream a partial rich message while it is being written

### Synopsis

Stream a partial rich message to a user while it is still being generated (sendRichMessageDraft).
The draft is EPHEMERAL: it is a preview that disappears after about 30 seconds.
Once generation is complete, send the finalized message with message send-rich to persist it in chat.
Reusing the same --draft-id animates the change. Direct file uploads are not supported for drafts.

```
tgctl message draft-rich [flags]
```

### Examples

```
  tgctl message draft-rich --chat 12345 --draft-id 1 --rich-message '{"blocks":[{"type":"thinking","text":"Generating report..."}]}' --can-stop
  tgctl message draft-rich --chat 12345 --draft-id 1 --rich-message '{"blocks":[{"type":"paragraph","text":"Partial answer"}]}'
```

### Options

```
      --can-stop                show the user a button to stop further drafts
      --chat int                target private chat id (drafts are private-chat only)
      --draft-id int            id of this draft; reuse it to animate the change (non-zero)
  -h, --help                    help for draft-rich
      --keep-on-stop            keep the draft visible when the user presses stop
      --message-thread-id int   target message thread
      --rich-message string     JSON object describing the rich message draft (InputRichMessage)
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

* [tgctl message](tgctl_message.md)	 - Send and manage messages

