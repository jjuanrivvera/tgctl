package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExpansionVerbs_MockedAPI exercises every verb added in the 52→109 surface expansion against
// the mocked Bot API: each command must reach its method, render the result, and exit cleanly.
// It complements TestNewVerbs_MockedAPI (the earlier 31→52 batch) with the same harness.
func TestExpansionVerbs_MockedAPI(t *testing.T) {
	cases := []struct {
		name   string
		method string
		result string
		want   string
		args   []string
	}{
		// message: chat action, edit* family, stop-poll, and bulk ops
		{"message action", "sendChatAction", `true`, "true",
			[]string{"message", "action", "--chat", "@me", "--action", "typing"}},
		{"message edit-caption", "editMessageCaption", `{"message_id":5,"chat":{"id":7}}`, "5",
			[]string{"message", "edit-caption", "--chat", "@me", "--message-id", "5", "--caption", "new"}},
		{"message edit-media", "editMessageMedia", `{"message_id":5,"chat":{"id":7}}`, "5",
			[]string{"message", "edit-media", "--chat", "@me", "--message-id", "5", "--media", `{"type":"photo","media":"https://e.com/a.jpg"}`}},
		{"message edit-reply-markup", "editMessageReplyMarkup", `{"message_id":5,"chat":{"id":7}}`, "5",
			[]string{"message", "edit-reply-markup", "--chat", "@me", "--message-id", "5", "--reply-markup", `{"inline_keyboard":[]}`}},
		{"message edit-live-location", "editMessageLiveLocation", `{"message_id":5,"chat":{"id":7}}`, "5",
			[]string{"message", "edit-live-location", "--chat", "@me", "--message-id", "5", "--latitude", "3.4", "--longitude", "-76.5"}},
		{"message stop-live-location", "stopMessageLiveLocation", `{"message_id":5,"chat":{"id":7}}`, "5",
			[]string{"message", "stop-live-location", "--chat", "@me", "--message-id", "5"}},
		{"message stop-poll", "stopPoll", `{"id":"1","question":"Lunch?","total_voter_count":3,"is_closed":true}`, "Lunch?",
			[]string{"message", "stop-poll", "--chat", "@g", "--message-id", "5"}},
		{"message copy-batch", "copyMessages", `[{"message_id":10},{"message_id":11}]`, "10",
			[]string{"message", "copy-batch", "--chat", "@d", "--from-chat", "@s", "--message-ids", "[10,11]", "-o", "json"}},
		{"message forward-batch", "forwardMessages", `[{"message_id":10}]`, "10",
			[]string{"message", "forward-batch", "--chat", "@d", "--from-chat", "@s", "--message-ids", "[10]", "-o", "json"}},
		{"message delete-batch", "deleteMessages", `true`, "true",
			[]string{"message", "delete-batch", "--chat", "@g", "--message-ids", "[10,11]"}},

		// chat: admin setters, menu button, permissions, boosts
		{"chat delete-photo", "deleteChatPhoto", `true`, "true",
			[]string{"chat", "delete-photo", "--chat", "@g"}},
		{"chat set-permissions", "setChatPermissions", `true`, "true",
			[]string{"chat", "set-permissions", "--chat", "@g", "--permissions", `{"can_send_messages":true}`}},
		{"chat set-sticker-set", "setChatStickerSet", `true`, "true",
			[]string{"chat", "set-sticker-set", "--chat", "@g", "--sticker-set-name", "MyPack"}},
		{"chat delete-sticker-set", "deleteChatStickerSet", `true`, "true",
			[]string{"chat", "delete-sticker-set", "--chat", "@g"}},
		{"chat menu-button", "getChatMenuButton", `{"type":"commands"}`, "commands",
			[]string{"chat", "menu-button", "--chat", "123", "-o", "json"}},
		{"chat set-menu-button", "setChatMenuButton", `true`, "true",
			[]string{"chat", "set-menu-button", "--chat", "123", "--menu-button", `{"type":"default"}`}},
		{"chat unpin-all", "unpinAllChatMessages", `true`, "true",
			[]string{"chat", "unpin-all", "--chat", "@g"}},
		{"chat boosts", "getUserChatBoosts", `{"boosts":[{"boost_id":"b1"}]}`, "b1",
			[]string{"chat", "boosts", "--chat", "@g", "--user", "123", "-o", "json"}},

		// member: custom title, join requests, sender-chat bans
		{"member set-title", "setChatAdministratorCustomTitle", `true`, "true",
			[]string{"member", "set-title", "--chat", "@g", "--user", "5", "--title", "Lead"}},
		{"member approve-join", "approveChatJoinRequest", `true`, "true",
			[]string{"member", "approve-join", "--chat", "@g", "--user", "5"}},
		{"member decline-join", "declineChatJoinRequest", `true`, "true",
			[]string{"member", "decline-join", "--chat", "@g", "--user", "5"}},
		{"member ban-sender", "banChatSenderChat", `true`, "true",
			[]string{"member", "ban-sender", "--chat", "@g", "--sender-chat", "-100123"}},
		{"member unban-sender", "unbanChatSenderChat", `true`, "true",
			[]string{"member", "unban-sender", "--chat", "@g", "--sender-chat", "-100123"}},

		// forum topics + general topic
		{"forum create", "createForumTopic", `{"message_thread_id":9,"name":"News","icon_color":7322096}`, "9",
			[]string{"forum", "create", "--chat", "@g", "--name", "News"}},
		{"forum edit", "editForumTopic", `true`, "true",
			[]string{"forum", "edit", "--chat", "@g", "--thread", "9", "--name", "Renamed"}},
		{"forum close", "closeForumTopic", `true`, "true",
			[]string{"forum", "close", "--chat", "@g", "--thread", "9"}},
		{"forum reopen", "reopenForumTopic", `true`, "true",
			[]string{"forum", "reopen", "--chat", "@g", "--thread", "9"}},
		{"forum delete", "deleteForumTopic", `true`, "true",
			[]string{"forum", "delete", "--chat", "@g", "--thread", "9"}},
		{"forum unpin-all", "unpinAllForumTopicMessages", `true`, "true",
			[]string{"forum", "unpin-all", "--chat", "@g", "--thread", "9"}},
		{"forum icon-stickers", "getForumTopicIconStickers", `[{"emoji":"🔥","custom_emoji_id":"e1","set_name":"s"}]`, "e1",
			[]string{"forum", "icon-stickers", "-o", "json"}},
		{"forum edit-general", "editGeneralForumTopic", `true`, "true",
			[]string{"forum", "edit-general", "--chat", "@g", "--name", "General"}},
		{"forum close-general", "closeGeneralForumTopic", `true`, "true",
			[]string{"forum", "close-general", "--chat", "@g"}},
		{"forum reopen-general", "reopenGeneralForumTopic", `true`, "true",
			[]string{"forum", "reopen-general", "--chat", "@g"}},
		{"forum hide-general", "hideGeneralForumTopic", `true`, "true",
			[]string{"forum", "hide-general", "--chat", "@g"}},
		{"forum unhide-general", "unhideGeneralForumTopic", `true`, "true",
			[]string{"forum", "unhide-general", "--chat", "@g"}},
		{"forum unpin-all-general", "unpinAllGeneralForumTopicMessages", `true`, "true",
			[]string{"forum", "unpin-all-general", "--chat", "@g"}},

		// verify chats/users
		{"verify chat", "verifyChat", `true`, "true",
			[]string{"verify", "chat", "--chat", "@g", "--custom-description", "Official"}},
		{"verify user", "verifyUser", `true`, "true",
			[]string{"verify", "user", "--user", "5"}},
		{"verify remove-chat", "removeChatVerification", `true`, "true",
			[]string{"verify", "remove-chat", "--chat", "@g"}},
		{"verify remove-user", "removeUserVerification", `true`, "true",
			[]string{"verify", "remove-user", "--user", "5"}},

		// stars economy
		{"stars transactions", "getStarTransactions", `{"transactions":[{"id":"t1"}]}`, "t1",
			[]string{"stars", "transactions", "--limit", "5", "-o", "json"}},
		{"stars gifts", "getAvailableGifts", `{"gifts":[{"id":"g1"}]}`, "g1",
			[]string{"stars", "gifts", "-o", "json"}},
		{"stars send-gift", "sendGift", `true`, "true",
			[]string{"stars", "send-gift", "--user", "5", "--gift-id", "g1", "--text", "Enjoy"}},
		{"stars refund", "refundStarPayment", `true`, "true",
			[]string{"stars", "refund", "--user", "5", "--charge-id", "abc"}},
		{"stars edit-subscription", "editUserStarSubscription", `true`, "true",
			[]string{"stars", "edit-subscription", "--user", "5", "--charge-id", "abc", "--canceled"}},
		{"stars set-emoji-status", "setUserEmojiStatus", `true`, "true",
			[]string{"stars", "set-emoji-status", "--user", "5"}},
		{"stars send-paid-media", "sendPaidMedia", `{"message_id":3,"chat":{"id":7}}`, "3",
			[]string{"stars", "send-paid-media", "--chat", "@g", "--star-count", "10", "--media", `[{"type":"photo","media":"https://e.com/a.jpg"}]`}},

		// bot: short description, default admin rights, close/logout
		{"bot set-short-description", "setMyShortDescription", `true`, "true",
			[]string{"bot", "set-short-description", "--short-description", "hi"}},
		{"bot get-short-description", "getMyShortDescription", `{"short_description":"hi"}`, "short_description",
			[]string{"bot", "get-short-description", "-o", "json"}},
		{"bot set-admin-rights", "setMyDefaultAdministratorRights", `true`, "true",
			[]string{"bot", "set-admin-rights", "--rights", `{"can_manage_chat":true}`}},
		{"bot get-admin-rights", "getMyDefaultAdministratorRights", `{"can_manage_chat":true}`, "can_manage_chat",
			[]string{"bot", "get-admin-rights", "-o", "json"}},
		{"bot close", "close", `true`, "true",
			[]string{"bot", "close"}},
		{"bot logout", "logOut", `true`, "true",
			[]string{"bot", "logout"}},

		// invite: export + subscription links
		{"invite export", "exportChatInviteLink", `"https://t.me/+xyz"`, "xyz",
			[]string{"invite", "export", "--chat", "@g"}},
		{"invite create-subscription", "createChatSubscriptionInviteLink", `{"invite_link":"https://t.me/+sub","name":"VIP","subscription_period":2592000}`, "sub",
			[]string{"invite", "create-subscription", "--chat", "@g", "--subscription-period", "2592000", "--subscription-price", "100"}},
		{"invite edit-subscription", "editChatSubscriptionInviteLink", `{"invite_link":"https://t.me/+sub","name":"VIP"}`, "sub",
			[]string{"invite", "edit-subscription", "--chat", "@g", "--invite-link", "https://t.me/+sub", "--name", "VIP"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := newServer(t, routes{tc.method: tc.result})
			out, _, err := run(t, srv, tc.args...)
			require.NoError(t, err)
			assert.Contains(t, out, tc.want)
		})
	}
}

// TestExpansionSetChatPhoto_Upload covers the one new file-upload verb (setChatPhoto) through the
// real multipart path, mirroring TestMediaPhoto_Upload.
func TestExpansionSetChatPhoto_Upload(t *testing.T) {
	srv := newServer(t, routes{"setChatPhoto": `true`})
	dir := t.TempDir()
	pic := filepath.Join(dir, "logo.png")
	require.NoError(t, os.WriteFile(pic, []byte("PNG"), 0o600))
	out, _, err := run(t, srv, "chat", "set-photo", "--chat", "@g", "--photo", pic)
	require.NoError(t, err)
	assert.Contains(t, out, "true")
}
