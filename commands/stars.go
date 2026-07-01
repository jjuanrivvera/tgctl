package commands

func init() {
	registerGroup(group{
		Use:   "stars",
		Short: "Telegram Stars: transactions, gifts, and paid media",
		Long:  "Inspect the bot's Star balance transactions, send gifts, refund Star payments, manage Star subscriptions, and send paid media.",
		Cmds: []methodCmd{
			{
				Use: "transactions", Aliases: []string{"txns"}, Method: "getStarTransactions", Kind: kindRead,
				Short:   "List the bot's Star transactions",
				Example: `  tgctl stars transactions --limit 20`,
				Flags: []flagSpec{
					{Name: "offset", Kind: flagInt, Usage: "number of transactions to skip"},
					{Name: "limit", Kind: flagInt, Usage: "max transactions to return (1-100)"},
				},
			},
			{
				Use: "gifts", Method: "getAvailableGifts", Kind: kindRead,
				Short:   "List the gifts the bot can send",
				Example: `  tgctl stars gifts -o json`,
			},
			{
				Use: "send-gift", Method: "sendGift", Kind: kindWrite,
				Short:   "Send a gift to a user or channel",
				Example: `  tgctl stars send-gift --user 12345 --gift-id 5170233102089322756 --text "Enjoy!"`,
				Flags: []flagSpec{
					{Name: "user", Param: "user_id", Kind: flagInt, Usage: "recipient user id (or use --chat)"},
					{Name: "chat", Param: "chat_id", Usage: "recipient channel chat id or @username (or use --user)"},
					{Name: "gift-id", Param: "gift_id", Required: true, Usage: "id of the gift to send (from stars gifts)"},
					{Name: "pay-for-upgrade", Param: "pay_for_upgrade", Kind: flagBool, Usage: "pay for the gift's upgrade to a unique gift"},
					{Name: "text", Usage: "text shown with the gift (0-128 chars)"},
					{Name: "text-parse-mode", Param: "text_parse_mode", Usage: "parse mode for --text (MarkdownV2 | HTML)"},
				},
			},
			{
				Use: "refund", Method: "refundStarPayment", Kind: kindWrite,
				Short:   "Refund a successful Star payment",
				Example: `  tgctl stars refund --user 12345 --charge-id abc123`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "charge-id", Param: "telegram_payment_charge_id", Required: true, Usage: "the telegram payment charge id to refund"},
				},
			},
			{
				Use: "edit-subscription", Method: "editUserStarSubscription", Kind: kindWrite,
				Short:   "Cancel or re-enable a user's Star subscription",
				Example: `  tgctl stars edit-subscription --user 12345 --charge-id abc123 --canceled`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "charge-id", Param: "telegram_payment_charge_id", Required: true, Usage: "the telegram payment charge id of the subscription"},
					{Name: "canceled", Param: "is_canceled", Kind: flagBool, Required: true, Usage: "true to cancel, false to re-enable before the period ends"},
				},
			},
			{
				Use: "set-emoji-status", Method: "setUserEmojiStatus", Kind: kindWrite,
				Short:   "Set a user's emoji status (requires the user's prior consent)",
				Example: `  tgctl stars set-emoji-status --user 12345 --emoji-status-custom-emoji-id 5170233102089322756`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "emoji-status-custom-emoji-id", Param: "emoji_status_custom_emoji_id", Usage: "custom emoji id for the status (omit to remove)"},
					{Name: "emoji-status-expiration-date", Param: "emoji_status_expiration_date", Kind: flagInt, Usage: "unix time the status expires"},
				},
			},
			{
				Use: "send-paid-media", Method: "sendPaidMedia", Kind: kindWrite,
				Short: "Send paid media that recipients unlock with Stars",
				Long:  "Send media locked behind a Star paywall. --media is a JSON array of InputPaidMedia objects.",
				Example: `  tgctl stars send-paid-media --chat @channel --star-count 50 \
    --media '[{"type":"photo","media":"https://e.com/a.jpg"}]'`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "star-count", Param: "star_count", Kind: flagInt, Required: true, Usage: "Stars a user must pay to unlock (1-2500)"},
					{Name: "media", Kind: flagJSON, Required: true, Usage: "JSON array of InputPaidMedia objects"},
					{Name: "payload", Usage: "bot-defined payload (not shown to users)"},
					{Name: "caption", Usage: "media caption (0-1024 chars)"},
					parseModeFlag(), silentFlag(), protectContentFlag(), businessFlag(),
				},
				Columns: []string{"message_id", "chat.id"},
			},
		},
	})
}
