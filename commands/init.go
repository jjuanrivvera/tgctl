package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/api"
	"github.com/jjuanrivvera/tgctl/internal/auth"
	"github.com/jjuanrivvera/tgctl/internal/config"
)

func init() {
	register(func(root *cobra.Command) {
		cmd := &cobra.Command{
			Use:     "init",
			Aliases: []string{"setup"},
			Short:   "First-run wizard: pick a base URL, capture a token, and smoke-test",
			Long:    "Interactively set up a profile: choose the base URL (default https://api.telegram.org), paste a bot token, verify it, and store it in the keyring.",
			Example: `  tgctl init
  tgctl init --profile staging`,
			RunE: func(cmd *cobra.Command, _ []string) error {
				profileName, cfg, err := resolveProfileName(cmd)
				if err != nil {
					return err
				}
				out := cmd.OutOrStdout()
				fmt.Fprintf(cmd.ErrOrStderr(), "Setting up profile %q.\n", profileName)

				base, err := promptLine(cmd, "Base URL [https://api.telegram.org]: ")
				if err != nil {
					return err
				}
				base = config.FirstNonEmpty(base, api.DefaultBaseURL)
				if err := config.ValidateBaseURL(base); err != nil {
					return err
				}

				token, err := promptSecret(cmd, "Bot token (from @BotFather): ")
				if err != nil {
					return err
				}
				authr, err := api.NewBotTokenAuth(token)
				if err != nil {
					return err
				}

				client := api.New(authr, api.WithBaseURL(base))
				me, err := client.GetMe(cmd.Context())
				if err != nil {
					return fmt.Errorf("smoke test failed (token or connectivity): %w", err)
				}

				dir, err := config.Dir()
				if err != nil {
					return err
				}
				if err := auth.New(dir).Set(profileName, token); err != nil {
					return err
				}
				if err := cfg.SetProfile(profileName, config.Profile{BaseURL: base, AuthMethod: authr.Method(), BotID: me.ID.String()}); err != nil {
					return err
				}
				cfg.CurrentProfile = profileName
				if err := cfg.Save(); err != nil {
					return err
				}

				fmt.Fprintf(out, "✓ Profile %q ready — authenticated as %s (id %s)\n", profileName, me.DisplayName(), me.ID)
				fmt.Fprintln(out, "  Try: tgctl bot info")
				return nil
			},
		}
		root.AddCommand(cmd)
	})
}
