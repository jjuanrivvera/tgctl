package commands

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
			},
		},
	})
}
