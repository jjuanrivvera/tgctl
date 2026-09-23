package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// paramProbe records the JSON body of the last call, so a test can assert the exact wire
// parameters a verb sends rather than just that it succeeded.
func paramProbe(t *testing.T, method string, got *map[string]any) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.True(t, strings.HasSuffix(r.URL.Path, "/"+method), "unexpected call to %s", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(got))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// Removing one member's reaction from one message (issue #29).
func TestMessageUnreact_ByUser(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "deleteMessageReaction", &got)

	out, _, err := run(t, srv, "message", "unreact", "--chat", "@group", "--message-id", "42", "--user", "12345")
	require.NoError(t, err)
	assert.Contains(t, out, "true")
	assert.Equal(t, "@group", got["chat_id"])
	assert.EqualValues(t, 42, got["message_id"])
	assert.EqualValues(t, 12345, got["user_id"])
	assert.NotContains(t, got, "actor_chat_id", "an unset flag must not be sent")
}

// The same verb, for a reaction a channel left on its own behalf.
func TestMessageUnreact_ByActorChat(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "deleteMessageReaction", &got)

	_, _, err := run(t, srv, "message", "unreact", "--chat", "@group", "--message-id", "7", "--actor-chat", "-1001234567890")
	require.NoError(t, err)
	assert.EqualValues(t, -1001234567890, got["actor_chat_id"])
	assert.NotContains(t, got, "user_id")
}

// The chat-wide sweep takes no message id: it is about an actor, not a message.
func TestMessageUnreactAll(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "deleteAllMessageReactions", &got)

	_, _, err := run(t, srv, "message", "unreact-all", "--chat", "@group", "--user", "12345")
	require.NoError(t, err)
	assert.Equal(t, "@group", got["chat_id"])
	assert.EqualValues(t, 12345, got["user_id"])
	assert.NotContains(t, got, "message_id")
}

func TestMemberSetTag(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "setChatMemberTag", &got)

	_, _, err := run(t, srv, "member", "set-tag", "--chat", "@group", "--user", "12345", "--tag", "Moderator")
	require.NoError(t, err)
	assert.Equal(t, "Moderator", got["tag"])
	assert.EqualValues(t, 12345, got["user_id"])
}

// Omitting --tag clears it, so the parameter must be absent rather than an empty string —
// only flags the user actually set are sent.
func TestMemberSetTag_OmittedClears(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "setChatMemberTag", &got)

	_, _, err := run(t, srv, "member", "set-tag", "--chat", "@group", "--user", "12345")
	require.NoError(t, err)
	assert.NotContains(t, got, "tag")
}

// Both reaction verbs remove someone else's content, so the agent guard and the MCP server
// must see them as destructive.
func TestReactionVerbs_AreDestructive(t *testing.T) {
	want := map[string]bool{"deleteMessageReaction": false, "deleteAllMessageReactions": false}
	for _, c := range APICommands() {
		if _, tracked := want[c.Method]; tracked {
			want[c.Method] = c.IsDestructive()
		}
	}
	for method, destructive := range want {
		assert.True(t, destructive, "%s should be classified destructive", method)
	}
}

// A reaction has exactly one author, so naming none or naming two is a question the API was
// never asked. Answer it locally, with the flag to use, instead of relaying a 400.
func TestUnreact_RequiresExactlyOneActor(t *testing.T) {
	srv := newServer(t, routes{
		"deleteMessageReaction":     `true`,
		"deleteAllMessageReactions": `true`,
	})

	cases := []struct {
		name, wantErr string
		args          []string
	}{{
		name:    "unreact with no actor",
		args:    []string{"message", "unreact", "--chat", "@g", "--message-id", "1"},
		wantErr: "one of --user",
	}, {
		name:    "unreact with both",
		args:    []string{"message", "unreact", "--chat", "@g", "--message-id", "1", "--user", "1", "--actor-chat", "-100"},
		wantErr: "not both",
	}, {
		name:    "unreact-all with no actor",
		args:    []string{"message", "unreact-all", "--chat", "@g"},
		wantErr: "one of --user",
	}, {
		name:    "unreact-all with both",
		args:    []string{"message", "unreact-all", "--chat", "@g", "--user", "1", "--actor-chat", "-100"},
		wantErr: "not both",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, srv, tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
