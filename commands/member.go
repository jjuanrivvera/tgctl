package commands

func init() {
	registerGroup(group{
		Use:   "member",
		Short: "Moderate chat members (ban, restrict, promote)",
		Long:  "Administrative actions on members. The bot must be an admin with the relevant rights.",
		Cmds: []methodCmd{
			{
				Use: "ban", Method: "banChatMember", Kind: kindDestructive,
				Short: "Ban a user from a chat",
				Example: `  tgctl member ban --chat @group --user 12345
  tgctl member ban --chat @group --user 12345 --revoke-messages`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "until", Param: "until_date", Kind: flagInt, Usage: "unix time to auto-unban (0 = forever)"},
					{Name: "revoke-messages", Param: "revoke_messages", Kind: flagBool, Usage: "delete all the user's messages"},
				},
			},
			{
				Use: "unban", Method: "unbanChatMember", Kind: kindWrite,
				Short:   "Unban a previously banned user",
				Example: `  tgctl member unban --chat @group --user 12345 --only-if-banned`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "only-if-banned", Param: "only_if_banned", Kind: flagBool, Usage: "do nothing if the user is not banned"},
				},
			},
			{
				Use: "restrict", Method: "restrictChatMember", Kind: kindWrite,
				Short: "Restrict what a member can do",
				Example: `  tgctl member restrict --chat @group --user 12345 \
    --permissions '{"can_send_messages":false}'`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "permissions", Kind: flagJSON, Required: true, Usage: "ChatPermissions object as JSON"},
					{Name: "until", Param: "until_date", Kind: flagInt, Usage: "unix time the restriction lifts"},
				},
			},
			{
				Use: "promote", Method: "promoteChatMember", Kind: kindWrite,
				Short:   "Promote or demote an administrator",
				Example: `  tgctl member promote --chat @group --user 12345 --can-delete-messages --can-pin-messages`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "can-manage-chat", Param: "can_manage_chat", Kind: flagBool, Usage: "can access the admin log, etc."},
					{Name: "can-delete-messages", Param: "can_delete_messages", Kind: flagBool, Usage: "can delete others' messages"},
					{Name: "can-restrict-members", Param: "can_restrict_members", Kind: flagBool, Usage: "can restrict/ban members"},
					{Name: "can-promote-members", Param: "can_promote_members", Kind: flagBool, Usage: "can add new admins"},
					{Name: "can-change-info", Param: "can_change_info", Kind: flagBool, Usage: "can change chat title/photo"},
					{Name: "can-invite-users", Param: "can_invite_users", Kind: flagBool, Usage: "can invite new users"},
					{Name: "can-pin-messages", Param: "can_pin_messages", Kind: flagBool, Usage: "can pin messages"},
				},
			},
			{
				Use: "set-title", Method: "setChatAdministratorCustomTitle", Kind: kindWrite,
				Short:   "Set a custom title for an administrator the bot promoted",
				Example: `  tgctl member set-title --chat @group --user 12345 --title "Community Lead"`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "title", Param: "custom_title", Required: true, Usage: "custom admin title (0-16 chars, no emoji)"},
				},
			},
			{
				Use: "approve-join", Method: "approveChatJoinRequest", Kind: kindWrite,
				Short:   "Approve a chat join request",
				Example: `  tgctl member approve-join --chat @group --user 12345`,
				Flags:   []flagSpec{chatFlag(), userFlag()},
			},
			{
				Use: "decline-join", Method: "declineChatJoinRequest", Kind: kindWrite,
				Short:   "Decline a chat join request",
				Example: `  tgctl member decline-join --chat @group --user 12345`,
				Flags:   []flagSpec{chatFlag(), userFlag()},
			},
			{
				Use: "ban-sender", Method: "banChatSenderChat", Kind: kindDestructive,
				Short:   "Ban a channel from posting as itself in a chat",
				Example: `  tgctl member ban-sender --chat @group --sender-chat -1001234567890`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "sender-chat", Param: "sender_chat_id", Kind: flagInt, Required: true, Usage: "id of the channel/chat to ban as a sender"},
				},
			},
			{
				Use: "unban-sender", Method: "unbanChatSenderChat", Kind: kindWrite,
				Short:   "Unban a channel that was banned as a sender",
				Example: `  tgctl member unban-sender --chat @group --sender-chat -1001234567890`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "sender-chat", Param: "sender_chat_id", Kind: flagInt, Required: true, Usage: "id of the channel/chat to unban"},
				},
			},
		},
	})
}
