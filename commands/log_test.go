package commands

import (
	"testing"

	"github.com/njayp/ophis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog_Empty(t *testing.T) {
	out, _, err := run(t, nil, "log", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", out)
}

func TestLog_RecordsOutboundSend(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":42,"chat":{"id":7},"text":"hello from tgctl"}`})

	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "hello from tgctl")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "hello from tgctl")
	mustContain(t, out, `"direction": "out"`)
	mustContain(t, out, `"chat_id": 7`)
	mustContain(t, out, `"message_id": 42`)
}

func TestLog_NoStore_SkipsRecording(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7},"text":"ephemeral"}`})

	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "ephemeral", "--no-store")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", out, "--no-store must skip recording the send")
}

func TestLog_FilterByChatAndKind(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{
		"sendMessage": `{"message_id":1,"chat":{"id":7},"text":"to seven"}`,
	})
	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "to seven")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "--chat", "7", "--kind", "text", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "to seven")

	out, _, err = runIn(t, dir, srv, "tok:val", "log", "--chat", "999", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", out)

	out, _, err = runIn(t, dir, srv, "tok:val", "log", "--kind", "photo", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", out)
}

func TestLog_InvalidSince(t *testing.T) {
	_, _, err := run(t, nil, "log", "--since", "not-a-time")
	require.Error(t, err)
	mustContain(t, err.Error(), "invalid --since")
}

func TestLog_SinceDuration(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7},"text":"recent"}`})
	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "recent")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "--since", "1h", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "recent")

	out, _, err = runIn(t, dir, srv, "tok:val", "log", "--since", "-1h", "-o", "json")
	require.NoError(t, err)
	_ = out // a negative duration is still a valid Go duration (since = now + 1h, i.e. the future)
}

func TestLog_Search(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{
		"sendMessage": `{"message_id":1,"chat":{"id":7},"text":"deploy failed on staging"}`,
	})
	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "deploy failed on staging")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "search", "deploy", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "staging")

	out, _, err = runIn(t, dir, srv, "tok:val", "log", "search", "nomatch", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", out)
}

func TestLog_Show(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":55,"chat":{"id":7},"text":"show me"}`})
	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "show me")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "show", "55", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "show me")
	mustContain(t, out, `"raw"`)
}

func TestLog_Show_NotFound(t *testing.T) {
	_, _, err := run(t, nil, "log", "show", "999")
	require.Error(t, err)
	mustContain(t, err.Error(), "no recorded message")
}

func TestLog_Show_InvalidID(t *testing.T) {
	_, _, err := run(t, nil, "log", "show", "not-a-number")
	require.Error(t, err)
	mustContain(t, err.Error(), "invalid message id")
}

func TestLog_Prune(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{"sendMessage": `{"message_id":1,"chat":{"id":7},"text":"old news"}`})
	_, _, err := runIn(t, dir, srv, "tok:val", "message", "send", "--chat", "7", "--text", "old news")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "prune", "--older-than", "0s")
	require.NoError(t, err)
	mustContain(t, out, "pruned 1 message")

	logOut, _, err := runIn(t, dir, srv, "tok:val", "log", "-o", "json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", logOut)
}

func TestLog_Prune_InvalidDuration(t *testing.T) {
	_, _, err := run(t, nil, "log", "prune", "--older-than", "not-a-duration")
	require.Error(t, err)
	mustContain(t, err.Error(), "invalid --older-than")
}

func TestLog_Prune_MissingRequiredFlag(t *testing.T) {
	_, _, err := run(t, nil, "log", "prune")
	require.Error(t, err)
}

func TestUpdatesGet_RecordsInbound(t *testing.T) {
	dir := t.TempDir()
	srv := newServer(t, routes{
		"getUpdates": `[{"update_id":100,"message":{"message_id":9,"chat":{"id":321},"text":"hi from user"}}]`,
	})
	_, _, err := runIn(t, dir, srv, "tok:val", "updates", "get", "--limit", "5")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "tok:val", "log", "--chat", "321", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "hi from user")
	mustContain(t, out, `"direction": "in"`)
}

func TestLog_MarkedReadAndDestructive(t *testing.T) {
	root := NewRootCmd()
	list := findCmd(root, "log")
	require.NotNil(t, list)
	assert.Equal(t, "true", list.Annotations[annReadOnly])

	search := findCmd(root, "log", "search")
	require.NotNil(t, search)
	assert.Equal(t, "true", search.Annotations[annReadOnly])

	show := findCmd(root, "log", "show")
	require.NotNil(t, show)
	assert.Equal(t, "true", show.Annotations[annReadOnly])

	prune := findCmd(root, "log", "prune")
	require.NotNil(t, prune)
	assert.Equal(t, "true", prune.Annotations[annDestructive])
}

func TestLog_NotExcludedFromMCP(t *testing.T) {
	sel := ophis.ExcludeCmdsContaining(excludedFromMCP...)
	for _, p := range [][]string{{"log"}, {"log", "search"}, {"log", "show"}, {"log", "prune"}} {
		cmd := findCmd(NewRootCmd(), p...)
		require.NotNil(t, cmd, "command %v should exist", p)
		assert.True(t, sel(cmd), "log command %v must be exposed as an MCP tool", p)
	}
}
