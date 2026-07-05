package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/config"
	"github.com/jjuanrivvera/tgctl/internal/store"
)

// logDefaultColumns orders the common case (a wide table of everything) sensibly; -o json/yaml
// always carries every field regardless, and --columns overrides this like any other command.
var logDefaultColumns = []string{"id", "ts", "direction", "chat_id", "kind", "text", "message_id"}

func init() {
	register(func(root *cobra.Command) {
		var chatID int64
		var since, kind string
		var limit int

		logCmd := &cobra.Command{
			Use:   "log",
			Short: "Query tgctl's local send/receive history",
			Long: `tgctl records every outbound send (and, in polling/webhook mode, inbound updates)
to a local SQLite database — one per bot profile — because the Bot API itself has no history
endpoint. This lets a restarted or compacted session, or any external tool, answer "what did
you send/receive, when, to whom". Disable recording for a single call with --no-store; this
command itself always reads regardless of --no-store (it does not write).`,
			Example: `  tgctl log
  tgctl log --chat 123456789 --since 24h
  tgctl log --kind photo --limit 20 -o json
  tgctl log search "deploy failed"
  tgctl log show 42
  tgctl log prune --older-than 2160h`,
			Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				f, err := buildLogFilter(chatID, since, kind, limit)
				if err != nil {
					return err
				}
				return withReadStore(cmd, func(st *store.Store) error {
					msgs, err := st.Query(cmd.Context(), f)
					if err != nil {
						return err
					}
					return renderMessages(cmd, msgs)
				})
			},
		}
		bindLogFilterFlags(logCmd, &chatID, &since, &kind, &limit)
		markKind(logCmd, kindRead)

		search, show, prune := logSearchCmd(), logShowCmd(), logPruneCmd()
		markKind(search, kindRead)
		markKind(show, kindRead)
		markKind(prune, kindDestructive)
		logCmd.AddCommand(search, show, prune)
		root.AddCommand(logCmd)
	})
}

func bindLogFilterFlags(cmd *cobra.Command, chatID *int64, since, kind *string, limit *int) {
	cmd.Flags().Int64Var(chatID, "chat", 0, "filter by chat id")
	cmd.Flags().StringVar(since, "since", "", "only messages at/after this time: a Go duration (24h) or RFC3339/YYYY-MM-DD")
	cmd.Flags().StringVar(kind, "kind", "", "filter by kind: text|photo|document|voice|edit|...")
	cmd.Flags().IntVar(limit, "limit", store.DefaultLimit, "max rows to return")
}

func buildLogFilter(chatID int64, since, kind string, limit int) (store.Filter, error) {
	t, err := parseSince(since)
	if err != nil {
		return store.Filter{}, err
	}
	return store.Filter{ChatID: chatID, Since: t, Kind: kind, Limit: limit}, nil
}

// parseSince accepts a Go duration ("24h", meaning "since 24h ago") or an absolute timestamp
// (RFC3339 or a bare YYYY-MM-DD date), matching the issue's `--since 24h` example while still
// allowing a precise cutoff.
func parseSince(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if d, err := time.ParseDuration(s); err == nil {
		return time.Now().UTC().Add(-d), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("invalid --since %q (want a Go duration like 24h, or RFC3339/YYYY-MM-DD)", s)
}

func logSearchCmd() *cobra.Command {
	var chatID int64
	var since, kind string
	var limit int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search recorded message/caption text",
		Long: `Search uses FTS5 MATCH when the linked SQLite build supports it (operators: AND/OR/
NOT, prefix*, "phrases"); otherwise it degrades to a plain substring scan automatically — check
"tgctl doctor" or the store's FTSEnabled to see which mode is active.`,
		Example: `  tgctl log search "deploy failed"
  tgctl log search "deploy* AND staging" --chat 123456789`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := buildLogFilter(chatID, since, kind, limit)
			if err != nil {
				return err
			}
			return withReadStore(cmd, func(st *store.Store) error {
				msgs, err := st.Search(cmd.Context(), args[0], f)
				if err != nil {
					return err
				}
				return renderMessages(cmd, msgs)
			})
		},
	}
	bindLogFilterFlags(cmd, &chatID, &since, &kind, &limit)
	return cmd
}

func logShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <message_id>",
		Short: "Show one recorded message, including its full raw API payload",
		Example: `  tgctl log show 42
  tgctl log show 42 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			messageID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid message id %q: %w", args[0], err)
			}
			return withReadStore(cmd, func(st *store.Store) error {
				msg, ok, err := st.Show(cmd.Context(), messageID)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("no recorded message with message_id %d", messageID)
				}
				return render(cmd, mustJSON(msg))
			})
		},
	}
	return cmd
}

func logPruneCmd() *cobra.Command {
	var olderThan string
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Delete recorded messages older than a duration",
		Example: `  tgctl log prune --older-than 2160h   # 90 days
  tgctl log prune --older-than 720h    # 30 days`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := time.ParseDuration(olderThan)
			if err != nil {
				return fmt.Errorf("invalid --older-than %q (want a Go duration like 2160h): %w", olderThan, err)
			}
			return withReadStore(cmd, func(st *store.Store) error {
				n, err := st.Prune(cmd.Context(), d)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "pruned %d message(s) older than %s\n", n, olderThan)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&olderThan, "older-than", "", "delete messages recorded before now minus this Go duration (required)")
	_ = cmd.MarkFlagRequired("older-than")
	return cmd
}

// withReadStore opens the active profile's store for a `log` subcommand and always closes it
// afterward. Unlike the write-path's openStoreForWrite, a failure here IS a real command error:
// reading local history is what `tgctl log` exists to do, so silently returning "no messages"
// on an unopenable store would be misleading rather than merely degraded.
func withReadStore(cmd *cobra.Command, fn func(*store.Store) error) error {
	profileName, _, err := resolveProfileName(cmd)
	if err != nil {
		return err
	}
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	path, err := store.PathFor(dir, profileName)
	if err != nil {
		return err
	}
	st, err := store.Open(path)
	if err != nil {
		return fmt.Errorf("open message store: %w", err)
	}
	defer func() { _ = st.Close() }()
	return fn(st)
}

// renderMessages applies the default column order (unless the user set --columns) then renders.
func renderMessages(cmd *cobra.Command, msgs []store.Message) error {
	if !cmd.Flags().Changed("columns") {
		if err := cmd.Flags().Set("columns", strings.Join(logDefaultColumns, ",")); err != nil {
			return err
		}
	}
	return render(cmd, mustJSON(msgs))
}
