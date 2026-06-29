// Package commands wires the Telegram Bot API client (internal/api) into a Cobra command
// tree. Command groups self-register via register() so a fresh tree can be built per process
// (and per test). Shared concerns — client construction and rendering — live here once.
package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/api"
	"github.com/jjuanrivvera/tgctl/internal/auth"
	"github.com/jjuanrivvera/tgctl/internal/config"
	"github.com/jjuanrivvera/tgctl/internal/output"
)

// registrations are applied to each fresh root command. Command files append to it in init().
var registrations []func(*cobra.Command)

func register(fn func(*cobra.Command)) { registrations = append(registrations, fn) }

const rootLong = `tgctl is a fast, scriptable command-line tool for the Telegram Bot API.

It wraps the Bot API methods (sendMessage, getChat, getUpdates, ...) behind ergonomic
commands with table/json/yaml/csv output, named profiles for multiple bots, OS-keyring
token storage, and an MCP server so AI agents can drive it safely.

Get a bot token from @BotFather, then:

  tgctl auth login                       # store the token in your OS keyring
  tgctl bot info                         # who am I?
  tgctl message send --chat @me --text "hello from tgctl"
  tgctl updates get --limit 5 -o json    # poll recent updates as JSON

Every command honors --dry-run (prints the equivalent curl), -o/--output, and --jq.`

// NewRootCmd builds a fresh command tree with all registered groups attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "tgctl",
		Short:         "Command-line tool for the Telegram Bot API",
		Long:          rootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       "", // set by cmd/tgctl via SetVersionTemplate; version cmd has the detail
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("output")
			if !output.Format(f).Valid() {
				return fmt.Errorf("invalid --output %q (want table|json|yaml|csv|id)", f)
			}
			return nil
		},
	}
	addGlobalFlags(root)
	for _, fn := range registrations {
		fn(root)
	}
	return root
}

func addGlobalFlags(root *cobra.Command) {
	pf := root.PersistentFlags()
	pf.StringP("output", "o", "table", "output format: table|json|yaml|csv|id")
	pf.String("profile", "", "profile/instance to use (env TGCTL_PROFILE)")
	pf.String("base-url", "", "Bot API base URL (default https://api.telegram.org)")
	pf.Bool("dry-run", false, "print the equivalent curl and make no request")
	pf.Bool("show-token", false, "do not redact the bot token in --dry-run output")
	pf.BoolP("verbose", "v", false, "log raw API responses to stderr")
	pf.Bool("no-color", false, "disable colored output")
	pf.StringSlice("columns", nil, "explicit, ordered table/csv columns")
	// --quiet has no -q short: the `api` escape hatch uses -q for repeatable key=value params.
	pf.Bool("quiet", false, "suppress notes on stderr")
	pf.String("jq", "", "gojq expression applied to the result before rendering")
	pf.Float64("rps", 0, "client-side requests-per-second cap (0 = default)")
}

// clientFromCmd builds an API client from the resolved profile, token, and global flags.
// Token precedence: $TGCTL_TOKEN > $TELEGRAM_BOT_TOKEN > the profile's keyring entry.
func clientFromCmd(cmd *cobra.Command) (*api.Client, error) {
	f := cmd.Flags()
	profileFlag, _ := f.GetString("profile")
	baseURLFlag, _ := f.GetString("base-url")
	dryRun, _ := f.GetBool("dry-run")
	showToken, _ := f.GetBool("show-token")
	verbose, _ := f.GetBool("verbose")
	rps, _ := f.GetFloat64("rps")

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	profileName := cfg.ResolveProfileName(profileFlag)
	prof, _ := cfg.Profile(profileName)

	token := config.FirstNonEmpty(os.Getenv("TGCTL_TOKEN"), os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		dir, derr := config.Dir()
		if derr != nil {
			return nil, derr
		}
		token, err = auth.New(dir).Get(profileName)
		if err != nil {
			return nil, fmt.Errorf("no bot token for profile %q — run `tgctl auth login` (or set $TGCTL_TOKEN)", profileName)
		}
	}

	authr, err := api.NewBotTokenAuth(token)
	if err != nil {
		return nil, err
	}

	baseURL := config.FirstNonEmpty(baseURLFlag, prof.BaseURL, api.DefaultBaseURL)
	if err := config.ValidateBaseURL(baseURL); err != nil {
		return nil, err
	}

	opts := []api.Option{
		api.WithBaseURL(baseURL),
		api.WithDryRun(dryRun),
		api.WithShowToken(showToken),
		api.WithVerbose(verbose),
		api.WithDryRunWriter(cmd.ErrOrStderr()),
	}
	if rps > 0 {
		opts = append(opts, api.WithRPS(rps))
	}
	return api.New(authr, opts...), nil
}

// render writes data using the format/columns/jq/quiet flags, to the command's streams.
func render(cmd *cobra.Command, data json.RawMessage) error {
	f := cmd.Flags()
	format, _ := f.GetString("output")
	columns, _ := f.GetStringSlice("columns")
	noColor, _ := f.GetBool("no-color")
	quiet, _ := f.GetBool("quiet")
	jq, _ := f.GetString("jq")
	return output.Render(data, output.Options{
		Format:  output.Format(format),
		Columns: columns,
		NoColor: noColor,
		Quiet:   quiet,
		JQ:      jq,
		Out:     cmd.OutOrStdout(),
		Err:     cmd.ErrOrStderr(),
	})
}

// resolveProfileName returns the active profile name from flags/env/config without building
// a client — used by auth/config/doctor commands.
func resolveProfileName(cmd *cobra.Command) (string, *config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", nil, err
	}
	pf, _ := cmd.Flags().GetString("profile")
	return cfg.ResolveProfileName(pf), cfg, nil
}
