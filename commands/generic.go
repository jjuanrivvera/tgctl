package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/api"
)

// The Telegram Bot API is RPC-method-oriented (sendMessage, getChat, ...), not CRUD on
// resources, so tgctl uses a generic *method-command* builder instead of a generic CRUD
// resource (DECISIONS.md). A group is a noun (message, chat, ...) whose verbs each map 1:1
// to a Bot API method. Adding a method is a few declarative lines in a group file — zero
// edits to this shared builder.

// cmdKind classifies a command for retry safety and MCP/agent-guard annotations.
type cmdKind int

const (
	kindRead        cmdKind = iota // read-only: idempotent, safe to auto-retry, MCP readOnlyHint
	kindWrite                      // creates/changes state: MCP openWorldHint
	kindDestructive                // irreversible (delete/leave/ban): MCP destructiveHint
)

// MCP tool annotation keys (the singular MCP hint keys; ophis reads these from cmd.Annotations).
const (
	annReadOnly    = "readOnlyHint"
	annDestructive = "destructiveHint"
	annOpenWorld   = "openWorldHint"
	annIdempotent  = "idempotentHint"
)

type flagKind int

const (
	flagString flagKind = iota
	flagInt
	flagBool
	flagStringSlice
	flagJSON // value parsed as JSON and sent as-is (objects/arrays, e.g. reply_markup)
)

// flagSpec declares one CLI flag and the Bot API parameter it maps to.
type flagSpec struct {
	Name     string
	Param    string // defaults to Name with '-' → '_'
	Kind     flagKind
	Required bool
	Short    string
	Default  string
	Usage    string
}

func (f flagSpec) param() string {
	if f.Param != "" {
		return f.Param
	}
	return strings.ReplaceAll(f.Name, "-", "_")
}

// fileSpec declares an upload field. The value may be a local file path (sent as
// multipart/form-data), or an http(s) URL / existing Telegram file_id (sent as a string
// param) — matching the Bot API's own flexibility.
type fileSpec struct {
	Name     string
	Param    string
	Required bool
	Usage    string
}

func (f fileSpec) param() string {
	if f.Param != "" {
		return f.Param
	}
	return f.Name
}

// methodCmd declares one verb (a Bot API method).
type methodCmd struct {
	Use     string
	Aliases []string
	Method  string
	Short   string
	Long    string
	Example string
	Kind    cmdKind
	Flags   []flagSpec
	Files   []fileSpec
	Columns []string
}

// group is a noun with its verbs.
type group struct {
	Use     string
	Aliases []string
	Short   string
	Long    string
	Cmds    []methodCmd
}

// apiCmdInfo records a built API command for the MCP server and agent guard to classify.
type apiCmdInfo struct {
	Path   string // e.g. "message send"
	Method string
	Kind   cmdKind
}

var registeredAPICmds []apiCmdInfo

// APICommands returns the classification of every API-backed command (for agent guard tests).
func APICommands() []apiCmdInfo { return registeredAPICmds }

// Path returns the command path, Method the Bot API method, and Kind the classification.
func (a apiCmdInfo) PathString() string  { return a.Path }
func (a apiCmdInfo) IsRead() bool        { return a.Kind == kindRead }
func (a apiCmdInfo) IsDestructive() bool { return a.Kind == kindDestructive }

// registerGroup adds a group's commands to the root tree and the classification registry.
func registerGroup(g group) {
	for _, mc := range g.Cmds {
		registeredAPICmds = append(registeredAPICmds, apiCmdInfo{
			Path:   g.Use + " " + mc.Use,
			Method: mc.Method,
			Kind:   mc.Kind,
		})
	}
	register(func(root *cobra.Command) {
		parent := &cobra.Command{
			Use:     g.Use,
			Aliases: g.Aliases,
			Short:   g.Short,
			Long:    g.Long,
		}
		for _, mc := range g.Cmds {
			parent.AddCommand(buildMethodCmd(mc))
		}
		root.AddCommand(parent)
	})
}

