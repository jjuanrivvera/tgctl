package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

// TestAllCommands_DryRun exercises every generated command's flag binding, param assembly,
// and file handling via --dry-run (which performs no request but runs the whole build path).
func TestAllCommands_DryRun(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		method string
	}{
		{"bot set-name", []string{"bot", "set-name", "--name", "Bot"}, "setMyName"},
		{"bot get-name", []string{"bot", "get-name"}, "getMyName"},
		{"bot set-description", []string{"bot", "set-description", "--description", "hi"}, "setMyDescription"},
		{"bot get-description", []string{"bot", "get-description"}, "getMyDescription"},
		{"message edit", []string{"message", "edit", "--chat", "@me", "--message-id", "1", "--text", "x"}, "editMessageText"},
		{"message forward", []string{"message", "forward", "--chat", "@a", "--from-chat", "@b", "--message-id", "1"}, "forwardMessage"},
		{"message copy", []string{"message", "copy", "--chat", "@a", "--from-chat", "@b", "--message-id", "1"}, "copyMessage"},
		{"message pin", []string{"message", "pin", "--chat", "@a", "--message-id", "1"}, "pinChatMessage"},
		{"message unpin", []string{"message", "unpin", "--chat", "@a"}, "unpinChatMessage"},
		{"chat administrators", []string{"chat", "administrators", "--chat", "@a"}, "getChatAdministrators"},
		{"chat member", []string{"chat", "member", "--chat", "@a", "--user", "5"}, "getChatMember"},
		{"chat leave", []string{"chat", "leave", "--chat", "@a"}, "leaveChat"},
		{"member ban", []string{"member", "ban", "--chat", "@a", "--user", "5", "--revoke-messages"}, "banChatMember"},
		{"member unban", []string{"member", "unban", "--chat", "@a", "--user", "5"}, "unbanChatMember"},
		{"member restrict", []string{"member", "restrict", "--chat", "@a", "--user", "5", "--permissions", `{"can_send_messages":false}`}, "restrictChatMember"},
		{"member promote", []string{"member", "promote", "--chat", "@a", "--user", "5", "--can-pin-messages"}, "promoteChatMember"},
		{"webhook set", []string{"webhook", "set", "--url", "https://e.com/b", "--max-connections", "10"}, "setWebhook"},
		{"webhook delete", []string{"webhook", "delete", "--drop-pending"}, "deleteWebhook"},
		{"commands delete", []string{"commands", "delete"}, "deleteMyCommands"},
		{"media document URL", []string{"media", "document", "--chat", "@me", "--document", "https://e.com/f.pdf"}, "sendDocument"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, errb, err := run(t, nil, append(tc.args, "--dry-run")...)
			require.NoError(t, err)
			assert.Contains(t, errb, tc.method)
		})
	}
}

func TestWebhookInfo(t *testing.T) {
	srv := newServer(t, routes{"getWebhookInfo": `{"url":"https://e.com/b","pending_update_count":3}`})
	out, _, err := run(t, srv, "webhook", "info")
	require.NoError(t, err)
	mustContain(t, out, "URL")
}

func TestChatAdministrators(t *testing.T) {
	srv := newServer(t, routes{"getChatAdministrators": `[{"status":"creator","user":{"id":1,"username":"boss"}}]`})
	out, _, err := run(t, srv, "chat", "administrators", "--chat", "@g")
	require.NoError(t, err)
	mustContain(t, out, "creator")
}

func TestAuthLogout(t *testing.T) {
	keyring.MockInit()
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	dir := t.TempDir()
	_, _, err := runIn(t, dir, srv, "", "auth", "login", "--token", "111:AAA")
	require.NoError(t, err)
	out, _, err := runIn(t, dir, srv, "", "auth", "logout")
	require.NoError(t, err)
	mustContain(t, out, "logged out")
}

func TestInitWizard(t *testing.T) {
	keyring.MockInit()
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"wizbot"}`})
	// stdin supplies the base URL (blank → default is overridden by --base-url anyway) then the token.
	stdin := "\n111:WIZARDTOKEN\n"
	out, errb, err := runNoToken(t, srv, stdin, "init")
	require.NoError(t, err)
	_ = errb
	mustContain(t, out, "wizbot")
}

func TestAliasExpansionEndToEnd(t *testing.T) {
	keyring.MockInit()
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	dir := t.TempDir()
	_, _, err := runIn(t, dir, srv, "", "alias", "set", "ping", "bot info")
	require.NoError(t, err)
	// ExpandAliases must turn ["ping"] into ["bot","info"] using the saved config.
	t.Setenv("XDG_CONFIG_HOME", dir)
	got := ExpandAliases([]string{"ping"})
	assert.Equal(t, []string{"bot", "info"}, got)

	out, _, err := runIn(t, dir, srv, "", "alias", "list")
	require.NoError(t, err)
	mustContain(t, out, "ping = bot info")

	_, _, err = runIn(t, dir, srv, "", "alias", "remove", "ping")
	require.NoError(t, err)
}
