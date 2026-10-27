package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusinessReads_WireParams(t *testing.T) {
	t.Run("getBusinessConnection", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "getBusinessConnection", &got)

		out, _, err := run(t, srv, "business", "connection", "--business-connection-id", "bc_123")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
	})

	t.Run("getBusinessAccountGifts with all flags", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "getBusinessAccountGifts", &got)

		out, _, err := run(t, srv, "business", "gifts",
			"--business-connection-id", "bc_123",
			"--exclude-unsaved",
			"--exclude-saved",
			"--exclude-unlimited",
			"--exclude-limited",
			"--exclude-unique",
			"--sort-by-price",
			"--offset", "cursor_abc",
			"--limit", "50",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, true, got["exclude_unsaved"])
		assert.Equal(t, true, got["exclude_saved"])
		assert.Equal(t, true, got["exclude_unlimited"])
		assert.Equal(t, true, got["exclude_limited"])
		assert.Equal(t, true, got["exclude_unique"])
		assert.Equal(t, true, got["sort_by_price"])
		assert.Equal(t, "cursor_abc", got["offset"])
		assert.EqualValues(t, 50, got["limit"])
	})

	t.Run("getBusinessAccountStarBalance", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "getBusinessAccountStarBalance", &got)

		out, _, err := run(t, srv, "business", "star-balance", "--business-connection-id", "bc_123")
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
	})
}

func TestBusinessWrites_WireParams(t *testing.T) {
	t.Run("readBusinessMessage", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "readBusinessMessage", &got)

		out, _, err := run(t, srv, "business", "read-message",
			"--business-connection-id", "bc_123",
			"--chat", "12345678",
			"--message-id", "42",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "12345678", got["chat_id"])
		assert.EqualValues(t, 42, got["message_id"])
	})

	t.Run("deleteBusinessMessages", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "deleteBusinessMessages", &got)

		out, _, err := run(t, srv, "business", "delete-messages",
			"--business-connection-id", "bc_123",
			"--message-ids", "[10,11,12]",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		ids, ok := got["message_ids"].([]any)
		require.True(t, ok)
		assert.Len(t, ids, 3)
		assert.EqualValues(t, 10, ids[0])
		assert.EqualValues(t, 11, ids[1])
		assert.EqualValues(t, 12, ids[2])
	})

	t.Run("setBusinessAccountBio", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountBio", &got)

		out, _, err := run(t, srv, "business", "set-bio",
			"--business-connection-id", "bc_123",
			"--bio", "Official business bio",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "Official business bio", got["bio"])
	})

	t.Run("setBusinessAccountBio empty to remove", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountBio", &got)

		out, _, err := run(t, srv, "business", "set-bio",
			"--business-connection-id", "bc_123",
			"--bio", "",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "", got["bio"])
	})

	t.Run("setBusinessAccountName", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountName", &got)

		out, _, err := run(t, srv, "business", "set-name",
			"--business-connection-id", "bc_123",
			"--first-name", "Acme",
			"--last-name", "Corp",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "Acme", got["first_name"])
		assert.Equal(t, "Corp", got["last_name"])
	})

	t.Run("setBusinessAccountUsername", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountUsername", &got)

		out, _, err := run(t, srv, "business", "set-username",
			"--business-connection-id", "bc_123",
			"--username", "acmeshop",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "acmeshop", got["username"])
	})

	t.Run("setBusinessAccountUsername empty to remove", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountUsername", &got)

		out, _, err := run(t, srv, "business", "set-username",
			"--business-connection-id", "bc_123",
			"--username", "",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, "", got["username"])
	})

	t.Run("setBusinessAccountGiftSettings", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "setBusinessAccountGiftSettings", &got)

		out, _, err := run(t, srv, "business", "set-gift-settings",
			"--business-connection-id", "bc_123",
			"--show-gift-button",
			"--accepted-gift-types", `{"unlimited_gifts":true,"unique_gifts":false}`,
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.Equal(t, true, got["show_gift_button"])
		types, ok := got["accepted_gift_types"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, true, types["unlimited_gifts"])
		assert.Equal(t, false, types["unique_gifts"])
	})

	t.Run("removeBusinessAccountProfilePhoto", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "removeBusinessAccountProfilePhoto", &got)

		out, _, err := run(t, srv, "business", "remove-photo",
			"--business-connection-id", "bc_123",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
	})

	t.Run("transferBusinessAccountStars", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "transferBusinessAccountStars", &got)

		out, _, err := run(t, srv, "business", "transfer-stars",
			"--business-connection-id", "bc_123",
			"--star-count", "250",
		)
		require.NoError(t, err)
		assert.Contains(t, out, "true")
		assert.Equal(t, "bc_123", got["business_connection_id"])
		assert.EqualValues(t, 250, got["star_count"])
	})
}

