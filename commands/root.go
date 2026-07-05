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
	"github.com/jjuanrivvera/tgctl/internal/store"
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
	groupCommands(root)
	return root
}

// commandGroups assigns each top-level command to a gh-style section so `--help` reads like a
// first-party tool instead of one flat alphabetical list.
var commandGroups = map[string]string{
	"message": "messaging", "media": "messaging", "file": "messaging",
	"callback": "messaging", "inline": "messaging",
	"chat": "chats", "member": "chats", "invite": "chats", "user": "chats", "updates": "chats",
	"forum": "chats", "verify": "chats",
	"bot": "config", "commands": "config", "webhook": "config", "stars": "config",
	"auth": "meta", "config": "meta", "init": "meta", "doctor": "meta", "log": "meta",
	"alias": "meta", "api": "meta", "version": "meta", "completion": "meta",
	"mcp": "agents", "agent": "agents",
}

func groupCommands(root *cobra.Command) {
	root.AddGroup(
		&cobra.Group{ID: "messaging", Title: "Messaging:"},
		&cobra.Group{ID: "chats", Title: "Chats & members:"},
		&cobra.Group{ID: "config", Title: "Bot configuration:"},
		&cobra.Group{ID: "agents", Title: "AI agents:"},
		&cobra.Group{ID: "meta", Title: "Setup & meta:"},
	)
	for _, c := range root.Commands() {
		if id, ok := commandGroups[c.Name()]; ok {
			c.GroupID = id
		}
	}
}

func addGlobalFlags(root *cobra.Command) {
	pf := root.PersistentFlags()
	pf.StringP("output", "o", "table", "output format: table|json|yaml|csv|id")
	// A tgctl "profile" is one bot, so the user-facing flag is --bot. --profile stays as a
	// hidden, still-working alias so existing scripts don't break.
	pf.String("bot", "", "bot to use: a named profile/credential (env TGCTL_BOT)")
	pf.String("profile", "", "deprecated alias for --bot")
	_ = pf.MarkHidden("profile")
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
	pf.Bool("no-store", false, "disable local SQLite send/receive history for this invocation (see tgctl log)")
}

// clientFromCmd builds an API client from the resolved profile, token, and global flags.
// Token precedence: $TGCTL_TOKEN > $TELEGRAM_BOT_TOKEN > the profile's keyring entry.
//
// Every caller must `defer client.Close()` once err is nil: the client may hold an open message
// store file handle (its Recorder), and Close is always safe to call (a no-op when there is no
// recorder to close).
func clientFromCmd(cmd *cobra.Command) (*api.Client, error) {
	f := cmd.Flags()
	profileFlag := resolveBotFlag(cmd)
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
	// The store hook is best-effort and additive: a disabled/unavailable store never prevents
	// building a client, so a send still works even when local history doesn't (DECISIONS.md).
	// Dry-run makes no API call at all, so there is nothing to record — skip opening the store
	// entirely rather than create a DB file (and a handle every caller must remember to close)
	// for a command that will never write to it.
	if !dryRun {
		if st := openStoreForWrite(cmd); st != nil {
			quiet, _ := f.GetBool("quiet")
			opts = append(opts, api.WithRecorder(&storeRecorder{st: st, quiet: quiet}))
		}
	}
	return api.New(authr, opts...), nil
}

// openStoreForWrite opens the active profile's local message store for the write/record path:
// outbound sends via storeRecorder (above) and inbound updates via
// commands/updates.go/commands/webhook_listen.go, honoring --no-store. It resolves the profile
// itself (rather than taking one as a parameter) so every write-path call site — including the
// ones that never build an api.Client — shares this one helper. Unlike openStoreForRead
// (commands/log.go), failure here is never fatal: nil means "recording is disabled for this
// call", and callers must keep going without it.
func openStoreForWrite(cmd *cobra.Command) *store.Store {
	if noStore, _ := cmd.Flags().GetBool("no-store"); noStore {
		return nil
	}
	profileName, _, err := resolveProfileName(cmd)
	if err != nil {
		warnStoreUnavailable(cmd, err)
		return nil
	}
	dir, err := config.Dir()
	if err != nil {
		warnStoreUnavailable(cmd, err)
		return nil
	}
	path, err := store.PathFor(dir, profileName)
	if err != nil {
		warnStoreUnavailable(cmd, err)
		return nil
	}
	st, err := store.Open(path)
	if err != nil {
		warnStoreUnavailable(cmd, err)
		return nil
	}
	return st
}

// warnStoreUnavailable notes a store-open failure on stderr, respecting --quiet. It is never
// an error returned to the caller: see openStoreForWrite.
func warnStoreUnavailable(cmd *cobra.Command, err error) {
	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		return
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "tgctl: warning: message store unavailable (%v) — continuing without local history\n", err)
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
	return cfg.ResolveProfileName(resolveBotFlag(cmd)), cfg, nil
}

// resolveBotFlag returns the selected bot from the --bot flag, falling back to the deprecated
// --profile alias, so existing scripts that pass --profile keep working unchanged.
func resolveBotFlag(cmd *cobra.Command) string {
	if bot, _ := cmd.Flags().GetString("bot"); bot != "" {
		return bot
	}
	prof, _ := cmd.Flags().GetString("profile")
	return prof
}
