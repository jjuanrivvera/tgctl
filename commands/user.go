package commands

func init() {
	registerGroup(group{
		Use:   "user",
		Short: "Read user information",
		Long:  "Inspect a user's public data: profile photos and audios, and the last messages on their personal chat.",
		Cmds: []methodCmd{
			{
				Use: "photos", Method: "getUserProfilePhotos", Kind: kindRead,
				Short: "List a user's profile photos",
				Example: `  tgctl user photos --user 12345
  tgctl user photos --user 12345 --limit 1 -o json`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "offset", Kind: flagInt, Usage: "number of photos to skip"},
					{Name: "limit", Kind: flagInt, Usage: "max photos to return (1-100)"},
				},
				Columns: []string{"total_count"},
			},
			{
				Use: "audios", Method: "getUserProfileAudios", Kind: kindRead,
				Short: "List a user's profile audios",
				Example: `  tgctl user audios --user 12345
  tgctl user audios --user 12345 --limit 1 -o json`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "offset", Kind: flagInt, Usage: "number of audios to skip"},
					{Name: "limit", Kind: flagInt, Usage: "max audios to return (1-100)"},
				},
				Columns: []string{"total_count"},
			},
			{
				Use: "chat-messages", Method: "getUserPersonalChatMessages", Kind: kindRead,
				Short: "Read the last messages on a user's personal chat",
				Long: `Read the most recent messages from the chat a user has pinned to their profile
(getUserPersonalChatMessages). --limit is required by the API and capped at 20.`,
				Example: `  tgctl user chat-messages --user 12345 --limit 5
  tgctl user chat-messages --user 12345 --limit 20 -o json`,
				Flags: []flagSpec{
					userFlag(),
					{Name: "limit", Kind: flagInt, Required: true, Usage: "how many messages to return (1-20)"},
				},
				Columns: []string{"message_id", "date", "text"},
			},
		},
	})
}
