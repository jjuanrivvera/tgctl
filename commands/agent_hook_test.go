package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestHookScript_BashExecution exercises the generated hook script with real bash to verify
// the adversarial cases: obfuscation, path-invoked binaries, alias paths, the raw api
// escape hatch, and the benign lookalikes that must stay allowed. Gated on a POSIX shell
// being available so it is safe in the regular suite.
func TestHookScript_BashExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	// Generate the hook from the real classification so blocked_cmds/blocked_tools are
	// fully populated (canonical + alias paths).
	hookContent := hookScript(classifyAPICommands(false))
	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "tgctl-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookContent), 0o755); err != nil { // #nosec G306 -- hook must be executable
		t.Fatalf("write hook: %v", err)
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  "Bash",
			"tool_input": map[string]any{"command": command},
		})
		return string(b)
	}
	mcpPayload := func(toolName string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  toolName,
			"tool_input": map[string]any{},
		})
		return string(b)
	}

	runHook := func(t *testing.T, payload string, extraEnv []string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		if extraEnv != nil {
			cmd.Env = append(os.Environ(), extraEnv...)
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		// The hook always exits 0; the decision is in the JSON output.
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}

	isDenied := func(output string) bool {
		return strings.Contains(output, `"permissionDecision":"deny"`)
	}

	cases := []struct {
		name       string
		payload    string
		wantDenied bool
	}{
		// --- direct blocked commands ---
		{"message_delete_denied", bashPayload("tgctl message delete --chat @c --id 5"), true},
		{"chat_leave_denied", bashPayload("tgctl chat leave --chat @group"), true},
		{"webhook_delete_denied", bashPayload("tgctl webhook delete"), true},
		// --- cobra alias paths (bypass before the fix) ---
		{"msg_alias_delete_denied", bashPayload("tgctl msg delete --chat @c --id 5"), true},
		{"delete_many_alias_denied", bashPayload("tgctl message delete-many --chat @c --ids 1,2"), true},
		{"cmds_alias_delete_denied", bashPayload("tgctl cmds delete"), true},
		// --- alias minting ---
		{"alias_set_denied", bashPayload(`tgctl alias set kill "message delete"`), true},
		// --- obfuscation ---
		{"quote_split_denied", bashPayload(`tgctl message de""lete --chat @c --id 5`), true},
		{"single_quote_split_denied", bashPayload(`tgctl chat le''ave --chat @g`), true},
		{"backslash_denied", bashPayload(`tgctl message de\lete --chat @c --id 5`), true},
		{"newline_continuation_denied", bashPayload("tgctl message \\\ndelete --chat @c --id 5"), true},
		// --- command position after separators ---
		{"after_semicolon_denied", bashPayload("true; tgctl message delete --chat @c --id 5"), true},
		{"after_pipe_denied", bashPayload("echo hi | tgctl message delete --chat @c --id 5"), true},
		{"after_and_denied", bashPayload("true && tgctl chat leave --chat @g"), true},
		{"trailing_separator_denied", bashPayload("tgctl webhook delete;true"), true},
		{"env_prefix_denied", bashPayload("env TGCTL_TOKEN=x tgctl chat leave --chat @g"), true},
		// --- path-invoked binaries (bug class 2) ---
		{"relative_path_binary_denied", bashPayload("./bin/tgctl message delete --chat @c --id 5"), true},
		{"absolute_path_binary_denied", bashPayload("/usr/local/bin/tgctl message delete --chat @c --id 5"), true},
		{"absolute_path_api_denied", bashPayload("/usr/local/bin/tgctl api deleteWebhook"), true},
		// --- raw api escape hatch (method position; names are case-insensitive) ---
		{"api_delete_method_denied", bashPayload("tgctl api deleteMessage -q chat_id=1 -q message_id=2"), true},
		{"api_uppercase_method_denied", bashPayload("tgctl api DELETEMESSAGE -q chat_id=1"), true},
		{"api_send_method_denied", bashPayload("tgctl api sendMessage -q chat_id=1 -q text=hi"), true},
		{"api_flag_before_method_denied", bashPayload("tgctl api -q chat_id=1 deleteWebhook"), true},
		{"api_compound_get_then_delete_denied", bashPayload("tgctl api getMe;tgctl api deleteWebhook"), true},
		// --- benign lookalikes that must stay allowed ---
		{"bot_info_allowed", bashPayload("tgctl bot info"), false},
		{"send_with_delete_in_arg_allowed", bashPayload(`tgctl message send --chat @c --text "how to delete a message"`), false},
		{"cat_file_allowed", bashPayload("cat message_delete.go"), false},
		{"api_get_allowed", bashPayload("tgctl api getMe"), false},
		{"api_get_with_delete_in_param_allowed", bashPayload("tgctl api getChat -q chat_id=@delete_club"), false},
		{"other_binary_allowed", bashPayload("mytgctl message delete --chat @c --id 5"), false},
		{"other_binary_api_allowed", bashPayload("mytgctl api deleteWebhook"), false},
		// --- MCP branch ---
		{"mcp_message_delete_denied", mcpPayload("mcp__tgctl__tg_message_delete"), true},
		{"mcp_unpin_all_denied", mcpPayload("mcp__tgctl__tg_chat_unpin-all"), true},
		{"mcp_bot_info_allowed", mcpPayload("mcp__tgctl__tg_bot_info"), false},
		{"mcp_near_miss_allowed", mcpPayload("mcp__tgctl__tg_message_delete_x"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := runHook(t, tc.payload, nil)
			if denied := isDenied(output); denied != tc.wantDenied {
				t.Errorf("want denied=%v, got denied=%v\noutput: %s", tc.wantDenied, denied, output)
			}
		})
	}
}

