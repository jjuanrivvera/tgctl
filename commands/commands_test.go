package commands

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBotInfo(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":123456,"is_bot":true,"first_name":"Test","username":"testbot","can_join_groups":true}`})
	out, _, err := run(t, srv, "bot", "info")
	require.NoError(t, err)
	mustContain(t, out, "USERNAME")
	mustContain(t, out, "testbot")
}

func TestBotInfo_JSON(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":123456,"is_bot":true,"first_name":"T","username":"b"}`})
	out, _, err := run(t, srv, "bot", "info", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, `"username": "b"`)
}

func TestMessageSend(t *testing.T) {
	srv := newServer(t, routes{"sendMessage": `{"message_id":42,"chat":{"id":7,"type":"private"},"date":1700000000,"text":"hi"}`})
	out, _, err := run(t, srv, "message", "send", "--chat", "@me", "--text", "hi")
	require.NoError(t, err)
	mustContain(t, out, "MESSAGE_ID")
	mustContain(t, out, "42")
}

func TestMessageDelete_ScalarResult(t *testing.T) {
	srv := newServer(t, routes{"deleteMessage": `true`})
	out, _, err := run(t, srv, "message", "delete", "--chat", "@me", "--message-id", "42")
	require.NoError(t, err)
	assert.Equal(t, "true\n", out)
}

func TestChatMembersCount(t *testing.T) {
	srv := newServer(t, routes{"getChatMemberCount": `1234`})
	out, _, err := run(t, srv, "chat", "members-count", "--chat", "@group")
	require.NoError(t, err)
	assert.Equal(t, "1234\n", out)
}

func TestUpdatesGet_Array(t *testing.T) {
	srv := newServer(t, routes{"getUpdates": `[{"update_id":100,"message":{"message_id":1,"text":"a","from":{"id":1,"username":"u"}}}]`})
	out, _, err := run(t, srv, "updates", "get", "--limit", "5")
	require.NoError(t, err)
	mustContain(t, out, "UPDATE_ID")
	mustContain(t, out, "100")
}

func TestCommandsList_CSV(t *testing.T) {
	srv := newServer(t, routes{"getMyCommands": `[{"command":"start","description":"Begin"},{"command":"help","description":"Help"}]`})
	out, _, err := run(t, srv, "commands", "list", "-o", "csv")
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	assert.Equal(t, "command,description", lines[0])
	assert.Equal(t, "start,Begin", lines[1])
}

func TestAPIError_SurfacesHint(t *testing.T) {
	srv := newServer(t, routes{}) // every method 404s
	_, _, err := run(t, srv, "bot", "info")
	require.Error(t, err)
	mustContain(t, err.Error(), "method name")
}

func TestRawAPI(t *testing.T) {
	srv := newServer(t, routes{"getChat": `{"id":7,"type":"group","title":"Room"}`})
	out, _, err := run(t, srv, "api", "getChat", "-q", "chat_id=@room", "--idempotent")
	require.NoError(t, err)
	mustContain(t, out, "Room")
}

func TestDryRun_NoRequest(t *testing.T) {
	// No server is contacted; the curl line goes to stderr and stdout stays empty.
	_, errb, err := run(t, nil, "message", "send", "--chat", "@me", "--text", "hi", "--dry-run")
	require.NoError(t, err)
	mustContain(t, errb, "curl -sS -X POST")
	mustContain(t, errb, "sendMessage")
	mustContain(t, errb, "123456:<redacted>")
}

func TestMediaPhoto_Upload(t *testing.T) {
	srv := newServer(t, routes{"sendPhoto": `{"message_id":9,"chat":{"id":7}}`})
	dir := t.TempDir()
	pic := dir + "/p.jpg"
	require.NoError(t, os.WriteFile(pic, []byte("JPEG"), 0o600))
	out, _, err := run(t, srv, "media", "photo", "--chat", "@me", "--photo", pic)
	require.NoError(t, err)
	mustContain(t, out, "9")
}

func TestMediaPhoto_URLAsParam(t *testing.T) {
	srv := newServer(t, routes{"sendPhoto": `{"message_id":9}`})
	out, _, err := run(t, srv, "media", "photo", "--chat", "@me", "--photo", "https://example.com/p.png")
	require.NoError(t, err)
	mustContain(t, out, "9")
}

func TestJQFilter(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":1,"username":"bot","first_name":"B"}`})
	out, _, err := run(t, srv, "bot", "info", "-o", "json", "--jq", ".username")
	require.NoError(t, err)
	mustContain(t, out, `"bot"`)
}
