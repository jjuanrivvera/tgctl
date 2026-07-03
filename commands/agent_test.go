package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findCmd walks the tree to the command at the given path (e.g. "auth","login").
func findCmd(root *cobra.Command, path ...string) *cobra.Command {
	cur := root
	for _, name := range path {
		var next *cobra.Command
		for _, c := range cur.Commands() {
			if c.Name() == name {
				next = c
				break
			}
		}
		if next == nil {
			return nil
		}
		cur = next
	}
	return cur
}

func TestMCPExcludesSetupCommands(t *testing.T) {
	sel := ophis.ExcludeCmdsContaining(excludedFromMCP...)
	root := NewRootCmd()

	for _, p := range [][]string{{"auth", "login"}, {"config", "view"}, {"alias", "set"}, {"api"}, {"version"}, {"doctor"}} {
		cmd := findCmd(root, p...)
		require.NotNil(t, cmd, "command %v should exist", p)
		assert.False(t, sel(cmd), "setup/secret command %v must be excluded from the MCP surface", p)
	}
	for _, p := range [][]string{{"message", "send"}, {"bot", "info"}, {"chat", "get"}} {
		cmd := findCmd(root, p...)
		require.NotNil(t, cmd)
		assert.True(t, sel(cmd), "API command %v must be exposed as an MCP tool", p)
	}
}

func TestMCPCommandRegistered(t *testing.T) {
	require.NotNil(t, findCmd(NewRootCmd(), "mcp"), "the mcp subtree must be registered")
}

func TestClassifyAPICommands(t *testing.T) {
	c := classifyAPICommands(false)

	has := func(set []apiCmdInfo, method string) bool {
		for _, x := range set {
			if x.Method == method {
				return true
			}
		}
		return false
	}
	assert.True(t, has(c.Read, "getMe"), "getMe is read-only")
	assert.True(t, has(c.Write, "sendMessage"), "sendMessage is a write")
	assert.True(t, has(c.Destructive, "deleteMessage"), "deleteMessage is destructive")
	assert.True(t, has(c.Destructive, "leaveChat"), "leaveChat is destructive")
	assert.True(t, has(c.Destructive, "banChatMember"), "banChatMember is destructive")
	assert.False(t, has(c.Write, "deleteMessage"), "deleteMessage must not be a mere write")

	// Irreversible operations promoted out of the write bucket: a Stars refund is an
	// unrecoverable money movement and unpinAll* bulk-destroys the pin set.
	assert.True(t, has(c.Destructive, "refundStarPayment"), "refundStarPayment is destructive")
	assert.True(t, has(c.Destructive, "unpinAllChatMessages"), "unpinAllChatMessages is destructive")
	assert.True(t, has(c.Destructive, "unpinAllForumTopicMessages"), "unpinAllForumTopicMessages is destructive")
	assert.True(t, has(c.Destructive, "unpinAllGeneralForumTopicMessages"), "unpinAllGeneralForumTopicMessages is destructive")

	// Invariant: nothing in the read (allowed) bucket may mutate remote state. Every
	// Bot API read is a get* method, so a non-get method in Read is a misclassification.
	for _, r := range c.Read {
		assert.Truef(t, strings.HasPrefix(strings.ToLower(r.Method), "get"),
			"read bucket contains non-get method %s (%s) — a write classified as read", r.Method, r.Path)
	}

	// --all-writes promotes ordinary writes into the hard-block bucket.
	strict := classifyAPICommands(true)
	assert.True(t, has(strict.Destructive, "sendMessage"))
	assert.Empty(t, strict.Write)
}

// claudeCodeSettingsMarker separates the hook-script section from the settings fragment
// in the claude-code guard output.
const claudeCodeSettingsMarker = "# ----- merge into .claude/settings.json -----\n"

