package api

import (
	"fmt"
	"strings"
)

// Authenticator applies a profile's credential to an outgoing Bot API request. The Bot API
// has exactly one credential scheme (a bot token embedded in the URL path), but we keep the
// single-method interface from the cliwright standard so the design scales down cleanly and
// the redaction logic lives in one place (GOAL.md §1 "auth providers").
type Authenticator interface {
	// RequestURL builds the full URL for a method call, embedding the credential.
	RequestURL(baseURL, method string) string
	// RedactedURL returns the same URL with the secret masked, for --dry-run and logs.
	RedactedURL(baseURL, method string) string
	// FileURL builds the download URL for a file_path returned by getFile.
	FileURL(baseURL, filePath string) string
	// RedactedFileURL returns the file download URL with the secret masked.
	RedactedFileURL(baseURL, filePath string) string
	// Method is the non-secret auth method name recorded in the profile.
	Method() string
}

// BotTokenAuth is the Telegram bot-token authenticator. The token has the shape
// "<bot_id>:<hash>"; the bot_id prefix is not secret (it's the bot's user id), the hash is.
type BotTokenAuth struct {
	Token string
}

// NewBotTokenAuth validates the token shape and returns the authenticator.
func NewBotTokenAuth(token string) (*BotTokenAuth, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("empty bot token")
	}
	id, hash, ok := strings.Cut(token, ":")
	if !ok || id == "" || hash == "" {
		return nil, fmt.Errorf("malformed bot token: want \"<bot_id>:<hash>\" (get one from @BotFather)")
	}
	return &BotTokenAuth{Token: token}, nil
}

func (a *BotTokenAuth) Method() string { return "bot-token" }

func (a *BotTokenAuth) RequestURL(baseURL, method string) string {
	return fmt.Sprintf("%s/bot%s/%s", strings.TrimRight(baseURL, "/"), a.Token, method)
}

func (a *BotTokenAuth) RedactedURL(baseURL, method string) string {
	return fmt.Sprintf("%s/bot%s/%s", strings.TrimRight(baseURL, "/"), RedactToken(a.Token), method)
}

// FileURL builds the URL to download a file. The Bot API serves files from a /file/ prefix
// (https://core.telegram.org/bots/api#getfile), distinct from the method-call path.
func (a *BotTokenAuth) FileURL(baseURL, filePath string) string {
	return fmt.Sprintf("%s/file/bot%s/%s", strings.TrimRight(baseURL, "/"), a.Token, strings.TrimLeft(filePath, "/"))
}

func (a *BotTokenAuth) RedactedFileURL(baseURL, filePath string) string {
	return fmt.Sprintf("%s/file/bot%s/%s", strings.TrimRight(baseURL, "/"), RedactToken(a.Token), strings.TrimLeft(filePath, "/"))
}

// BotID returns the non-secret numeric prefix of the token (the bot's user id).
func (a *BotTokenAuth) BotID() string {
	id, _, _ := strings.Cut(a.Token, ":")
	return id
}

// RedactToken masks a bot token for display: it keeps the non-secret bot_id prefix and
// replaces the secret hash with a fixed marker so it can never leak into a curl line or log.
func RedactToken(token string) string {
	id, _, ok := strings.Cut(token, ":")
	if !ok || id == "" {
		return "***"
	}
	return id + ":<redacted>"
}
