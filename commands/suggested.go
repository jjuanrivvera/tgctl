package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "suggested",
		Short: "Manage suggested posts in channel direct messages",
		Long:  "Approve or decline suggested posts sent to a channel's direct messages chat.",
		Cmds: []methodCmd{
			{
				Use: "approve", Method: "approveSuggestedPost", Kind: kindDestructive,
				Short:   "Approve a suggested post",
				Long:    "Approve a suggested post message for publication in the corresponding channel. Requires the can_post_messages administrator right.",
				Example: "  tgctl suggested approve --chat -100123456789 --message-id 42\n  tgctl suggested approve --chat @channel --message-id 42 --send-date 1780000000",
				Flags: []flagSpec{
					chatFlag(),
					messageIDFlag(),
					{Name: "send-date", Param: "send_date", Kind: flagInt, Usage: "point in time (Unix timestamp) when the post is expected to be published"},
				},
				PreCall: requireValidSendDate,
			},
			{
				Use: "decline", Method: "declineSuggestedPost", Kind: kindDestructive,
				Short:   "Decline a suggested post",
				Long:    "Decline a suggested post message with an optional feedback comment. Requires the can_manage_direct_messages administrator right.",
				Example: "  tgctl suggested decline --chat -100123456789 --message-id 42\n  tgctl suggested decline --chat -100123456789 --message-id 42 --comment \"Off-topic\"",
				Flags: []flagSpec{
					chatFlag(),
					messageIDFlag(),
					{Name: "comment", Usage: "optional feedback comment for the post author (0-128 chars)"},
				},
				PreCall: requireValidDeclineComment,
			},
		},
	})
}

// requireValidSendDate enforces that --send-date is a non-zero, positive timestamp
// before dispatching across the wire, catching mistyped or negative inputs early.
func requireValidSendDate(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	v, ok := params["send_date"]
	if !ok {
		return nil
	}
	date, ok := v.(int64)
	if !ok || date <= 0 {
		return fmt.Errorf("--send-date must be a positive Unix timestamp")
	}
	return nil
}

// requireValidDeclineComment enforces the Bot API limit of 128 characters for decline
// feedback comments locally to avoid wasting an unnecessary roundtrip on validation errors.
func requireValidDeclineComment(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	comment, ok := params["comment"].(string)
	if !ok {
		return nil
	}
	if len([]rune(comment)) > 128 {
		return fmt.Errorf("--comment must be at most 128 characters, got %d", len([]rune(comment)))
	}
	return nil
}
