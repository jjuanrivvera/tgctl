package commands

func init() {
	registerGroup(group{
		Use:   "media",
		Short: "Send files: photos, documents, and video",
		Long: `Send media to a chat. Each --<kind> flag accepts a local file path (uploaded as
multipart/form-data), an http(s) URL (Telegram fetches it), or an existing file_id.`,
		Cmds: []methodCmd{
			{
				Use: "photo", Method: "sendPhoto", Kind: kindWrite,
				Short: "Send a photo",
				Example: `  tgctl media photo --chat @me --photo ./cat.jpg --caption "my cat"
  tgctl media photo --chat @me --photo https://example.com/pic.png`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "photo caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "photo", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "document", Method: "sendDocument", Kind: kindWrite,
				Short:   "Send a document/file",
				Example: `  tgctl media document --chat @me --document ./report.pdf --caption "Q2 report"`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "document caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "document", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
			{
				Use: "video", Method: "sendVideo", Kind: kindWrite,
				Short:   "Send a video",
				Example: `  tgctl media video --chat @me --video ./clip.mp4`,
				Flags:   []flagSpec{chatFlag(), {Name: "caption", Usage: "video caption"}, parseModeFlag(), silentFlag()},
				Files:   []fileSpec{{Name: "video", Required: true, Usage: "local path, URL, or file_id"}},
				Columns: []string{"message_id", "chat.id"},
			},
		},
	})
}

func silentFlag() flagSpec {
	return flagSpec{Name: "silent", Param: "disable_notification", Kind: flagBool, Usage: "send without a notification sound"}
}
