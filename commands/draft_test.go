package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageDraft(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "sendMessageDraft", &got)

	_, _, err := run(t, srv, "message", "draft", "--chat", "12345", "--draft-id", "7",
		"--text", "Looking that up", "--can-stop")
	require.NoError(t, err)
	assert.EqualValues(t, 12345, got["chat_id"])
	assert.EqualValues(t, 7, got["draft_id"])
	assert.Equal(t, "Looking that up", got["text"])
	assert.Equal(t, true, got["can_stop"])
	assert.NotContains(t, got, "keep_on_stop", "an unset flag must not be sent")
}

// An empty --text is meaningful here: it shows the "Thinking…" placeholder. It must travel as
// an empty string rather than be dropped as unset.
func TestMessageDraft_EmptyTextIsSent(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "sendMessageDraft", &got)

	_, _, err := run(t, srv, "message", "draft", "--chat", "1", "--draft-id", "2", "--text", "")
	require.NoError(t, err)
	text, present := got["text"]
	require.True(t, present, "--text set to empty must still be sent: it is the placeholder")
	assert.Equal(t, "", text)
}

// The API states the id must be non-zero, and zero is what you get by forgetting what the
// flag means rather than by choosing it.
func TestMessageDraft_RejectsZeroDraftID(t *testing.T) {
	srv := newServer(t, routes{"sendMessageDraft": `true`})
	_, _, err := run(t, srv, "message", "draft", "--chat", "1", "--draft-id", "0", "--text", "hi")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-zero")
}

// answer-guest takes ONE result object, where `inline answer` takes an array. Sending the
// wrong shape is the easy mistake, so the flag is documented as an object and travels as one.
func TestInlineAnswerGuest(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "answerGuestQuery", &got)

	_, _, err := run(t, srv, "inline", "answer-guest", "--query-id", "AAxx",
		"--result", `{"type":"article","id":"1","title":"Hi","input_message_content":{"message_text":"Hi"}}`)
	require.NoError(t, err)
	assert.Equal(t, "AAxx", got["guest_query_id"])
	result, ok := got["result"].(map[string]any)
	require.True(t, ok, "result must be sent as an object, not an array or a string: %#v", got["result"])
	assert.Equal(t, "article", result["type"])
}

func TestInlineAnswerGuest_RejectsInvalidJSON(t *testing.T) {
	srv := newServer(t, routes{"answerGuestQuery": `{}`})
	_, _, err := run(t, srv, "inline", "answer-guest", "--query-id", "AAxx", "--result", "{not json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "valid JSON")
}

// `inline answer` takes an array of results and `answer-guest` takes one object. Pasting the
// array over here is the obvious slip between neighbours, and JSON validity alone lets it
// through — so the check is on the shape, not just the syntax.
func TestInlineAnswerGuest_RejectsAnArray(t *testing.T) {
	srv := newServer(t, routes{"answerGuestQuery": `{}`})

	_, _, err := run(t, srv, "inline", "answer-guest", "--query-id", "AAxx",
		"--result", `[{"type":"article","id":"1","title":"Hi"}]`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not an array")

	_, _, err = run(t, srv, "inline", "answer-guest", "--query-id", "AAxx", "--result", "null")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "got null")
}
