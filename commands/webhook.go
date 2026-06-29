package commands

import "github.com/spf13/cobra"

func init() {
	registerGroup(group{
		Use:   "webhook",
		Short: "Manage the bot's webhook",
		Long:  "Inspect, set, and delete the webhook used to receive updates over HTTPS (instead of polling).",
		Cmds: []methodCmd{
			{
				Use: "info", Method: "getWebhookInfo", Kind: kindRead,
				Short:   "Show the current webhook status",
				Example: `  tgctl webhook info -o json`,
				Columns: []string{"url", "pending_update_count", "last_error_message"},
			},
			{
				Use: "set", Method: "setWebhook", Kind: kindWrite,
				Short: "Set the webhook URL",
				Example: `  tgctl webhook set --url https://example.com/bot --max-connections 40
  tgctl webhook set --url https://example.com/bot --secret-token s3cr3t --drop-pending`,
				Flags: []flagSpec{
					{Name: "url", Required: true, Usage: "HTTPS URL to receive updates"},
					{Name: "secret-token", Param: "secret_token", Usage: "secret echoed in the X-Telegram-Bot-Api-Secret-Token header"},
					{Name: "max-connections", Param: "max_connections", Kind: flagInt, Usage: "max concurrent connections (1-100)"},
					{Name: "allowed-updates", Param: "allowed_updates", Kind: flagStringSlice, Usage: "update types to receive"},
					{Name: "drop-pending", Param: "drop_pending_updates", Kind: flagBool, Usage: "drop queued updates"},
				},
			},
			{
				Use: "delete", Method: "deleteWebhook", Kind: kindDestructive,
				Short:   "Delete the webhook (switch back to polling)",
				Example: `  tgctl webhook delete --drop-pending`,
				Flags: []flagSpec{
					{Name: "drop-pending", Param: "drop_pending_updates", Kind: flagBool, Usage: "drop queued updates"},
				},
			},
		},
		// `listen` is a local receiver (a value-add beyond the API), not a single Bot API
		// method — so it's a hand-written Extra command, not part of the generic surface.
		Extra: []func() *cobra.Command{webhookListenCmd},
	})
}
