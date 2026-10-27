package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "business",
		Short: "Manage Telegram Business connections and accounts",
		Long:  "Inspect business connections, manage business account profiles, read/delete messages, and manage Stars and gifts.",
		Cmds: []methodCmd{
			{
				Use: "connection", Method: "getBusinessConnection", Kind: kindRead,
				Short:   "Get information about a business connection",
				Long:    "Get information about the connection of the bot with a business account.",
				Example: "  tgctl business connection --business-connection-id bc123",
				Flags: []flagSpec{
					businessConnectionFlag(),
				},
			},
			{
				Use: "gifts", Method: "getBusinessAccountGifts", Kind: kindRead,
				Short:   "Get gifts received by a business account",
				Long:    "Get the list of gifts received by the business account.",
				Example: "  tgctl business gifts --business-connection-id bc123\n  tgctl business gifts --business-connection-id bc123 --exclude-unlimited --limit 50",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "exclude-unsaved", Param: "exclude_unsaved", Kind: flagBool, Usage: "exclude gifts that are not saved to the profile"},
					{Name: "exclude-saved", Param: "exclude_saved", Kind: flagBool, Usage: "exclude gifts that are saved to the profile"},
					{Name: "exclude-unlimited", Param: "exclude_unlimited", Kind: flagBool, Usage: "exclude unlimited gifts"},
					{Name: "exclude-limited", Param: "exclude_limited", Kind: flagBool, Usage: "exclude limited gifts"},
					{Name: "exclude-unique", Param: "exclude_unique", Kind: flagBool, Usage: "exclude unique gifts"},
					{Name: "sort-by-price", Param: "sort_by_price", Kind: flagBool, Usage: "sort results by gift price instead of send date"},
					{Name: "offset", Usage: "offset of the first gift to return from previous request"},
					{Name: "limit", Kind: flagInt, Usage: "number of gifts to return (1-100)"},
				},
				PreCall: rangeCheck("limit", "limit", 1, 100),
			},
			{
				Use: "star-balance", Method: "getBusinessAccountStarBalance", Kind: kindRead,
				Short:   "Get Star balance of a business account",
				Long:    "Get the Star balance of the business account.",
				Example: "  tgctl business star-balance --business-connection-id bc123",
				Flags: []flagSpec{
					businessConnectionFlag(),
				},
			},
			{
				Use: "read-message", Method: "readBusinessMessage", Kind: kindWrite,
				Short:   "Mark a message from a connected business account as read",
				Long:    "Mark a message from a connected business account as read on behalf of the business account.",
				Example: "  tgctl business read-message --business-connection-id bc123 --chat 12345678 --message-id 42",
				Flags: []flagSpec{
					businessConnectionFlag(),
					chatFlag(),
					messageIDFlag(),
				},
			},
			{
				Use: "delete-messages", Method: "deleteBusinessMessages", Kind: kindDestructive,
				Short:   "Delete messages on behalf of a connected business account",
				Long:    "Delete 1 to 100 messages on behalf of a connected business account.",
				Example: "  tgctl business delete-messages --business-connection-id bc123 --message-ids '[10,11,12]'",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "message-ids", Param: "message_ids", Kind: flagJSON, Required: true, Usage: "JSON array of 1-100 message identifiers to delete"},
				},
				PreCall: requireBusinessMessageIDs,
			},
			{
				Use: "set-bio", Method: "setBusinessAccountBio", Kind: kindWrite,
				Short:   "Set bio of a connected business account",
				Long:    "Change the bio of the business account. An empty or omitted bio removes it (0-140 characters).",
				Example: "  tgctl business set-bio --business-connection-id bc123 --bio \"Welcome to our shop!\"\n  tgctl business set-bio --business-connection-id bc123 --bio \"\"",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "bio", Usage: "new business account bio (0-140 characters); empty to remove"},
				},
				PreCall: requireValidBusinessBio,
			},
			{
				Use: "set-name", Method: "setBusinessAccountName", Kind: kindWrite,
				Short:   "Set name of a connected business account",
				Long:    "Change the first and last name of the business account.",
				Example: "  tgctl business set-name --business-connection-id bc123 --first-name \"Acme Corp\"\n  tgctl business set-name --business-connection-id bc123 --first-name \"Acme\" --last-name \"HQ\"",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "first-name", Param: "first_name", Required: true, Usage: "new first name of the business account (1-64 characters)"},
					{Name: "last-name", Param: "last_name", Usage: "new last name of the business account (0-64 characters)"},
				},
				PreCall: requireValidBusinessName,
			},
			{
				Use: "set-username", Method: "setBusinessAccountUsername", Kind: kindWrite,
				Short:   "Set username of a connected business account",
				Long:    "Change the username of the business account. An empty username removes it (5-32 characters).",
				Example: "  tgctl business set-username --business-connection-id bc123 --username acmeshop\n  tgctl business set-username --business-connection-id bc123 --username \"\"",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "username", Usage: "new business account username (5-32 characters); empty to remove"},
				},
				PreCall: requireValidBusinessUsername,
			},
			{
				Use: "set-gift-settings", Method: "setBusinessAccountGiftSettings", Kind: kindWrite,
				Short:   "Set gift settings of a connected business account",
				Long:    "Change the gift settings of a connected business account.",
				Example: "  tgctl business set-gift-settings --business-connection-id bc123 --show-gift-button --accepted-gift-types '{\"unlimited_gifts\":true}'",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "show-gift-button", Param: "show_gift_button", Kind: flagBool, Required: true, Usage: "show the gift button on the profile of the business account"},
					{Name: "accepted-gift-types", Param: "accepted_gift_types", Kind: flagJSON, Required: true, Usage: "JSON object describing accepted gift types"},
				},
				PreCall: requireValidGiftSettings,
			},
			{
				Use: "remove-photo", Method: "removeBusinessAccountProfilePhoto", Kind: kindDestructive,
				Short:   "Remove profile photo of a connected business account",
				Long:    "Remove the profile photo of the connected business account.",
				Example: "  tgctl business remove-photo --business-connection-id bc123",
				Flags: []flagSpec{
					businessConnectionFlag(),
				},
			},
			{
				Use: "transfer-stars", Method: "transferBusinessAccountStars", Kind: kindDestructive,
				Short:   "Transfer Stars from a business account to the bot",
				Long:    "Transfer Telegram Stars from the connected business account balance to the bot's balance.",
				Example: "  tgctl business transfer-stars --business-connection-id bc123 --star-count 500",
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "star-count", Param: "star_count", Kind: flagInt, Required: true, Usage: "number of Stars to transfer (positive integer)"},
				},
				PreCall: requirePositiveStarCount,
			},
		},
	})
}

