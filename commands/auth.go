package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/api"
	"github.com/jjuanrivvera/tgctl/internal/auth"
	"github.com/jjuanrivvera/tgctl/internal/config"
)

func init() {
	register(func(root *cobra.Command) {
		authCmd := &cobra.Command{
			Use:   "auth",
			Short: "Manage bot tokens and verify authentication",
			Long:  "Capture, verify, and remove the bot token for a profile. Tokens are stored in your OS keyring, never in the config file.",
		}
		authCmd.AddCommand(authLoginCmd(), authLogoutCmd(), authStatusCmd())
		root.AddCommand(authCmd)
	})
}

func authLoginCmd() *cobra.Command {
	var token string
	var noVerify bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store a bot token and verify it",
		Long:  "Capture a bot token (from @BotFather), verify it against getMe, and save it to the keyring for the active profile.",
		Example: `  tgctl auth login                      # prompt for the token (hidden input)
  tgctl auth login --token 123:ABC...   # non-interactive
  tgctl auth login --bot staging        # store under a named bot/profile`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profileName, cfg, err := resolveProfileName(cmd)
			if err != nil {
				return err
			}
			if token == "" {
				token, err = promptSecret(cmd, "Bot token (from @BotFather): ")
				if err != nil {
					return err
				}
			}
			authr, err := api.NewBotTokenAuth(token)
			if err != nil {
				return err
			}

			baseFlag, _ := cmd.Flags().GetString("base-url")
			base := config.FirstNonEmpty(baseFlag, api.DefaultBaseURL)
			if err := config.ValidateBaseURL(base); err != nil {
				return err
			}

			botID := authr.BotID()
			if !noVerify {
				client := api.New(authr, api.WithBaseURL(base))
				me, err := client.GetMe(cmd.Context())
				if err != nil {
					return fmt.Errorf("token verification failed: %w", err)
				}
				botID = me.ID.String()
				fmt.Fprintf(cmd.ErrOrStderr(), "verified as %s (id %s)\n", me.DisplayName(), me.ID)
			}

			dir, err := config.Dir()
			if err != nil {
				return err
			}
			if err := auth.New(dir).Set(profileName, token); err != nil {
				return fmt.Errorf("store token: %w", err)
			}
			if err := cfg.SetProfile(profileName, config.Profile{BaseURL: base, AuthMethod: authr.Method(), BotID: botID}); err != nil {
				return err
			}
			if cfg.CurrentProfile == "" {
				cfg.CurrentProfile = profileName
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "logged in to profile %q\n", profileName)
			return nil
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "bot token (omit to be prompted with hidden input)")
	cmd.Flags().BoolVar(&noVerify, "no-verify", false, "skip the getMe verification call")
	return cmd
}

func authLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the stored token for the active profile",
		RunE: func(cmd *cobra.Command, _ []string) error {
			profileName, _, err := resolveProfileName(cmd)
			if err != nil {
				return err
			}
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			if err := auth.New(dir).Delete(profileName); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "logged out of profile %q\n", profileName)
			return nil
		},
	}
}

func authStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"whoami"},
		Short:   "Show the active profile, base URL, and token validity",
		Example: `  tgctl auth status
  tgctl whoami -o json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profileName, cfg, err := resolveProfileName(cmd)
			if err != nil {
				return err
			}
			prof, _ := cfg.Profile(profileName)
			base := config.FirstNonEmpty(prof.BaseURL, api.DefaultBaseURL)

			// auth status is a real check: if the token is missing or invalid it exits
			// non-zero (so `tgctl auth status && …` works), while still printing why.
			client, err := clientFromCmd(cmd)
			if err != nil {
				return fmt.Errorf("not authenticated (profile %q): %w", profileName, err)
			}
			me, err := client.GetMe(cmd.Context())
			if err != nil {
				return fmt.Errorf("token invalid for profile %q: %w", profileName, err)
			}
			status := map[string]any{
				"profile":  profileName,
				"base_url": base,
				"valid":    true,
				"bot":      me.DisplayName(),
				"bot_id":   me.ID.String(),
			}
			return render(cmd, mustJSON(status))
		},
	}
	return cmd
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