// TestHookScript_BashExecutionNoJq exercises the no-jq fallback path by shadowing jq with
// an empty PATH entry while keeping the POSIX utilities the hook needs.
func TestHookScript_BashExecutionNoJq(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	hookContent := hookScript(classifyAPICommands(false))
	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "tgctl-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookContent), 0o755); err != nil { // #nosec G306 -- hook must be executable
		t.Fatalf("write hook: %v", err)
	}

	// Build a PATH without jq: a bin dir holding symlinks to everything the hook uses.
	binDir := filepath.Join(tmpDir, "nojq-bin")
	if err := os.Mkdir(binDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for _, tool := range []string{"cat", "tr", "grep", "sed", "printf", "env"} {
		p, lerr := exec.LookPath(tool)
		if lerr != nil {
			continue // shell builtins (printf) need no symlink
		}
		if serr := os.Symlink(p, filepath.Join(binDir, tool)); serr != nil {
			t.Fatalf("symlink %s: %v", tool, serr)
		}
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  "Bash",
			"tool_input": map[string]any{"command": command},
		})
		return string(b)
	}

	runHookNoJq := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		env := make([]string, 0, len(os.Environ()))
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "PATH=") {
				env = append(env, e)
			}
		}
		cmd.Env = append(env, "PATH="+binDir)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}

	isDenied := func(output string) bool {
		return strings.Contains(output, `"permissionDecision":"deny"`)
	}

	cases := []struct {
		name       string
		payload    string
		wantDenied bool
	}{
		{"nojq_message_delete_denied", bashPayload("tgctl message delete --chat @c --id 5"), true},
		{"nojq_obfuscated_delete_denied", bashPayload(`tgctl message de""lete --chat @c --id 5`), true},
		{"nojq_path_binary_denied", bashPayload("./bin/tgctl message delete --chat @c --id 5"), true},
		{"nojq_api_delete_denied", bashPayload("tgctl api deleteWebhook"), true},
		{"nojq_cat_file_allowed", bashPayload("cat message_delete.go"), false},
		{"nojq_send_allowed", bashPayload(`tgctl message send --chat @c --text "delete this later"`), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := runHookNoJq(t, tc.payload)
			if denied := isDenied(output); denied != tc.wantDenied {
				t.Errorf("want denied=%v, got denied=%v\noutput: %s", tc.wantDenied, denied, output)
			}
		})
	}
}
