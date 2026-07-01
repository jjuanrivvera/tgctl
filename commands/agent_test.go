package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
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

	// --all-writes promotes ordinary writes into the hard-block bucket.
	strict := classifyAPICommands(true)
	assert.True(t, has(strict.Destructive, "sendMessage"))
	assert.Empty(t, strict.Write)
}

func TestAgentGuard_ClaudeCode(t *testing.T) {
	out, _, err := run(t, nil, "agent", "guard", "--host", "claude-code")
	require.NoError(t, err)

	var settings struct {
		Permissions struct {
			Deny  []string `json:"deny"`
			Ask   []string `json:"ask"`
			Allow []string `json:"allow"`
		} `json:"permissions"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &settings))
	assert.Contains(t, settings.Permissions.Deny, "Bash(tgctl message delete:*)")
	assert.Contains(t, settings.Permissions.Deny, "mcp__tgctl__tg_message_delete")
	assert.Contains(t, settings.Permissions.Ask, "Bash(tgctl message send:*)")
	assert.Contains(t, settings.Permissions.Allow, "Bash(tgctl bot info:*)")
	// A destructive op must never appear in allow.
	assert.NotContains(t, settings.Permissions.Allow, "Bash(tgctl chat leave:*)")
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
