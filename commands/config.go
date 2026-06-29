package commands

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/config"
)

func init() {
	register(func(root *cobra.Command) {
		configCmd := &cobra.Command{
			Use:   "config",
			Short: "Inspect and edit tgctl configuration",
			Long:  "View the config file, switch profiles, and set per-profile options. Secrets live in the keyring and are never shown here.",
		}
		configCmd.AddCommand(configPathCmd(), configViewCmd(), configSetCmd(), configUseCmd(), configListProfilesCmd())
		root.AddCommand(configCmd)
	})
}

func configPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := config.Path()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), p)
			return nil
		},
	}
}

func configViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "view",
		Aliases: []string{"show"},
		Short:   "Show the current configuration (secrets redacted)",
		Example: `  tgctl config view
  tgctl config view -o json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			// The token is never stored in the config — only non-secret bits. We still note
			// that a credential lives in the keyring so `view` is self-explanatory.
			view := map[string]any{
				"config_path":     cfg.FilePath(),
				"current_profile": cfg.CurrentProfile,
				"profiles":        cfg.Profiles,
				"aliases":         cfg.Aliases,
				"token_storage":   "OS keyring (run `tgctl auth status` to verify)",
			}
			return render(cmd, mustJSON(view))
		},
	}
}

func configSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a per-profile option (key: base_url)",
		Long:  "Set a non-secret option on the active profile. Supported keys: base_url.",
		Example: `  tgctl config set base_url https://api.telegram.org
  tgctl --profile staging config set base_url http://localhost:8081`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			profileName, cfg, err := resolveProfileName(cmd)
			if err != nil {
				return err
			}
			prof, _ := cfg.Profile(profileName)
			switch key {
			case "base_url", "base-url":
				if err := config.ValidateBaseURL(value); err != nil {
					return err
				}
				prof.BaseURL = value
			default:
				return fmt.Errorf("unknown config key %q (supported: base_url)", key)
			}
			if err := cfg.SetProfile(profileName, prof); err != nil {
				return err
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "set %s=%s on profile %q\n", key, value, profileName)
			return nil
		},
	}
	return cmd
}

func configUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile>",
		Short: "Switch the active profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if _, ok := cfg.Profile(args[0]); !ok {
				return fmt.Errorf("no such profile %q — create it with `tgctl auth login --profile %s`", args[0], args[0])
			}
			cfg.CurrentProfile = args[0]
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "now using profile %q\n", args[0])
			return nil
		},
	}
}

func configListProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list-profiles",
		Aliases: []string{"profiles"},
		Short:   "List configured profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			names := cfg.ProfileNames()
			sort.Strings(names)
			rows := make([]map[string]any, 0, len(names))
			for _, n := range names {
				p, _ := cfg.Profile(n)
				rows = append(rows, map[string]any{
					"profile":  n,
					"current":  n == cfg.CurrentProfile,
					"base_url": config.FirstNonEmpty(p.BaseURL, "https://api.telegram.org"),
					"bot_id":   p.BotID,
				})
			}
			b, _ := json.Marshal(rows)
			return render(cmd, b)
		},
	}
}
