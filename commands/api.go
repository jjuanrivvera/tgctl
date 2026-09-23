package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/api"
)

func init() {
	register(func(root *cobra.Command) {
		var data string
		var query []string
		var uploads []string
		var idempotent bool
		cmd := &cobra.Command{
			Use:   "api <method> [-d body] [-q key=value ...] [-F name=@file ...]",
			Short: "Call any Bot API method directly (raw escape hatch)",
			Long: `Invoke an arbitrary Bot API method with a JSON body and/or key=value parameters.

This is the documented escape hatch for methods tgctl does not wrap as first-class
commands. It honors --dry-run and -o/--output like every other command. By default a
raw call is treated as a write (not auto-retried); pass --idempotent for read-only
methods (getX) so transient failures retry safely.

Uploading files: some methods refuse anything but an uploaded file (setMyProfilePhoto,
setChatPhoto, uploadStickerFile, ...). Pass -F name=@path — repeatable — and the request
switches to multipart/form-data, sending each file as the part <name> while -d/-q travel
alongside as fields. That is what lets a JSON body reference an upload as "attach://<name>":

  tgctl api setMyProfilePhoto \
    -d '{"photo":{"type":"static","photo":"attach://pic"}}' -F pic=@logo.jpg`,
			Example: `  tgctl api getMe
  tgctl api sendMessage -q chat_id=@me -q text="hi from tgctl"
  tgctl api getChat -q chat_id=@telegram --idempotent
  tgctl api sendMessage -d '{"chat_id":"@me","text":"json body"}'
  tgctl api setChatPhoto -q chat_id=@mygroup -F photo=@cover.jpg
  tgctl api setMyProfilePhoto -d '{"photo":{"type":"static","photo":"attach://pic"}}' -F pic=@logo.jpg`,
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
				files, err := parseUploads(uploads)
				if err != nil {
					return err
				}
				client, err := clientFromCmd(cmd)
				if err != nil {
					return err
				}
				defer func() { _ = client.Close() }()
				var raw json.RawMessage
				if len(files) > 0 {
					raw, err = client.Upload(cmd.Context(), method, params, files, idempotent)
				} else {
					raw, err = client.Call(cmd.Context(), method, params, idempotent)
				}
				if err != nil {
					return err
				}
				return render(cmd, raw)
			},
		}
		cmd.Flags().StringVarP(&data, "data", "d", "", "raw JSON request body")
		cmd.Flags().StringArrayVarP(&query, "query", "q", nil, "key=value parameter (repeatable)")
		cmd.Flags().StringArrayVarP(&uploads, "file", "F", nil, "name=@path to upload as multipart (repeatable)")
		cmd.Flags().BoolVar(&idempotent, "idempotent", false, "treat as read-only (safe to auto-retry)")
		root.AddCommand(cmd)
	})
}

// parseUploads turns the -F name=@path flags into the field→path map Upload wants. The @ is
// required, as in curl: it is what separates "upload this file" from "send this literal
// value", and a literal value already has a flag of its own (-q).
func parseUploads(specs []string) (map[string]string, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	files := make(map[string]string, len(specs))
	for _, spec := range specs {
		name, value, ok := strings.Cut(spec, "=")
		if !ok || name == "" {
			return nil, fmt.Errorf("invalid -F %q (want name=@path)", spec)
		}
		path, ok := strings.CutPrefix(value, "@")
		if !ok {
			return nil, fmt.Errorf("invalid -F %q: the value must be @<path> (for a literal value use -q %s=%s)", spec, name, value)
		}
		if path == "" {
			return nil, fmt.Errorf("invalid -F %q: empty path after @", spec)
		}
		if err := api.ValidateUploadPath(path); err != nil {
			return nil, fmt.Errorf("-F %s: %w", name, err)
		}
		if _, dup := files[name]; dup {
			return nil, fmt.Errorf("duplicate -F part %q", name)
		}
		files[name] = path
	}
	return files, nil
}
