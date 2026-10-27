package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuggestedApprove_WireParams(t *testing.T) {
	t.Run("minimal with numeric chat", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "approveSuggestedPost", &got)

		out, _, err := run(t, srv, "suggested", "approve", "--chat", "-100123456789", "--message-id", "42")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-100123456789", got["chat_id"])
		assert.EqualValues(t, 42, got["message_id"])
		assert.NotContains(t, got, "send_date")
	})

	t.Run("with channel username and send-date", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "approveSuggestedPost", &got)

		out, _, err := run(t, srv, "suggested", "approve", "--chat", "@channel", "--message-id", "42", "--send-date", "1780000000")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "@channel", got["chat_id"])
		assert.EqualValues(t, 42, got["message_id"])
		assert.EqualValues(t, 1780000000, got["send_date"])
	})
}

func TestSuggestedDecline_WireParams(t *testing.T) {
	t.Run("without comment", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "declineSuggestedPost", &got)

		out, _, err := run(t, srv, "suggested", "decline", "--chat", "-100123456789", "--message-id", "42")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-100123456789", got["chat_id"])
		assert.EqualValues(t, 42, got["message_id"])
		assert.NotContains(t, got, "comment")
	})

	t.Run("with comment", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "declineSuggestedPost", &got)

		out, _, err := run(t, srv, "suggested", "decline", "--chat", "-100123456789", "--message-id", "42", "--comment", "Off-topic submission")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "-100123456789", got["chat_id"])
		assert.EqualValues(t, 42, got["message_id"])
		assert.Equal(t, "Off-topic submission", got["comment"])
	})
}

func TestSuggestedValidation(t *testing.T) {
	tooLongComment := strings.Repeat("a", 129)
	cases := []struct {
		name string
		args []string
	}{
		{"approve missing chat", []string{"suggested", "approve", "--message-id", "42"}},
		{"approve missing message-id", []string{"suggested", "approve", "--chat", "@channel"}},
		{"approve zero send-date", []string{"suggested", "approve", "--chat", "@channel", "--message-id", "42", "--send-date", "0"}},
		{"approve negative send-date", []string{"suggested", "approve", "--chat", "@channel", "--message-id", "42", "--send-date", "-1"}},
		{"decline missing chat", []string{"suggested", "decline", "--message-id", "42"}},
		{"decline missing message-id", []string{"suggested", "decline", "--chat", "@channel"}},
		{"decline comment too long", []string{"suggested", "decline", "--chat", "@channel", "--message-id", "42", "--comment", tooLongComment}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
		})
	}
}

func TestSuggestedKindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	for _, tc := range []struct {
		path       string
		method     string
		adminRight string
	}{
		{"suggested approve", "approveSuggestedPost", "can_post_messages"},
		{"suggested decline", "declineSuggestedPost", "can_manage_direct_messages"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			found := false
			for _, c := range APICommands() {
				if c.Path == tc.path {
					found = true
					assert.Equal(t, kindDestructive, c.Kind)
					assert.Equal(t, tc.method, c.Method)
				}
			}
			assert.True(t, found, "path must be registered in APICommands")

			cmd, _, err := root.Find(strings.Fields(tc.path))
			require.NoError(t, err)
			assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
			assert.Equal(t, "true", cmd.Annotations[annDestructive])
			assert.Empty(t, cmd.Annotations[annReadOnly])
			assert.Contains(t, cmd.Long, tc.adminRight, "help must document required admin right")
		})
	}
}
