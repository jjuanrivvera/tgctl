package commands

func init() {
	registerGroup(group{
		Use:   "callback",
		Short: "Answer callback queries from inline keyboards",
		Long:  "Respond to the CallbackQuery a user triggers by tapping an inline-keyboard button (answerCallbackQuery).",
		Cmds: []methodCmd{
			{
				Use: "answer", Method: "answerCallbackQuery", Kind: kindWrite,
				Short: "Answer a callback query (toast, alert, or URL)",
				Example: `  tgctl callback answer --callback-query-id 12345 --text "Saved!"
  tgctl callback answer --callback-query-id 12345 --text "Not allowed" --show-alert`,
				Flags: []flagSpec{
					{Name: "callback-query-id", Param: "callback_query_id", Required: true, Usage: "id of the callback query to answer"},
					{Name: "text", Usage: "notification text shown to the user (0-200 chars)"},
					{Name: "show-alert", Param: "show_alert", Kind: flagBool, Usage: "show an alert dialog instead of a toast"},
					{Name: "url", Usage: "URL the client opens (game / t.me deep link)"},
					{Name: "cache-time", Param: "cache_time", Kind: flagInt, Usage: "seconds the result may be cached client-side"},
				},
			},
		},
	})
}