// buildMethodCmd turns a methodCmd into a cobra command: it binds the declared flags, stamps
// MCP annotations from the Kind, and on run assembles the params, calls the API (Upload when
// there are file fields), and renders the result.
func buildMethodCmd(mc methodCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:     mc.Use,
		Aliases: mc.Aliases,
		Short:   mc.Short,
		Long:    mc.Long,
		Example: mc.Example,
		Args:    cobra.NoArgs,
	}
	markKind(cmd, mc.Kind)
	bindFlags(cmd, mc)

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		params, err := collectParams(cmd, mc)
		if err != nil {
			return err
		}
		files, err := collectFiles(cmd, mc, params)
		if err != nil {
			return err
		}
		client, err := clientFromCmd(cmd)
		if err != nil {
			return err
		}
		idempotent := mc.Kind == kindRead
		var raw json.RawMessage
		if len(files) > 0 {
			raw, err = client.Upload(cmd.Context(), mc.Method, params, files, idempotent)
		} else {
			raw, err = client.Call(cmd.Context(), mc.Method, params, idempotent)
		}
		if err != nil {
			return err
		}
		if len(mc.Columns) > 0 && !cmd.Flags().Changed("columns") {
			// Apply the command's default columns unless the user overrode --columns.
			if err := cmd.Flags().Set("columns", strings.Join(mc.Columns, ",")); err != nil {
				return err
			}
		}
		return render(cmd, raw)
	}
	return cmd
}

// markKind stamps MCP annotations so the mcp server and agent guard can gate writes. A write
// sets only openWorldHint; a destructive verb adds destructiveHint; a read sets readOnlyHint.
func markKind(cmd *cobra.Command, kind cmdKind) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	switch kind {
	case kindRead:
		cmd.Annotations[annReadOnly] = "true"
		cmd.Annotations[annIdempotent] = "true"
	case kindWrite:
		cmd.Annotations[annOpenWorld] = "true"
	case kindDestructive:
		cmd.Annotations[annOpenWorld] = "true"
		cmd.Annotations[annDestructive] = "true"
	}
}

func bindFlags(cmd *cobra.Command, mc methodCmd) {
	f := cmd.Flags()
	for _, fs := range mc.Flags {
		switch fs.Kind {
		case flagString, flagJSON:
			f.StringP(fs.Name, fs.Short, fs.Default, fs.Usage)
		case flagInt:
			f.Int64P(fs.Name, fs.Short, 0, fs.Usage)
		case flagBool:
			f.BoolP(fs.Name, fs.Short, false, fs.Usage)
		case flagStringSlice:
			f.StringSliceP(fs.Name, fs.Short, nil, fs.Usage)
		}
		if fs.Required {
			_ = cmd.MarkFlagRequired(fs.Name)
		}
	}
	for _, fl := range mc.Files {
		cmd.Flags().String(fl.Name, "", fl.Usage)
		if fl.Required {
			_ = cmd.MarkFlagRequired(fl.Name)
		}
	}
}

// collectParams reads the set flags into the Bot API params map. Only flags the user actually
// set are sent, so API defaults apply otherwise.
func collectParams(cmd *cobra.Command, mc methodCmd) (map[string]any, error) {
	params := map[string]any{}
	f := cmd.Flags()
	for _, fs := range mc.Flags {
		if !f.Changed(fs.Name) {
			continue
		}
		switch fs.Kind {
		case flagString:
			v, _ := f.GetString(fs.Name)
			params[fs.param()] = v
		case flagInt:
			v, _ := f.GetInt64(fs.Name)
			params[fs.param()] = v
		case flagBool:
			v, _ := f.GetBool(fs.Name)
			params[fs.param()] = v
		case flagStringSlice:
			v, _ := f.GetStringSlice(fs.Name)
			params[fs.param()] = v
		case flagJSON:
			v, _ := f.GetString(fs.Name)
			var parsed any
			if err := json.Unmarshal([]byte(v), &parsed); err != nil {
				return nil, fmt.Errorf("--%s must be valid JSON: %w", fs.Name, err)
			}
			params[fs.param()] = parsed
		}
	}
	return params, nil
}

// collectFiles resolves file fields. A value that is an existing local file becomes a
// multipart upload (path-validated); an http(s) URL or a file_id is sent as a string param.
func collectFiles(cmd *cobra.Command, mc methodCmd, params map[string]any) (map[string]string, error) {
	files := map[string]string{}
	f := cmd.Flags()
	for _, fl := range mc.Files {
		if !f.Changed(fl.Name) {
			continue
		}
		v, _ := f.GetString(fl.Name)
		if v == "" {
			continue
		}
		if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
			params[fl.param()] = v // a URL: Telegram fetches it directly
			continue
		}
		if err := api.ValidateUploadPath(v); err != nil {
			// Not a readable local file → treat as an existing Telegram file_id.
			params[fl.param()] = v
			continue
		}
		files[fl.param()] = v
	}
	return files, nil
}