func TestAgentGuard_ClaudeCode(t *testing.T) {
	out, _, err := run(t, nil, "agent", "guard", "--host", "claude-code")
	require.NoError(t, err)

	// The output has two sections: the PreToolUse hook script and the settings fragment.
	idx := strings.Index(out, claudeCodeSettingsMarker)
	require.GreaterOrEqual(t, idx, 0, "claude-code output must contain the settings section marker")
	hook := out[:idx]
	settingsJSON := out[idx+len(claudeCodeSettingsMarker):]

	// Hook script: path-prefix-hardened command matching, api gate, MCP branch.
	assert.Contains(t, hook, "#!/usr/bin/env bash")
	assert.Contains(t, hook, `([^[:space:]]*/)?tgctl`, "hook regex must accept a path-prefixed binary")
	assert.Contains(t, hook, "'message delete'")
	assert.Contains(t, hook, "'msg delete-many'", "hook must block alias paths too")
	assert.Contains(t, hook, "'alias set'")
	assert.Contains(t, hook, "'mcp__tgctl__tg_message_delete'")
	assert.Contains(t, hook, "api_is_blocked")

	var settings struct {
		Permissions struct {
			Deny  []string `json:"deny"`
			Ask   []string `json:"ask"`
			Allow []string `json:"allow"`
		} `json:"permissions"`
		Hooks struct {
			PreToolUse []struct {
				Matcher string `json:"matcher"`
			} `json:"PreToolUse"`
		} `json:"hooks"`
	}
	require.NoError(t, json.Unmarshal([]byte(settingsJSON), &settings))
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl message delete:*)")
	assert.Contains(t, settings.Permissions.Deny, "mcp__tgctl__tg_message_delete")
	// Alias paths must be denied too, or `tgctl msg delete` / `tgctl message delete-many`
	// bypass the rules that only name the canonical path.
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl msg delete:*)")
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl message delete-many:*)")
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl cmds delete:*)")
	// Raw api escape hatch: destructive Bot API methods denied by name.
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl api deleteMessage:*)")
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl api banChatMember:*)")
	// Alias minting is denied.
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl alias set:*)")
	assert.Contains(t, settings.Permissions.Ask, "Bash(tgctl message send:*)")
	assert.Contains(t, settings.Permissions.Allow, "Bash(tgctl bot info:*)")
	// A destructive op must never appear in allow — canonical or alias path.
	for _, p := range settings.Permissions.Allow {
		assert.NotContains(t, settings.Permissions.Deny, p, "allow and deny must not overlap: %s", p)
	}
	assert.NotContains(t, settings.Permissions.Allow, "Bash(tgctl chat leave:*)")
	// The hook must be wired for both the Bash and MCP surfaces.
	require.Len(t, settings.Hooks.PreToolUse, 2)
	assert.Equal(t, "Bash", settings.Hooks.PreToolUse[0].Matcher)
	assert.Equal(t, "mcp__tgctl__", settings.Hooks.PreToolUse[1].Matcher)
}

// TestAPICommandAliasPaths pins the alias expansion the guard depends on.
func TestAPICommandAliasPaths(t *testing.T) {
	var deleteBatch *apiCmdInfo
	for i := range registeredAPICmds {
		if registeredAPICmds[i].Path == "message delete-batch" {
			deleteBatch = &registeredAPICmds[i]
			break
		}
	}
	require.NotNil(t, deleteBatch)
	all := deleteBatch.AllPaths()
	assert.Contains(t, all, "message delete-batch")
	assert.Contains(t, all, "message delete-many")
	assert.Contains(t, all, "msg delete-batch")
	assert.Contains(t, all, "msg delete-many")
}

func TestAgentGuard_Codex(t *testing.T) {
	out, _, err := run(t, nil, "agent", "guard", "--host", "codex")
	require.NoError(t, err)
	assert.Contains(t, out, "approval_policy")
	assert.Contains(t, out, "read-only")
	assert.Contains(t, out, "tgctl message delete")
}

func TestAgentGuard_OpenCode(t *testing.T) {
	out, _, err := run(t, nil, "agent", "guard", "--host", "opencode", "--all-writes")
	require.NoError(t, err)
	var cfg struct {
		Permission map[string]string `json:"permission"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &cfg))
	assert.Equal(t, "deny", cfg.Permission["Bash(tgctl message delete:*)"])
	// With --all-writes, sendMessage is hard-denied too.
	assert.Equal(t, "deny", cfg.Permission["Bash(tgctl message send:*)"])
	assert.Equal(t, "allow", cfg.Permission["Bash(tgctl bot info:*)"])
}

func TestAgentGuard_UnknownHost(t *testing.T) {
	_, _, err := run(t, nil, "agent", "guard", "--host", "bogus")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown --host")
}

func TestAgentGuard_WriteToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	_, _, err := run(t, nil, "agent", "guard", "--host", "claude-code", "--out", path)
	require.NoError(t, err)
	data, rerr := os.ReadFile(path)
	require.NoError(t, rerr)
	assert.Contains(t, string(data), "permissions")
}