func TestBusinessValidation(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"connection missing conn id", []string{"business", "connection"}},
		{"gifts limit too low", []string{"business", "gifts", "--business-connection-id", "bc", "--limit", "0"}},
		{"gifts limit too high", []string{"business", "gifts", "--business-connection-id", "bc", "--limit", "101"}},
		{"read-message missing chat", []string{"business", "read-message", "--business-connection-id", "bc", "--message-id", "1"}},
		{"read-message missing message-id", []string{"business", "read-message", "--business-connection-id", "bc", "--chat", "123"}},
		{"delete-messages invalid json", []string{"business", "delete-messages", "--business-connection-id", "bc", "--message-ids", "not-json"}},
		{"delete-messages empty array", []string{"business", "delete-messages", "--business-connection-id", "bc", "--message-ids", "[]"}},
		{"delete-messages not array", []string{"business", "delete-messages", "--business-connection-id", "bc", "--message-ids", `{"id":1}`}},
		{"set-bio too long", []string{"business", "set-bio", "--business-connection-id", "bc", "--bio", strings.Repeat("a", 141)}},
		{"set-name missing first-name", []string{"business", "set-name", "--business-connection-id", "bc"}},
		{"set-name first-name empty", []string{"business", "set-name", "--business-connection-id", "bc", "--first-name", ""}},
		{"set-name first-name too long", []string{"business", "set-name", "--business-connection-id", "bc", "--first-name", strings.Repeat("a", 65)}},
		{"set-name last-name too long", []string{"business", "set-name", "--business-connection-id", "bc", "--first-name", "Acme", "--last-name", strings.Repeat("a", 65)}},
		{"set-username too short", []string{"business", "set-username", "--business-connection-id", "bc", "--username", "abcd"}},
		{"set-username too long", []string{"business", "set-username", "--business-connection-id", "bc", "--username", strings.Repeat("a", 33)}},
		{"set-gift-settings missing types", []string{"business", "set-gift-settings", "--business-connection-id", "bc", "--show-gift-button"}},
		{"set-gift-settings types not object", []string{"business", "set-gift-settings", "--business-connection-id", "bc", "--show-gift-button", "--accepted-gift-types", `[1,2]`}},
		{"transfer-stars zero count", []string{"business", "transfer-stars", "--business-connection-id", "bc", "--star-count", "0"}},
		{"transfer-stars negative count", []string{"business", "transfer-stars", "--business-connection-id", "bc", "--star-count", "-10"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := run(t, nil, tc.args...)
			require.Error(t, err)
		})
	}
}

func TestBusinessKindsAndAnnotations(t *testing.T) {
	root := NewRootCmd()
	for _, tc := range []struct {
		path   string
		method string
		kind   cmdKind
	}{
		{"business connection", "getBusinessConnection", kindRead},
		{"business gifts", "getBusinessAccountGifts", kindRead},
		{"business star-balance", "getBusinessAccountStarBalance", kindRead},
		{"business read-message", "readBusinessMessage", kindWrite},
		{"business delete-messages", "deleteBusinessMessages", kindDestructive},
		{"business set-bio", "setBusinessAccountBio", kindWrite},
		{"business set-name", "setBusinessAccountName", kindWrite},
		{"business set-username", "setBusinessAccountUsername", kindWrite},
		{"business set-gift-settings", "setBusinessAccountGiftSettings", kindWrite},
		{"business remove-photo", "removeBusinessAccountProfilePhoto", kindDestructive},
		{"business transfer-stars", "transferBusinessAccountStars", kindDestructive},
	} {
		t.Run(tc.path, func(t *testing.T) {
			found := false
			for _, c := range APICommands() {
				if c.Path == tc.path {
					found = true
					assert.Equal(t, tc.kind, c.Kind)
					assert.Equal(t, tc.method, c.Method)
				}
			}
			assert.True(t, found, "path must be registered in APICommands")

			cmd, _, err := root.Find(strings.Fields(tc.path))
			require.NoError(t, err)
			switch tc.kind {
			case kindRead:
				assert.Equal(t, "true", cmd.Annotations[annReadOnly])
				assert.Equal(t, "true", cmd.Annotations[annIdempotent])
				assert.Empty(t, cmd.Annotations[annOpenWorld])
				assert.Empty(t, cmd.Annotations[annDestructive])
			case kindWrite:
				assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
				assert.Empty(t, cmd.Annotations[annReadOnly])
				assert.Empty(t, cmd.Annotations[annDestructive])
			case kindDestructive:
				assert.Equal(t, "true", cmd.Annotations[annOpenWorld])
				assert.Equal(t, "true", cmd.Annotations[annDestructive])
				assert.Empty(t, cmd.Annotations[annReadOnly])
			}
		})
	}
}

func TestBusinessTransferStars_NoRetryOn5xx(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/transferBusinessAccountStars") {
			attempts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok":          false,
				"error_code":  500,
				"description": "Internal server error",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	_, _, err := run(t, srv, "business", "transfer-stars",
		"--business-connection-id", "bc_123",
		"--star-count", "100",
	)
	require.Error(t, err)
	// Because transferBusinessAccountStars is destructive (non-idempotent write),
	// it must NOT retry on 5xx to prevent double-spending assets.
	assert.EqualValues(t, 1, attempts.Load(), "destructive transfer must not retry on 5xx")
}
