package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "inline",
		Short: "Answer inline and guest queries",
		Long:  "Respond to an inline query (a user typing @yourbot ...) with a list of results, or to a guest message with a single result.",
		Cmds: []methodCmd{
			{
				Use: "answer", Method: "answerInlineQuery", Kind: kindWrite,
				Short: "Answer an inline query with results",
				Long:  "Answer an inline query. --results is a JSON array of InlineQueryResult objects.",
				Example: `  tgctl inline answer --inline-query-id 999 \
    --results '[{"type":"article","id":"1","title":"Hi","input_message_content":{"message_text":"Hi"}}]'`,
				Flags: []flagSpec{
					{Name: "inline-query-id", Param: "inline_query_id", Required: true, Usage: "id of the inline query to answer"},
					{Name: "results", Kind: flagJSON, Required: true, Usage: "JSON array of InlineQueryResult objects (max 50)"},
					{Name: "cache-time", Param: "cache_time", Kind: flagInt, Usage: "seconds the result may be cached server-side"},
					{Name: "is-personal", Param: "is_personal", Kind: flagBool, Usage: "cache results per-user instead of globally"},
					{Name: "next-offset", Param: "next_offset", Usage: "offset a client sends to request the next page"},
					{Name: "button", Kind: flagJSON, Usage: "InlineQueryResultsButton object as JSON"},
				},
			},
			{
				Use: "answer-guest", Method: "answerGuestQuery", Kind: kindWrite,
				Short: "Reply to a guest message with a single result",
				Long: `Reply to a guest message (answerGuestQuery).

Unlike ` + "`inline answer`" + `, which returns a LIST the user picks from, this sends ONE
result as the reply: --result is a single InlineQueryResult object, not an array.`,
				Example: `  tgctl inline answer-guest --query-id AAxx     --result '{"type":"article","id":"1","title":"Hi","input_message_content":{"message_text":"Hi"}}'`,
				Flags: []flagSpec{
					{Name: "query-id", Param: "guest_query_id", Required: true, Usage: "id of the guest query to answer"},
					{Name: "result", Kind: flagJSON, Required: true, Usage: "a single InlineQueryResult object as JSON"},
				},
				PreCall: requireResultObject,
			},
			{
				Use: "answer-web-app", Method: "answerWebAppQuery", Kind: kindWrite,
				Short: "Answer a Web App query with a message",
				Long:  "Send one InlineQueryResult on behalf of the user to the originating chat. Returns a SentWebAppMessage.",
				Flags: []flagSpec{
					{Name: "query-id", Param: "web_app_query_id", Required: true, Usage: "Web App query identifier"},
					{Name: "result", Kind: flagJSON, Required: true, Usage: "single InlineQueryResult object as JSON"},
				}, PreCall: requireResultObject,
			},
			{
				Use: "prepare-message", Method: "savePreparedInlineMessage", Kind: kindWrite,
				Short: "Prepare a message for a Mini App user to send",
				Long:  "Store one InlineQueryResult for the user. Returns an identifier and expiration date.",
				Flags: []flagSpec{
					userFlag(),
					{Name: "result", Kind: flagJSON, Required: true, Usage: "single InlineQueryResult object as JSON"},
					{Name: "allow-user-chats", Kind: flagBool, Usage: "allow private chats with users"},
					{Name: "allow-bot-chats", Kind: flagBool, Usage: "allow private chats with bots"},
					{Name: "allow-group-chats", Kind: flagBool, Usage: "allow group and supergroup chats"},
					{Name: "allow-channel-chats", Kind: flagBool, Usage: "allow channel chats"},
				}, PreCall: requireResultObject,
			},
			{
				Use: "prepare-button", Method: "savePreparedKeyboardButton", Kind: kindWrite,
				Short: "Prepare a keyboard button for a Mini App user",
				Long:  "Store a KeyboardButton with exactly one of request_users, request_chat or request_managed_bot. Returns its identifier.",
				Flags: []flagSpec{
					userFlag(),
					{Name: "button", Kind: flagJSON, Required: true, Usage: "KeyboardButton object as JSON"},
				}, PreCall: requirePreparedButton,
			},
		},
	})
}

// requireResultObject rejects a --result that parsed as valid JSON but is not a single object.
// The neighbouring `inline answer` takes an ARRAY of results, so pasting one over here is the
// obvious slip, and JSON validity alone would let it through to fail at Telegram.
func requireResultObject(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	switch v := params["result"].(type) {
	case map[string]any:
		return nil
	case []any:
		return fmt.Errorf("--result takes ONE InlineQueryResult object, not an array (that is `inline answer`)")
	case nil:
		return fmt.Errorf("--result must be an InlineQueryResult object, got null")
	default:
		return fmt.Errorf("--result must be an InlineQueryResult object, got %T", v)
	}
}

func requirePreparedButton(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	button, ok := params["button"].(map[string]any)
	if !ok {
		return fmt.Errorf("--button must be a KeyboardButton object")
	}
	count := 0
	for _, field := range []string{"request_users", "request_chat", "request_managed_bot"} {
		if value, exists := button[field]; exists {
			if _, ok := value.(map[string]any); !ok {
				return fmt.Errorf("--button %s must be an object", field)
			}
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("--button must have exactly one of request_users, request_chat, or request_managed_bot")
	}
	for _, field := range []string{"request_contact", "request_location", "request_poll", "web_app"} {
		if _, exists := button[field]; exists {
			return fmt.Errorf("--button cannot combine %s with a request button", field)
		}
	}
	return nil
}
