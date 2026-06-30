package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

// TestNewVerbs_MockedAPI exercises every verb added in the surface expansion against the mocked
// Bot API: each command must reach its method, render the result, and exit cleanly. File-send
// commands are driven with an http(s) URL / file_id so they go through the JSON path (the local
// multipart upload path is covered separately in commands_test.go and the api package).
func TestNewVerbs_MockedAPI(t *testing.T) {
	cases := []struct {
		name   string
		method string
		result string
		want   string
		args   []string
	}{
		// media sends
		{"media audio", "sendAudio", `{"message_id":1,"chat":{"id":7}}`, "1",
			[]string{"media", "audio", "--chat", "@me", "--audio", "https://e.com/a.mp3", "--performer", "X"}},
		{"media voice", "sendVoice", `{"message_id":1,"chat":{"id":7}}`, "1",
			[]string{"media", "voice", "--chat", "@me", "--voice", "https://e.com/v.ogg", "--duration", "5"}},
		{"media animation", "sendAnimation", `{"message_id":1,"chat":{"id":7}}`, "1",
			[]string{"media", "animation", "--chat", "@me", "--animation", "https://e.com/a.gif"}},
		{"media video-note", "sendVideoNote", `{"message_id":1,"chat":{"id":7}}`, "1",
			[]string{"media", "video-note", "--chat", "@me", "--video-note", "EXISTINGFILEID", "--length", "240"}},
		{"media sticker", "sendSticker", `{"message_id":1,"chat":{"id":7}}`, "1",
			[]string{"media", "sticker", "--chat", "@me", "--sticker", "CAACfileid", "--emoji", "🔥"}},
		{"media media-group", "sendMediaGroup", `[{"message_id":1,"chat":{"id":7}}]`, "1",
			[]string{"media", "media-group", "--chat", "@me", "--media", `[{"type":"photo","media":"https://e.com/a.jpg"}]`}},

		// message rich-content sends + reactions
		{"message react", "setMessageReaction", `true`, "true",
			[]string{"message", "react", "--chat", "@g", "--message-id", "5", "--reaction", `[{"type":"emoji","emoji":"👍"}]`, "--is-big"}},
		{"message location", "sendLocation", `{"message_id":2,"chat":{"id":7}}`, "2",
			[]string{"message", "location", "--chat", "@me", "--latitude", "3.45", "--longitude", "-76.5"}},
		{"message venue", "sendVenue", `{"message_id":2,"chat":{"id":7}}`, "2",
			[]string{"message", "venue", "--chat", "@me", "--latitude", "3.45", "--longitude", "-76.5", "--title", "Office", "--address", "Av 1"}},
		{"message contact", "sendContact", `{"message_id":2,"chat":{"id":7}}`, "2",
			[]string{"message", "contact", "--chat", "@me", "--phone-number", "+15551234567", "--first-name", "Ada"}},
		{"message poll", "sendPoll", `{"message_id":3,"chat":{"id":7}}`, "3",
			[]string{"message", "poll", "--chat", "@g", "--question", "Lunch?", "--options", `[{"text":"A"},{"text":"B"}]`, "--type", "regular"}},
		{"message dice", "sendDice", `{"message_id":4,"chat":{"id":7},"dice":{"emoji":"🎯","value":5}}`, "4",
			[]string{"message", "dice", "--chat", "@me", "--emoji", "🎯"}},

		// chat admin
		{"chat set-title", "setChatTitle", `true`, "true",
			[]string{"chat", "set-title", "--chat", "@g", "--title", "New name"}},
		{"chat set-description", "setChatDescription", `true`, "true",
			[]string{"chat", "set-description", "--chat", "@g", "--description", "About"}},

		// invite links
		{"invite create", "createChatInviteLink", `{"invite_link":"https://t.me/+abc","name":"Launch","is_primary":false}`, "abc",
			[]string{"invite", "create", "--chat", "@g", "--name", "Launch", "--member-limit", "100"}},
		{"invite edit", "editChatInviteLink", `{"invite_link":"https://t.me/+abc","name":"Launch"}`, "abc",
			[]string{"invite", "edit", "--chat", "@g", "--invite-link", "https://t.me/+abc", "--member-limit", "10"}},
		{"invite revoke", "revokeChatInviteLink", `{"invite_link":"https://t.me/+abc","is_revoked":true}`, "abc",
			[]string{"invite", "revoke", "--chat", "@g", "--invite-link", "https://t.me/+abc"}},

		// user + file metadata
		{"user photos", "getUserProfilePhotos", `{"total_count":2,"photos":[]}`, "2",
			[]string{"user", "photos", "--user", "123", "--limit", "1"}},
		{"file info", "getFile", `{"file_id":"ABC","file_path":"photos/x.jpg","file_size":100}`, "x.jpg",
			[]string{"file", "info", "--file-id", "ABC"}},

		// callbacks + inline
		{"callback answer", "answerCallbackQuery", `true`, "true",
			[]string{"callback", "answer", "--callback-query-id", "999", "--text", "ok", "--show-alert"}},
		{"inline answer", "answerInlineQuery", `true`, "true",
			[]string{"inline", "answer", "--inline-query-id", "999", "--results", "[]", "--is-personal"}},
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

// fileServer answers both the getFile method call and a /file/ download with the given bytes.
func fileServer(t *testing.T, filePath, content string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/getFile"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"file_id":"ABC","file_path":"` + filePath + `","file_size":` + itoaTest(len(content)) + `}}`))
		case strings.Contains(r.URL.Path, "/file/"):
			_, _ = w.Write([]byte(content))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func itoaTest(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func TestFileDownload_ToFile(t *testing.T) {
	srv := fileServer(t, "photos/x.jpg", "HELLO")
	dir := t.TempDir()
	dest := filepath.Join(dir, "got.jpg")
	_, errb, err := run(t, srv, "file", "download", "--file-id", "ABC", "--dest", dest)
	require.NoError(t, err)
	mustContain(t, errb, "downloaded 5 bytes")

	data, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "HELLO", string(data))
}

func TestFileDownload_Stdout(t *testing.T) {
	srv := fileServer(t, "photos/x.jpg", "BYTES")
	out, _, err := run(t, srv, "file", "download", "--file-id", "ABC", "--dest", "-")
	require.NoError(t, err)
	mustContain(t, out, "BYTES")
}

func TestFileDownload_DryRun(t *testing.T) {
	_, errb, err := run(t, nil, "file", "download", "--file-id", "ABC", "--dry-run")
	require.NoError(t, err)
	mustContain(t, errb, "getFile")
	mustContain(t, errb, "would then GET")
}

func TestFileDownload_NoPath(t *testing.T) {
	srv := newServer(t, routes{"getFile": `{"file_id":"ABC"}`}) // no file_path → too big to fetch
	_, _, err := run(t, srv, "file", "download", "--file-id", "ABC")
	require.Error(t, err)
	mustContain(t, err.Error(), "no downloadable path")
}

func TestBotFlag_SelectsProfile(t *testing.T) {
	keyring.MockInit()
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	dir := t.TempDir()
	// --bot is the new name for profile selection; logging in under it creates that profile.
	_, _, err := runIn(t, dir, srv, "", "auth", "login", "--token", "111:AAA", "--bot", "prod")
	require.NoError(t, err)
	out, _, err := runIn(t, dir, srv, "", "config", "list-profiles", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "prod")
}
