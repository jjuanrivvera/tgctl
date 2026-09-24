package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStarsBalance(t *testing.T) {
	srv := newServer(t, routes{"getMyStarBalance": `{"amount":1500,"nanostar_amount":0}`})
	out, _, err := run(t, srv, "stars", "balance", "-o", "json")
	require.NoError(t, err)
	assert.Contains(t, out, "1500")
}

func TestStarsUserGifts_Filters(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "getUserGifts", &got)

	_, _, err := run(t, srv, "stars", "user-gifts", "--user", "12345",
		"--exclude-unlimited", "--sort-by-price", "--limit", "10")
	require.NoError(t, err)
	assert.EqualValues(t, 12345, got["user_id"])
	assert.Equal(t, true, got["exclude_unlimited"])
	assert.Equal(t, true, got["sort_by_price"])
	assert.EqualValues(t, 10, got["limit"])
	assert.NotContains(t, got, "exclude_unique", "filters the user did not ask for must not be sent")
}

// The two listings share the --exclude-* vocabulary, and chat-gifts adds the saved/unsaved
// pair that only makes sense for a channel's profile page.
func TestStarsChatGifts(t *testing.T) {
	var got map[string]any
	srv := paramProbe(t, "getChatGifts", &got)

	_, _, err := run(t, srv, "stars", "chat-gifts", "--chat", "@mychannel", "--exclude-unsaved")
	require.NoError(t, err)
	assert.Equal(t, "@mychannel", got["chat_id"])
	assert.Equal(t, true, got["exclude_unsaved"])
}

func TestStarsGiftOperations(t *testing.T) {
	t.Run("convert", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "convertGiftToStars", &got)
		_, _, err := run(t, srv, "stars", "convert-gift", "--business-connection-id", "BQ1", "--gift-id", "g1")
		require.NoError(t, err)
		assert.Equal(t, "BQ1", got["business_connection_id"])
		assert.Equal(t, "g1", got["owned_gift_id"])
	})

	t.Run("upgrade", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "upgradeGift", &got)
		_, _, err := run(t, srv, "stars", "upgrade-gift", "--business-connection-id", "BQ1",
			"--gift-id", "g1", "--keep-original-details", "--star-count", "25")
		require.NoError(t, err)
		assert.Equal(t, true, got["keep_original_details"])
		assert.EqualValues(t, 25, got["star_count"])
	})

	t.Run("transfer", func(t *testing.T) {
		var got map[string]any
		srv := paramProbe(t, "transferGift", &got)
		_, _, err := run(t, srv, "stars", "transfer-gift", "--business-connection-id", "BQ1",
			"--gift-id", "g1", "--new-owner", "999")
		require.NoError(t, err)
		assert.EqualValues(t, 999, got["new_owner_chat_id"])
	})
}

// Every gift operation moves an asset one way, so the agent guard and the MCP annotations
// must see them as destructive — not as ordinary writes.
func TestStarsGiftOperations_AreDestructive(t *testing.T) {
	oneWay := map[string]bool{
		"convertGiftToStars":      false,
		"upgradeGift":             false,
		"transferGift":            false,
		"giftPremiumSubscription": false,
	}
	for _, c := range APICommands() {
		if _, tracked := oneWay[c.Method]; tracked {
			oneWay[c.Method] = c.IsDestructive()
		}
	}
	for method, destructive := range oneWay {
		assert.True(t, destructive, "%s is one-way and should be classified destructive", method)
	}
}

// The API sells exactly three durations at exactly three prices, and this command spends real
// Stars: a mistyped amount is better met before the charge than after.
func TestStarsGiftPremium_Pricing(t *testing.T) {
	srv := newServer(t, routes{"giftPremiumSubscription": `true`})

	for _, ok := range [][2]string{{"3", "1000"}, {"6", "1500"}, {"12", "2500"}} {
		t.Run("valid "+ok[0]+"m", func(t *testing.T) {
			_, _, err := run(t, srv, "stars", "gift-premium", "--user", "1", "--months", ok[0], "--stars", ok[1])
			require.NoError(t, err)
		})
	}

	_, _, err := run(t, srv, "stars", "gift-premium", "--user", "1", "--months", "3", "--stars", "2500")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "costs 1000 stars")

	_, _, err = run(t, srv, "stars", "gift-premium", "--user", "1", "--months", "4", "--stars", "1000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be 3, 6 or 12")
}
