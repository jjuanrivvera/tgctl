## tgctl message

Send and manage messages

### Synopsis

Send, edit, delete, forward, copy, and pin messages. --chat accepts a numeric id or @username.

### Options

```
  -h, --help   help for message
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
* [tgctl message contact](tgctl_message_contact.md)	 - Send a phone contact
* [tgctl message copy](tgctl_message_copy.md)	 - Copy a message (without a 'forwarded from' header)
* [tgctl message delete](tgctl_message_delete.md)	 - Delete a message
* [tgctl message dice](tgctl_message_dice.md)	 - Send an animated emoji with a random value (dice, dart, etc.)
* [tgctl message edit](tgctl_message_edit.md)	 - Edit a message's text
* [tgctl message forward](tgctl_message_forward.md)	 - Forward a message to another chat
* [tgctl message location](tgctl_message_location.md)	 - Send a point on the map
* [tgctl message pin](tgctl_message_pin.md)	 - Pin a message in a chat
* [tgctl message poll](tgctl_message_poll.md)	 - Send a native poll
* [tgctl message react](tgctl_message_react.md)	 - Set (or clear) reactions on a message
* [tgctl message send](tgctl_message_send.md)	 - Send a text message
* [tgctl message unpin](tgctl_message_unpin.md)	 - Unpin a message (or the most recent pin) in a chat
* [tgctl message venue](tgctl_message_venue.md)	 - Send information about a venue

