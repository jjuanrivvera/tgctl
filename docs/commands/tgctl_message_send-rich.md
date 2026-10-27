## tgctl message send-rich

Send a rich structured message

### Synopsis

Send a rich structured message with block-level layout, media, tables, and lists (sendRichMessage).
--rich-message accepts a JSON object (InputRichMessage).
Pass --file name=@path (repeatable) to upload multipart parts referenced by attach://name in blocks.

```
tgctl message send-rich [flags]
```

### Examples

```
  tgctl message send-rich --chat @me --rich-message '{"blocks":[{"type":"section_heading","text":"Hello"},{"type":"paragraph","text":"World"}]}'
  tgctl message send-rich --chat @me --rich-message '{"blocks":[{"type":"photo","photo":"attach://pic"}]}' --file pic=@./image.jpg
```

### Options

```
      --allow-paid-broadcast                  allow up to 1000 messages per second by paying Telegram Stars
      --business-connection-id string         business connection id on whose behalf to act
      --chat string                           target chat: numeric id or @username
      --direct-messages-topic-id int          target direct messages topic id
      --ephemeral-message-parameters string   JSON object for ephemeral message parameters
  -F, --file stringArray                      name=@path to upload as multipart (repeatable)
  -h, --help                                  help for send-rich
      --message-effect-id string              unique identifier of the message effect
      --protect-content                       protect the content from forwarding and saving
      --reply-markup string                   inline/reply keyboard as JSON
      --reply-parameters string               JSON object describing message to reply to
      --rich-message string                   JSON object describing the rich message (InputRichMessage)
      --silent                                send without a notification sound
      --suggested-post-parameters string      JSON object for suggested post parameters
      --thread int                            forum topic thread id
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

