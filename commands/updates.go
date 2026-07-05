package commands

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "updates",
		Short: "Fetch incoming updates (long polling)",
		Long: `Read updates with getUpdates. Note: getUpdates conflicts with a set webhook —
delete the webhook first (tgctl webhook delete) if you want to poll.`,
		Cmds: []methodCmd{
			{
				Use: "get", Method: "getUpdates", Kind: kindRead,
				Short: "Get pending updates",
				Example: `  tgctl updates get --limit 5
  tgctl updates get --offset 123456789 --timeout 30 -o json
  tgctl updates get --allowed-updates message,callback_query`,
				Flags: []flagSpec{
					{Name: "offset", Kind: flagInt, Usage: "first update id to return (ack earlier ones)"},
					{Name: "limit", Kind: flagInt, Usage: "max updates to return (1-100)"},
					{Name: "timeout", Kind: flagInt, Usage: "long-poll seconds (0 = short poll)"},
					{Name: "allowed-updates", Param: "allowed_updates", Kind: flagStringSlice, Usage: "update types to receive"},
				},
				Columns:     []string{"update_id", "message.message_id", "message.from.username", "message.text"},
				PostSuccess: recordInboundUpdates,
			},
		},
	})
}

// recordInboundUpdates persists each incoming message update (direction 'in') to the local
// store — the counterpart to storeRecorder's direction 'out' — so `tgctl log` sees polled
// updates too. getUpdates has already succeeded and will be rendered right after this runs, so
// a store failure here is logged (inside recordInboundMessage) and swallowed, never surfaced
// as a command error (DECISIONS.md).
func recordInboundUpdates(cmd *cobra.Command, raw json.RawMessage) {
	if len(raw) == 0 {
		return // dry-run, or nothing returned
	}
	st := openStoreForWrite(cmd)
	if st == nil {
		return // disabled (--no-store) or unavailable; openStoreForWrite already warned
	}
	defer func() { _ = st.Close() }()

	var updates []struct {
		Message *telegramMessage `json:"message"`
	}
	if err := json.Unmarshal(raw, &updates); err != nil {
		return // getUpdates always returns this shape; a decode failure here is not our problem to raise
	}
	for _, u := range updates {
		if u.Message != nil {
			recordInboundMessage(cmd, st, u.Message)
		}
	}
}
