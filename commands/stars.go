package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "stars",
		Short: "Telegram Stars: transactions, gifts, and paid media",
		Long:  "Inspect the bot's Star balance and transactions, send and manage gifts, refund Star payments, manage Star subscriptions, gift Premium, and send paid media.",
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
				Use: "balance", Method: "getMyStarBalance", Kind: kindRead,
				Short:   "Show the bot's current Star balance",
				Example: "  tgctl stars balance\n  tgctl stars balance -o json",
			},
			{
				Use: "user-gifts", Method: "getUserGifts", Kind: kindRead,
				Short: "List the gifts a user owns",
				Long: `List the gifts owned and hosted by a user (getUserGifts). The --exclude-* flags
narrow the list by gift kind; --sort-by-price orders by value instead of by date.`,
				Example: `  tgctl stars user-gifts --user 12345
  tgctl stars user-gifts --user 12345 --exclude-unlimited --sort-by-price`,
				Flags: append([]flagSpec{
					{Name: "user", Param: "user_id", Kind: flagInt, Required: true, Usage: "owner of the gifts"},
				}, giftFilterFlags()...),
			},
			{
				Use: "chat-gifts", Method: "getChatGifts", Kind: kindRead,
				Short: "List the gifts a chat owns",
				Long: `List the gifts owned by a channel (getChatGifts). Whether unsaved gifts are visible
depends on the bot holding can_post_messages in that channel.`,
				Example: `  tgctl stars chat-gifts --chat @mychannel`,
				Flags: append([]flagSpec{
					chatFlag(),
					{Name: "exclude-unsaved", Param: "exclude_unsaved", Kind: flagBool, Usage: "skip gifts not saved to the profile page"},
					{Name: "exclude-saved", Param: "exclude_saved", Kind: flagBool, Usage: "skip gifts saved to the profile page"},
				}, giftFilterFlags()...),
			},
			{
				// Destructive: the gift stops being a gift. Converting is one-way — the stars
				// come back, the item does not.
				Use: "convert-gift", Method: "convertGiftToStars", Kind: kindDestructive,
				Short:   "Convert a regular gift into Telegram Stars",
				Long:    "Convert an owned regular gift into Stars (convertGiftToStars). One-way: the gift is consumed. Needs the can_convert_gifts_to_stars business right.",
				Example: `  tgctl stars convert-gift --business-connection-id BQ... --gift-id abc123`,
				Flags:   []flagSpec{businessConnectionFlag(), ownedGiftFlag()},
			},
			{
				// Destructive: upgrading replaces the regular gift with a unique one and may
				// spend stars doing it; neither step can be undone.
				Use: "upgrade-gift", Method: "upgradeGift", Kind: kindDestructive,
				Short:   "Upgrade a regular gift to a unique one",
				Long:    "Upgrade an owned regular gift to a unique gift (upgradeGift). Irreversible, and --star-count pays for it when the upgrade is not free. Needs the can_transfer_and_upgrade_gifts business right.",
				Example: `  tgctl stars upgrade-gift --business-connection-id BQ... --gift-id abc123 --keep-original-details`,
				Flags: []flagSpec{
					businessConnectionFlag(), ownedGiftFlag(),
					{Name: "keep-original-details", Param: "keep_original_details", Kind: flagBool, Usage: "keep the original text, sender and receiver"},
					{Name: "star-count", Param: "star_count", Kind: flagInt, Usage: "stars to pay when the upgrade is paid"},
				},
			},
			{
				// Destructive: the gift leaves this account for another owner.
				Use: "transfer-gift", Method: "transferGift", Kind: kindDestructive,
				Short:   "Transfer a unique gift to another owner",
				Long:    "Transfer an owned unique gift to another chat (transferGift). The gift leaves this account. The receiving chat must have been active in the last 24 hours.",
				Example: `  tgctl stars transfer-gift --business-connection-id BQ... --gift-id abc123 --new-owner 12345`,
				Flags: []flagSpec{
					businessConnectionFlag(), ownedGiftFlag(),
					{Name: "new-owner", Param: "new_owner_chat_id", Kind: flagInt, Required: true, Usage: "chat id that will own the gift"},
					{Name: "star-count", Param: "star_count", Kind: flagInt, Usage: "stars to pay when the transfer is paid"},
				},
			},
			{
				// Destructive: it spends the bot's stars, and a subscription cannot be
				// un-gifted. The agent guard should treat spending like any other one-way act.
				Use: "gift-premium", Method: "giftPremiumSubscription", Kind: kindDestructive,
				Short: "Gift a user a Telegram Premium subscription",
				Long: `Gift Telegram Premium to a user (giftPremiumSubscription). It SPENDS the bot's
Stars and cannot be undone. The API fixes the price: 3 months costs 1000 stars, 6 months 1500,
and 12 months 2500 — any other pairing is rejected before the request leaves.`,
				Example: `  tgctl stars gift-premium --user 12345 --months 3 --stars 1000 --text "Happy birthday"`,
				Flags: []flagSpec{
					{Name: "user", Param: "user_id", Kind: flagInt, Required: true, Usage: "who receives the subscription"},
					{Name: "months", Param: "month_count", Kind: flagInt, Required: true, Usage: "3, 6 or 12"},
					{Name: "stars", Param: "star_count", Kind: flagInt, Required: true, Usage: "1000 for 3 months, 1500 for 6, 2500 for 12"},
					{Name: "text", Usage: "text shown with the service message"},
					{Name: "text-parse-mode", Param: "text_parse_mode", Usage: "parse mode for --text"},
					{Name: "text-entities", Param: "text_entities", Kind: flagJSON, Usage: "JSON array of MessageEntity (instead of --text-parse-mode)"},
				},
				PreCall: requirePremiumPricing,
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
				// Destructive: a refund is an irreversible money movement — the payment
				// cannot be re-charged once returned.
				Use: "refund", Method: "refundStarPayment", Kind: kindDestructive,
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

// giftFilterFlags are the --exclude-* selectors both gift listings share. They are identical
// in name and meaning on getUserGifts and getChatGifts, so they live here once.
func giftFilterFlags() []flagSpec {
	return []flagSpec{
		{Name: "exclude-unlimited", Param: "exclude_unlimited", Kind: flagBool, Usage: "skip gifts sold in unlimited numbers"},
		{Name: "exclude-limited-upgradable", Param: "exclude_limited_upgradable", Kind: flagBool, Usage: "skip limited gifts that can become unique"},
		{Name: "exclude-limited-non-upgradable", Param: "exclude_limited_non_upgradable", Kind: flagBool, Usage: "skip limited gifts that cannot become unique"},
		{Name: "exclude-from-blockchain", Param: "exclude_from_blockchain", Kind: flagBool, Usage: "skip gifts held on the blockchain"},
		{Name: "exclude-unique", Param: "exclude_unique", Kind: flagBool, Usage: "skip unique gifts"},
		{Name: "sort-by-price", Param: "sort_by_price", Kind: flagBool, Usage: "order by value instead of by date"},
		{Name: "offset", Usage: "opaque cursor from a previous page"},
		{Name: "limit", Kind: flagInt, Usage: "max gifts to return (1-100)"},
	}
}

func businessConnectionFlag() flagSpec {
	return flagSpec{Name: "business-connection-id", Param: "business_connection_id", Required: true, Usage: "id of the business connection acting on the account"}
}

func ownedGiftFlag() flagSpec {
	return flagSpec{Name: "gift-id", Param: "owned_gift_id", Required: true, Usage: "id of the owned gift (from stars user-gifts)"}
}

// premiumPricing is the price list the API fixes for a gifted subscription.
var premiumPricing = map[int64]int64{3: 1000, 6: 1500, 12: 2500}

// requirePremiumPricing refuses a months/stars pairing the API does not sell. This one is
// worth catching locally beyond the usual round-trip argument: the command spends real Stars,
// and a mistyped --stars is the kind of error someone would rather meet before it is charged
// than after.
func requirePremiumPricing(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	months, _ := params["month_count"].(int64)
	stars, _ := params["star_count"].(int64)
	want, ok := premiumPricing[months]
	if !ok {
		return fmt.Errorf("--months must be 3, 6 or 12, got %d", months)
	}
	if stars != want {
		return fmt.Errorf("--months %d costs %d stars, not %d", months, want, stars)
	}
	return nil
}
