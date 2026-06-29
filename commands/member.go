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
		},
	})
}
