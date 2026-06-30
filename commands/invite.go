package commands

func init() {
	registerGroup(group{
		Use:   "invite",
		Short: "Manage chat invite links",
		Long:  "Create, edit, and revoke additional invite links for a chat (the bot must be an admin with can_invite_users).",
		Cmds: []methodCmd{
			{
				Use: "create", Method: "createChatInviteLink", Kind: kindWrite,
				Short: "Create a new invite link",
				Example: `  tgctl invite create --chat @group --name "Launch" --member-limit 100
  tgctl invite create --chat @group --creates-join-request`,
				Flags: []flagSpec{
					chatFlag(),
					{Name: "name", Usage: "invite link name (0-32 chars)"},
					{Name: "expire-date", Param: "expire_date", Kind: flagInt, Usage: "unix time the link expires"},
					{Name: "member-limit", Param: "member_limit", Kind: flagInt, Usage: "max users that may join via this link (1-99999)"},
					{Name: "creates-join-request", Param: "creates_join_request", Kind: flagBool, Usage: "users joining are placed in a join-request queue"},
				},
				Columns: []string{"invite_link", "name", "is_primary", "member_limit"},
			},
			{
				Use: "edit", Method: "editChatInviteLink", Kind: kindWrite,
				Short:   "Edit an existing invite link",
				Example: `  tgctl invite edit --chat @group --invite-link https://t.me/+abc --member-limit 10`,
				Flags: []flagSpec{
					chatFlag(),
					inviteLinkFlag(),
					{Name: "name", Usage: "invite link name (0-32 chars)"},
					{Name: "expire-date", Param: "expire_date", Kind: flagInt, Usage: "unix time the link expires"},
					{Name: "member-limit", Param: "member_limit", Kind: flagInt, Usage: "max users that may join via this link (1-99999)"},
					{Name: "creates-join-request", Param: "creates_join_request", Kind: flagBool, Usage: "users joining are placed in a join-request queue"},
				},
				Columns: []string{"invite_link", "name", "member_limit"},
			},
			{
				Use: "revoke", Method: "revokeChatInviteLink", Kind: kindDestructive,
				Short:   "Revoke an invite link (a new one is generated automatically)",
				Example: `  tgctl invite revoke --chat @group --invite-link https://t.me/+abc`,
				Flags:   []flagSpec{chatFlag(), inviteLinkFlag()},
				Columns: []string{"invite_link", "is_revoked"},
			},
		},
	})
}

func inviteLinkFlag() flagSpec {
	return flagSpec{Name: "invite-link", Param: "invite_link", Required: true, Usage: "the invite link to act on"}
}
