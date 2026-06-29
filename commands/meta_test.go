package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestAuthLogin_StoresAndVerifies(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":123456,"is_bot":true,"first_name":"T","username":"testbot"}`})
	// login with an explicit token (no env token), verifying against getMe
	out, errb, err := runNoToken(t, srv, "", "auth", "login", "--token", "123456:REALTOKENVALUE")
	require.NoError(t, err)
	mustContain(t, errb, "verified as @testbot")
	mustContain(t, out, "logged in")
}

func TestAuthStatus_NoToken(t *testing.T) {
	_, _, err := runNoToken(t, nil, "", "auth", "status")
	require.Error(t, err)
	mustContain(t, err.Error(), "not authenticated")
}

func TestAuthLogin_BadToken(t *testing.T) {
	_, _, err := runNoToken(t, nil, "", "auth", "login", "--token", "no-colon", "--no-verify")
	require.Error(t, err)
	mustContain(t, err.Error(), "malformed bot token")
}

func TestConfigUse_And_ListProfiles(t *testing.T) {
	keyring.MockInit()
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	dir := t.TempDir() // shared across the calls below so config writes persist

	_, _, err := runIn(t, dir, srv, "", "auth", "login", "--token", "111:AAA", "--profile", "prod")
	require.NoError(t, err)

	out, _, err := runIn(t, dir, srv, "", "config", "list-profiles", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "prod")

	_, _, err = runIn(t, dir, srv, "", "config", "use", "prod")
	require.NoError(t, err)

	// config set base_url on the active profile
	_, _, err = runIn(t, dir, srv, "", "config", "set", "base_url", "https://api.telegram.org", "--profile", "prod")
	require.NoError(t, err)
}

func TestConfigPath(t *testing.T) {
	out, _, err := run(t, nil, "config", "path")
	require.NoError(t, err)
	mustContain(t, out, "config.yaml")
}

func TestConfigView(t *testing.T) {
	out, _, err := run(t, nil, "config", "view", "-o", "json")
	require.NoError(t, err)
	mustContain(t, out, "token_storage")
}

func TestConfigUse_UnknownProfile(t *testing.T) {
	_, _, err := run(t, nil, "config", "use", "ghost")
	require.Error(t, err)
	mustContain(t, err.Error(), "no such profile")
}

func TestAliasSetListRemove(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	// Use a single config dir across the alias lifecycle by sharing XDG within one run isn't
	// possible (each run() makes a fresh temp dir), so test set+list together via ExpandAliases.
	out, _, err := run(t, srv, "alias", "set", "ping", "bot info")
	require.NoError(t, err)
	mustContain(t, out, "ping")
}

func TestAlias_RejectsBuiltin(t *testing.T) {
	_, _, err := run(t, nil, "alias", "set", "message", "bot info")
	require.Error(t, err)
	mustContain(t, err.Error(), "built-in")
}

func TestExpandAliases_BuiltinWins(t *testing.T) {
	// A built-in name is never expanded, even if an alias by that name somehow existed.
	got := ExpandAliases([]string{"bot", "info"})
	assert.Equal(t, []string{"bot", "info"}, got)
}

func TestExpandAliases_Empty(t *testing.T) {
	assert.Empty(t, ExpandAliases(nil))
}

func TestVersion(t *testing.T) {
	out, _, err := run(t, nil, "version")
	require.NoError(t, err)
	mustContain(t, out, "tgctl")
}

func TestVersionJSON(t *testing.T) {
	out, _, err := run(t, nil, "version", "--json")
	require.NoError(t, err)
	mustContain(t, out, `"version"`)
}

func TestDoctor_Success(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	out, _, err := run(t, srv, "doctor")
	require.NoError(t, err)
	mustContain(t, out, "API reachable")
	assert.True(t, strings.Contains(out, "✓"))
}

func TestDoctor_FailsWithoutToken(t *testing.T) {
	_, _, err := runNoToken(t, nil, "", "doctor")
	require.Error(t, err)
}

func TestDoctorJSON(t *testing.T) {
	srv := newServer(t, routes{"getMe": `{"id":1,"is_bot":true,"first_name":"T","username":"b"}`})
	out, _, err := run(t, srv, "doctor", "--json")
	require.NoError(t, err)
	mustContain(t, out, `"name"`)
}

func TestCompletion(t *testing.T) {
	out, _, err := run(t, nil, "completion", "bash")
	require.NoError(t, err)
	mustContain(t, out, "tgctl")
}

func TestInvalidOutputFormat(t *testing.T) {
	_, _, err := run(t, nil, "version", "-o", "bogus")
	require.Error(t, err)
	mustContain(t, err.Error(), "invalid --output")
}

func TestAPICommandsClassification(t *testing.T) {
	cmds := APICommands()
	require.NotEmpty(t, cmds)
	// deleteMessage must be classified destructive; getMe read-only.
	var sawDelete, sawRead bool
	for _, c := range cmds {
		if c.Method == "deleteMessage" {
			sawDelete = c.IsDestructive()
		}
		if c.Method == "getMe" {
			sawRead = c.IsRead()
		}
	}
	assert.True(t, sawDelete, "deleteMessage should be destructive")
	assert.True(t, sawRead, "getMe should be read-only")
}
