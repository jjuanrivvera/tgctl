package api

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// The Bot API carries its credential in the URL path, so any component that prints a URL
// prints the secret. Go's *url.Error is exactly such a component: http.Client.Do wraps every
// transport failure in one, and its Error() renders the URL verbatim. That is how a read
// timeout put a live token on stderr and into an agent transcript (issue #21). Redaction
// therefore has to happen on the error path, not only on the --dry-run path.

// tokenInURLRe matches a token where the Bot API puts it: after "/bot" for method calls and
// after "/file/bot" for downloads. The numeric bot id is kept — it is the bot's user id, not
// a secret — and only the hash is masked, so a redacted error still says which bot failed.
var tokenInURLRe = regexp.MustCompile(`(/(?:file/)?bot)(-?\d+):[A-Za-z0-9_-]+`)

// bareTokenRe matches a token that appears outside a URL (an env var echoed into a message,
// a config dump). It demands both halves of the "<bot_id>:<hash>" shape with a hash long
// enough that ordinary "id:value" text does not match.
var bareTokenRe = regexp.MustCompile(`\b(\d{5,}):[A-Za-z0-9_-]{20,}\b`)

// RedactSecrets masks any Telegram bot token found anywhere in s. It is pattern-based, so it
// works without knowing which token is active — that makes it usable as a last line of
// defense at an output boundary (see cmd/tgctl), where the credential is out of scope.
func RedactSecrets(s string) string {
	if s == "" {
		return s
	}
	s = tokenInURLRe.ReplaceAllString(s, "${1}${2}:<redacted>")
	return bareTokenRe.ReplaceAllString(s, "${1}:<redacted>")
}

// RedactSecrets masks this authenticator's own token in s, then applies the generic patterns.
// Knowing the exact token catches shapes the patterns would miss — a self-hosted Local Bot API
// Server's credential, or a hash using characters outside the usual alphabet.
func (a *BotTokenAuth) RedactSecrets(s string) string {
	if a.Token != "" {
		s = strings.ReplaceAll(s, a.Token, RedactToken(a.Token))
	}
	return RedactSecrets(s)
}

// secretRedactor is the optional Authenticator capability of masking its own credential.
// It stays out of the Authenticator interface so a custom authenticator is not forced to
// implement it: a type that does not gets the pattern-only redaction.
type secretRedactor interface {
	RedactSecrets(s string) string
}

// redact masks the active credential in an arbitrary string.
func (c *Client) redact(s string) string {
	if r, ok := c.auth.(secretRedactor); ok {
		return r.RedactSecrets(s)
	}
	return RedactSecrets(s)
}

// redactError returns err with the bot token masked everywhere it would print. Every error
// leaving the network path goes through it.
func (c *Client) redactError(err error) error {
	if err == nil {
		return nil
	}
	// The structured leak: rewrite *url.Error's URL field in place rather than re-wrapping,
	// so the whole error chain survives (errors.Is on the cause, net.Error assertions) and a
	// later errors.As inspection finds an already-clean URL instead of the raw one.
	var ue *url.Error
	if errors.As(err, &ue) {
		ue.URL = c.redact(ue.URL)
	}
	// Anything else that rendered the token into its own message.
	msg := err.Error()
	if red := c.redact(msg); red != msg {
		return &redactedError{msg: red, err: err}
	}
	return err
}

// redactedError substitutes a redacted message for an error's own while keeping the original
// reachable through Unwrap, so errors.Is/As still match the cause.
type redactedError struct {
	msg string
	err error
}

func (e *redactedError) Error() string { return e.msg }
func (e *redactedError) Unwrap() error { return e.err }
