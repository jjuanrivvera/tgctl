package commands

func init() {
	registerGroup(group{
		Use:   "forum",
		Short: "Manage forum topics in supergroups",
		Long:  "Create, edit, close, reopen, and delete forum topics, plus manage the General topic. The bot must be an admin with can_manage_topics.",
		Cmds: []methodCmd{
			{
				Use: "create", Method: "createForumTopic", Kind: kindWrite,
				Short:   "Create a forum topic",
				Example: `  tgctl forum create --chat @group --name "Announcements" --icon-color 7322096`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "name", Required: true, Usage: "topic name (1-128 chars)"},
					{Name: "icon-color", Param: "icon_color", Kind: flagInt, Usage: "RGB color of the topic icon (one of the allowed palette values)"},
					{Name: "icon-custom-emoji-id", Param: "icon_custom_emoji_id", Usage: "custom emoji id shown as the topic icon"},
				},
				Columns: []string{"message_thread_id", "name", "icon_color"},
			},
			{
				Use: "edit", Method: "editForumTopic", Kind: kindWrite,
				Short:   "Edit a forum topic's name or icon",
				Example: `  tgctl forum edit --chat @group --thread 42 --name "Renamed"`,
				Flags: []flagSpec{
					chatFlag(), threadReqFlag(),
					{Name: "name", Usage: "new topic name (1-128 chars; omit to keep)"},
					{Name: "icon-custom-emoji-id", Param: "icon_custom_emoji_id", Usage: "new custom emoji id (empty string removes the icon)"},
				},
			},
			{
				Use: "close", Method: "closeForumTopic", Kind: kindWrite,
				Short:   "Close a forum topic",
				Example: `  tgctl forum close --chat @group --thread 42`,
				Flags:   []flagSpec{chatFlag(), threadReqFlag()},
			},
			{
				Use: "reopen", Method: "reopenForumTopic", Kind: kindWrite,
				Short:   "Reopen a closed forum topic",
				Example: `  tgctl forum reopen --chat @group --thread 42`,
				Flags:   []flagSpec{chatFlag(), threadReqFlag()},
			},
			{
				Use: "delete", Method: "deleteForumTopic", Kind: kindDestructive,
				Short:   "Delete a forum topic and all its messages",
				Example: `  tgctl forum delete --chat @group --thread 42`,
				Flags:   []flagSpec{chatFlag(), threadReqFlag()},
			},
			{
				// Destructive: bulk-unpin whose previous pin set cannot be recovered.
				Use: "unpin-all", Method: "unpinAllForumTopicMessages", Kind: kindDestructive,
				Short:   "Unpin all messages in a forum topic",
				Example: `  tgctl forum unpin-all --chat @group --thread 42`,
				Flags:   []flagSpec{chatFlag(), threadReqFlag()},
			},
			{
				Use: "icon-stickers", Method: "getForumTopicIconStickers", Kind: kindRead,
				Short:   "List the custom emoji stickers usable as topic icons",
				Example: `  tgctl forum icon-stickers -o json`,
				Columns: []string{"emoji", "custom_emoji_id", "set_name"},
			},
			{
				Use: "edit-general", Method: "editGeneralForumTopic", Kind: kindWrite,
				Short:   "Rename the General forum topic",
				Example: `  tgctl forum edit-general --chat @group --name "General chat"`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "name", Required: true, Usage: "new name for the General topic (1-128 chars)"},
				},
			},
			{
				Use: "close-general", Method: "closeGeneralForumTopic", Kind: kindWrite,
				Short:   "Close the General forum topic",
				Example: `  tgctl forum close-general --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "reopen-general", Method: "reopenGeneralForumTopic", Kind: kindWrite,
				Short:   "Reopen the General forum topic (also unhides it)",
				Example: `  tgctl forum reopen-general --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "hide-general", Method: "hideGeneralForumTopic", Kind: kindWrite,
				Short:   "Hide the General forum topic (also closes it)",
				Example: `  tgctl forum hide-general --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				Use: "unhide-general", Method: "unhideGeneralForumTopic", Kind: kindWrite,
				Short:   "Unhide the General forum topic",
				Example: `  tgctl forum unhide-general --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
			{
				// Destructive: bulk-unpin whose previous pin set cannot be recovered.
				Use: "unpin-all-general", Method: "unpinAllGeneralForumTopicMessages", Kind: kindDestructive,
				Short:   "Unpin all messages in the General forum topic",
				Example: `  tgctl forum unpin-all-general --chat @group`,
				Flags:   []flagSpec{chatFlag()},
			},
		},
	})
}

// threadReqFlag is the required forum-topic thread id (message_thread_id) for topic operations.
func threadReqFlag() flagSpec {
	return flagSpec{Name: "thread", Param: "message_thread_id", Kind: flagInt, Required: true, Usage: "forum topic thread id"}
}
