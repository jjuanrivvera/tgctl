package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// check is one diagnostic line.
type check struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}

func init() {
	register(func(root *cobra.Command) {
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "doctor",
			Short: "Diagnose configuration, credentials, and connectivity",
			Long:  "Run a series of checks (config, token, API reachability, clock) and exit non-zero if any fail.",
			Example: `  tgctl doctor
  tgctl doctor --json`,
			RunE: func(cmd *cobra.Command, _ []string) error {
				checks := runDoctor(cmd)
				allOK := true
				for _, c := range checks {
					if !c.OK {
						allOK = false
					}
				}
				if jsonOut {
					b, _ := json.MarshalIndent(checks, "", "  ")
					fmt.Fprintln(cmd.OutOrStdout(), string(b))
				} else {
					for _, c := range checks {
						mark := "✓"
						if !c.OK {
							mark = "✗"
						}
						line := fmt.Sprintf("%s %s", mark, c.Name)
						if c.Detail != "" {
							line += ": " + c.Detail
						}
						fmt.Fprintln(cmd.OutOrStdout(), line)
					}
				}
				if !allOK {
					return fmt.Errorf("one or more checks failed")
				}
				return nil
			},
		}
		cmd.Flags().BoolVar(&jsonOut, "json", false, "output checks as JSON")
		root.AddCommand(cmd)
	})
}

func runDoctor(cmd *cobra.Command) []check {
	var checks []check

	profileName, cfg, err := resolveProfileName(cmd)
	if err != nil {
		return []check{{Name: "config", OK: false, Detail: err.Error()}}
	}
	checks = append(checks, check{Name: "config loaded", OK: true, Detail: cfg.FilePath()})
	checks = append(checks, check{Name: "active profile", OK: true, Detail: profileName})

	client, err := clientFromCmd(cmd)
	if err != nil {
		checks = append(checks, check{Name: "credentials resolvable", OK: false, Detail: err.Error()})
		return checks
	}
	defer func() { _ = client.Close() }()
	checks = append(checks, check{Name: "credentials resolvable", OK: true, Detail: "token found"})
	checks = append(checks, check{Name: "base URL", OK: true, Detail: client.BaseURL()})

	me, err := client.GetMe(cmd.Context())
	if err != nil {
		checks = append(checks, check{Name: "API reachable + token valid", OK: false, Detail: err.Error()})
	} else {
		checks = append(checks, check{Name: "API reachable + token valid", OK: true, Detail: me.DisplayName()})
	}

	// Clock sanity: a wildly wrong local clock breaks TLS and confuses logs.
	now := time.Now()
	checks = append(checks, check{Name: "local clock", OK: now.Year() >= 2024, Detail: now.Format(time.RFC3339)})

	return checks
}
