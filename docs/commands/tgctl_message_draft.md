## tgctl message draft

Stream a partial message while it is still being written

### Synopsis

Show a user a live, partial message while the real one is still being generated
(sendMessageDraft) — the "typing out an answer" effect, for a private chat.

The draft is EPHEMERAL: it is a preview that disappears after about 30 seconds, so once the
text is final you still have to send it with `message send` for it to exist in the
chat. Reusing the same --draft-id animates the change; a new id replaces it outright. An empty
--text shows a "Thinking…" placeholder.

```
tgctl message draft [flags]
```

### Examples

```
  tgctl message draft --chat 12345 --draft-id 1 --text "Looking that up"
  tgctl message draft --chat 12345 --draft-id 1 --text "" --can-stop
```

### Options

```
      --can-stop                show the user a button to stop further drafts
      --chat int                target private chat id (drafts are private-chat only)
      --draft-id int            id of this draft; reuse it to animate the change (non-zero)
      --entities string         JSON array of MessageEntity objects (instead of --parse-mode)
  -h, --help                    help for draft
      --keep-on-stop            keep the draft visible when the user presses stop
      --message-thread-id int   target message thread
      --parse-mode string       MarkdownV2 | HTML | Markdown
      --text string             partial text so far (empty shows a "Thinking…" placeholder)
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

