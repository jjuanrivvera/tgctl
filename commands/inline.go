package commands

func init() {
	registerGroup(group{
		Use:   "inline",
		Short: "Answer inline queries",
		Long:  "Respond to an inline query (a user typing @yourbot ...) with a list of results (answerInlineQuery).",
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
		},
	})
}
