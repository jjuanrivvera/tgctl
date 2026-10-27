package commands

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManifestMatchesRegistryAndOfficialIndex(t *testing.T) {
	var manifest struct {
		Total     int `json:"api_method_total"`
		Resources []struct {
			Name  string   `json:"name"`
			Verbs []string `json:"verbs"`
		} `json:"resources"`
	}
	data, err := os.ReadFile("../api-manifest.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &manifest))
	var snapshot struct {
		Source  string   `json:"source"`
		Version string   `json:"version"`
		Methods []string `json:"methods"`
	}
	data, err = os.ReadFile("testdata/bot_api_10_3_methods.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &snapshot))
	require.NotEmpty(t, snapshot.Source)
	require.NotEmpty(t, snapshot.Version)
	require.Len(t, snapshot.Methods, manifest.Total)
	require.True(t, sort.StringsAreSorted(snapshot.Methods))
	official := map[string]bool{}
	for _, method := range snapshot.Methods {
		require.False(t, official[method], method)
		official[method] = true
	}
	paths := map[string]bool{}
	root := NewRootCmd()
	for _, resource := range manifest.Resources {
		for _, verb := range resource.Verbs {
			path := resource.Name + " " + verb
			require.False(t, paths[path], path)
			paths[path] = true
			cmd, remaining, err := root.Find(strings.Fields(path))
			require.NoError(t, err)
			require.Empty(t, remaining, path)
			require.Equal(t, "tgctl "+path, cmd.CommandPath())
		}
	}
	methods := map[string]bool{}
	for _, info := range APICommands() {
		require.True(t, paths[info.Path], info.Path)
		delete(paths, info.Path)
		require.False(t, methods[info.Method], info.Method)
		methods[info.Method] = true
		require.True(t, official[info.Method], info.Method)
	}
	require.Empty(t, paths)
}