func requireBusinessMessageIDs(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	v, ok := params["message_ids"]
	if !ok {
		return nil
	}
	slice, ok := v.([]any)
	if !ok {
		return fmt.Errorf("--message-ids must be a JSON array of message ids")
	}
	if len(slice) < 1 || len(slice) > 100 {
		return fmt.Errorf("--message-ids must contain between 1 and 100 message ids, got %d", len(slice))
	}
	return nil
}

func requireValidBusinessBio(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	bio, ok := params["bio"].(string)
	if !ok {
		return nil
	}
	if len([]rune(bio)) > 140 {
		return fmt.Errorf("--bio must be at most 140 characters, got %d", len([]rune(bio)))
	}
	return nil
}

func requireValidBusinessName(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	fn, _ := params["first_name"].(string)
	if len([]rune(fn)) < 1 || len([]rune(fn)) > 64 {
		return fmt.Errorf("--first-name must be between 1 and 64 characters, got %d", len([]rune(fn)))
	}
	if ln, ok := params["last_name"].(string); ok && len([]rune(ln)) > 64 {
		return fmt.Errorf("--last-name must be at most 64 characters, got %d", len([]rune(ln)))
	}
	return nil
}

func requireValidBusinessUsername(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	u, ok := params["username"].(string)
	if !ok || u == "" {
		return nil
	}
	if len([]rune(u)) < 5 || len([]rune(u)) > 32 {
		return fmt.Errorf("--username must be between 5 and 32 characters or empty to remove, got %d", len([]rune(u)))
	}
	return nil
}

func requireValidGiftSettings(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	v, ok := params["accepted_gift_types"]
	if !ok {
		return nil
	}
	if _, isObj := v.(map[string]any); !isObj {
		return fmt.Errorf("--accepted-gift-types must be a JSON object")
	}
	return nil
}

func requirePositiveStarCount(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	v, ok := params["star_count"].(int64)
	if !ok || v <= 0 {
		return fmt.Errorf("--star-count must be greater than zero")
	}
	return nil
}
