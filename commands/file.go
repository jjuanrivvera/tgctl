package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "file",
		Short: "Inspect and download files",
		Long:  "Resolve a file_id to its metadata (getFile) and download the file's bytes to disk.",
		Cmds: []methodCmd{
			{
				Use: "info", Method: "getFile", Kind: kindRead,
				Short:   "Resolve a file_id to its path and size (getFile)",
				Example: `  tgctl file info --file-id BAADBAADrwADBREAAYag...`,
				Flags: []flagSpec{
					{Name: "file-id", Param: "file_id", Required: true, Usage: "the file_id from a message's document/photo/audio/etc."},
				},
				Columns: []string{"file_id", "file_path", "file_size"},
			},
		},
		// `download` resolves the file and streams its bytes locally — two steps (getFile then a
		// /file/ GET), not a single Bot API method — so it's a hand-written Extra, like
		// `webhook listen`. It is not in api-manifest.json (the manifest tracks pure API methods).
		Extra: []func() *cobra.Command{fileDownloadCmd},
	})
}

// fileDownloadCmd resolves a file_id with getFile, then streams the file from the Bot API's
// /file/ endpoint to a local path (or stdout). Files larger than 20MB can't be fetched via the
// Bot API and have no file_path; we surface that explicitly rather than writing an empty file.
func fileDownloadCmd() *cobra.Command {
	var fileID, dest string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a file by file_id to a local path",
		Long: `Resolve a file_id with getFile and download its bytes. The destination defaults to the
file's base name in the current directory; pass --dest to choose a path, or --dest - to write
to stdout. Honors --dry-run (it prints the getFile request without downloading).`,
		Example: `  tgctl file download --file-id BAADBAADrwAD...
  tgctl file download --file-id BAADBAADrwAD... --dest ./photo.jpg
  tgctl file download --file-id BAADBAADrwAD... --dest - > out.bin`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := clientFromCmd(cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()
			ctx := cmd.Context()
			raw, err := client.Call(ctx, "getFile", map[string]any{"file_id": fileID}, true)
			if err != nil {
				return err
			}
			if dry, _ := cmd.Flags().GetBool("dry-run"); dry {
				// getFile already printed its curl; note the second step a real run would take.
				fmt.Fprintln(cmd.ErrOrStderr(), "dry-run: would then GET the file's /file/ URL and write it to disk")
				return nil
			}

			var f struct {
				FilePath string `json:"file_path"`
				FileSize int64  `json:"file_size"`
			}
			if err := json.Unmarshal(raw, &f); err != nil {
				return fmt.Errorf("parse getFile result: %w", err)
			}
			if f.FilePath == "" {
				return fmt.Errorf("no downloadable path for file_id %q — files over 20MB can't be fetched via the Bot API", fileID)
			}

			target := dest
			if target == "" {
				target = filepath.Base(f.FilePath)
			}
			w := cmd.OutOrStdout()
			if target != "-" {
				out, err := os.Create(target) //nolint:gosec // G304: dest is the user's chosen output path
				if err != nil {
					return err
				}
				defer func() { _ = out.Close() }()
				w = out
			}

			n, err := client.DownloadFile(ctx, f.FilePath, w)
			if err != nil {
				return err
			}
			if target != "-" {
				fmt.Fprintf(cmd.ErrOrStderr(), "downloaded %d bytes to %s\n", n, target)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&fileID, "file-id", "", "file_id to download (from a message's document/photo/etc.)")
	cmd.Flags().StringVar(&dest, "dest", "", "destination path (default: the file's base name; '-' for stdout)")
	_ = cmd.MarkFlagRequired("file-id")
	return cmd
}
