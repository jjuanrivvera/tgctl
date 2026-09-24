package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAudios(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "getUserProfileAudios", &got)

	out, _, err := run(t, srv, "user", "audios", "--user", "12345", "--limit", "5")
	require.NoError(t, err)
	assert.EqualValues(t, 12345, got["user_id"])
	assert.EqualValues(t, 5, got["limit"])
	assert.NotContains(t, got, "offset")
	assert.Contains(t, out, "true")
}

func TestUserChatMessages(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "getUserPersonalChatMessages", &got)

	_, _, err := run(t, srv, "user", "chat-messages", "--user", "12345", "--limit", "20")
	require.NoError(t, err)
	assert.EqualValues(t, 20, got["limit"])
}

// --limit is required by the API, not optional like it is on `user photos`.
func TestUserChatMessages_LimitRequired(t *testing.T) {
	srv := newServer(t, routes{"getUserPersonalChatMessages": `[]`})
	_, _, err := run(t, srv, "user", "chat-messages", "--user", "12345")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "limit")
}

func TestMemberAnswerJoinQuery(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "answerChatJoinRequestQuery", &got)

	_, _, err := run(t, srv, "member", "answer-join-query", "--query-id", "AAxx", "--result", "approve")
	require.NoError(t, err)
	assert.Equal(t, "AAxx", got["chat_join_request_query_id"])
	assert.Equal(t, "approve", got["result"])
}

// The three words are the whole vocabulary of --result: a typo should not cost a round trip,
// least of all in a join flow where someone is waiting on the other side.
func TestMemberAnswerJoinQuery_RejectsUnknownResult(t *testing.T) {
	srv := newServer(t, routes{"answerChatJoinRequestQuery": `true`})

	for _, bad := range []string{"accept", "reject", "APPROVE", ""} {
		t.Run("result="+bad, func(t *testing.T) {
			_, _, err := run(t, srv, "member", "answer-join-query", "--query-id", "AAxx", "--result", bad)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "approve, decline or queue")
		})
	}

	for _, ok := range []string{"approve", "decline", "queue"} {
		t.Run("result="+ok, func(t *testing.T) {
			_, _, err := run(t, srv, "member", "answer-join-query", "--query-id", "AAxx", "--result", ok)
			require.NoError(t, err)
		})
	}
}
