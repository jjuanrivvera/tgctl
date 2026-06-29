package commands

func init() {
	registerGroup(group{
		Use:   "updates",
		Short: "Fetch incoming updates (long polling)",
		Long: `Read updates with getUpdates. Note: getUpdates conflicts with a set webhook —
delete the webhook first (tgctl webhook delete) if you want to poll.`,
		Cmds: []methodCmd{
			{
				Use: "get", Method: "getUpdates", Kind: kindRead,
				Short: "Get pending updates",
				Example: `  tgctl updates get --limit 5
  tgctl updates get --offset 123456789 --timeout 30 -o json
  tgctl updates get --allowed-updates message,callback_query`,
				Flags: []flagSpec{
					{Name: "offset", Kind: flagInt, Usage: "first update id to return (ack earlier ones)"},
					{Name: "limit", Kind: flagInt, Usage: "max updates to return (1-100)"},
					{Name: "timeout", Kind: flagInt, Usage: "long-poll seconds (0 = short poll)"},
					{Name: "allowed-updates", Param: "allowed_updates", Kind: flagStringSlice, Usage: "update types to receive"},
				},
				Columns: []string{"update_id", "message.message_id", "message.from.username", "message.text"},
			},
		},
	})
}
