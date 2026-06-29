package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	register(func(root *cobra.Command) {
		var data string
		var query []string
		var idempotent bool
		cmd := &cobra.Command{
			Use:   "api <method> [-d body] [-q key=value ...]",
			Short: "Call any Bot API method directly (raw escape hatch)",
			Long: `Invoke an arbitrary Bot API method with a JSON body and/or key=value parameters.

This is the documented escape hatch for methods tgctl does not wrap as first-class
commands. It honors --dry-run and -o/--output like every other command. By default a
raw call is treated as a write (not auto-retried); pass --idempotent for read-only
methods (getX) so transient failures retry safely.`,
			Example: `  tgctl api getMe
  tgctl api sendMessage -q chat_id=@me -q text="hi from tgctl"
  tgctl api getChat -q chat_id=@telegram --idempotent
  tgctl api sendMessage -d '{"chat_id":"@me","text":"json body"}'`,
			Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				method := args[0]
				params := map[string]any{}
				if data != "" {
					if err := json.Unmarshal([]byte(data), &params); err != nil {
						return fmt.Errorf("invalid -d JSON body: %w", err)
					}
				}
				for _, kv := range query {
					k, v, ok := strings.Cut(kv, "=")
					if !ok {
						return fmt.Errorf("invalid -q %q (want key=value)", kv)
					}
					params[k] = v
				}
				client, err := clientFromCmd(cmd)
				if err != nil {
					return err
				}
				raw, err := client.Call(cmd.Context(), method, params, idempotent)
				if err != nil {
					return err
				}
				return render(cmd, raw)
			},
		}
		cmd.Flags().StringVarP(&data, "data", "d", "", "raw JSON request body")
		cmd.Flags().StringArrayVarP(&query, "query", "q", nil, "key=value parameter (repeatable)")
		cmd.Flags().BoolVar(&idempotent, "idempotent", false, "treat as read-only (safe to auto-retry)")
		root.AddCommand(cmd)
	})
}
