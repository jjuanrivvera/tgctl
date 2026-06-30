package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBotTokenAuth_Validation(t *testing.T) {
	_, err := NewBotTokenAuth("")
	require.Error(t, err)

	_, err = NewBotTokenAuth("no-colon-here")
	require.Error(t, err)

	_, err = NewBotTokenAuth(":onlyhash")
	require.Error(t, err)

	a, err := NewBotTokenAuth("  123456:GOODHASH123  ")
	require.NoError(t, err)
	assert.Equal(t, "123456", a.BotID())
	assert.Equal(t, "bot-token", a.Method())
}

func TestBotTokenAuth_URLs(t *testing.T) {
	a, _ := NewBotTokenAuth("123456:SECRETHASH")
	url := a.RequestURL("https://api.telegram.org/", "sendMessage")
	assert.Equal(t, "https://api.telegram.org/bot123456:SECRETHASH/sendMessage", url)

	red := a.RedactedURL("https://api.telegram.org", "sendMessage")
	assert.Equal(t, "https://api.telegram.org/bot123456:<redacted>/sendMessage", red)
	assert.False(t, strings.Contains(red, "SECRETHASH"))
}

func TestBotTokenAuth_FileURLs(t *testing.T) {
	a, _ := NewBotTokenAuth("123456:SECRETHASH")
	url := a.FileURL("https://api.telegram.org/", "photos/file_1.jpg")
	assert.Equal(t, "https://api.telegram.org/file/bot123456:SECRETHASH/photos/file_1.jpg", url)

	red := a.RedactedFileURL("https://api.telegram.org", "/photos/file_1.jpg")
	assert.Equal(t, "https://api.telegram.org/file/bot123456:<redacted>/photos/file_1.jpg", red)
	assert.False(t, strings.Contains(red, "SECRETHASH"))
}

func TestRedactToken(t *testing.T) {
	assert.Equal(t, "123:<redacted>", RedactToken("123:abcdef"))
	assert.Equal(t, "***", RedactToken("garbage"))
}
