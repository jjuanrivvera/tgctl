package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
				Use: "set-tag", Method: "setChatMemberTag", Kind: kindWrite,
				Short: "Set (or clear) a regular member's tag",
				Long: `Set the tag shown next to a regular member's name in a group (setChatMemberTag).
The bot must be an administrator with the can_manage_tags right. Omit --tag to clear it.`,
				Example: `  tgctl member set-tag --chat @group --user 12345 --tag "Moderator"
  tgctl member set-tag --chat @group --user 12345   # clear the tag`,
				Flags: []flagSpec{
					chatFlag(), userFlag(),
					{Name: "tag", Usage: "new tag (0-16 chars, no emoji; omit to clear)"},
				},
			},
			{
				Use: "approve-join", Method: "approveChatJoinRequest", Kind: kindWrite,
				Short:   "Approve a chat join request",
				Example: `  tgctl member approve-join --chat @group --user 12345`,
				Flags:   []flagSpec{chatFlag(), userFlag()},
			},
			{
				Use: "answer-join-query", Method: "answerChatJoinRequestQuery", Kind: kindWrite,
				Short: "Answer a join-request query from a Mini App flow",
				Long: `Resolve a chat-join-request query (answerChatJoinRequestQuery).

This is the Mini App path, not the plain one: it answers a QUERY the bot received, identified
by its id, where ` + "`member approve-join` / `decline-join`" + ` act on a chat and a user
directly. --result takes approve, decline, or queue to leave the decision to another admin.`,
				Example: `  tgctl member answer-join-query --query-id AAxx... --result approve
  tgctl member answer-join-query --query-id AAxx... --result queue`,
				Flags: []flagSpec{
					{Name: "query-id", Param: "chat_join_request_query_id", Required: true, Usage: "id of the join-request query being answered"},
					{Name: "result", Required: true, Usage: "approve | decline | queue"},
				},
				PreCall: requireJoinQueryResult,
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

// requireJoinQueryResult rejects a --result the API does not define. The three words are the
// whole vocabulary of the parameter, and a typo would otherwise cost a round trip to be told
// so by Telegram, in a chat-join flow where the user is waiting on the other side.
func requireJoinQueryResult(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	result, _ := params["result"].(string)
	switch result {
	case "approve", "decline", "queue":
		return nil
	default:
		return fmt.Errorf("--result %q is not one of approve, decline or queue", result)
	}
}
