package commands

func init() {
	registerGroup(group{
		Use:   "user",
		Short: "Read user information",
		Long:  "Inspect a user's public data, such as their profile photos (getUserProfilePhotos).",
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
		},
	})
}
