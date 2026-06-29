package api

import (
	"fmt"
	"strings"
)

// APIError is a typed Telegram Bot API failure. The Bot API returns
// {"ok":false,"error_code":N,"description":"...","parameters":{...}}; we surface all of it
// and, crucially, append an actionable hint keyed by the status so the user knows the next
// move instead of staring at a bare "request failed" (GOAL.md §1).
type APIError struct {
	StatusCode  int         // HTTP status (often mirrors error_code)
	Code        int         // Telegram error_code
	Description string      // Telegram's human description
	Parameters  *RespParams // retry_after / migrate_to_chat_id, when present
	Body        string      // raw body, for --verbose debugging
	Method      string      // the Bot API method that failed, for context
}

// RespParams is the optional `parameters` object on an error (and some results).
type RespParams struct {
	// RetryAfter is the seconds to wait before retrying, set on a 429.
	RetryAfter int `json:"retry_after,omitempty"`
	// MigrateToChatID is the supergroup's new id when a group was upgraded.
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
}

func (e *APIError) Error() string {
	var b strings.Builder
	if e.Method != "" {
		fmt.Fprintf(&b, "%s: ", e.Method)
	}
	desc := e.Description
	if desc == "" {
		desc = "request failed"
	}
	fmt.Fprintf(&b, "Telegram API error %d: %s", e.Code, desc)
	if hint := e.hint(); hint != "" {
		fmt.Fprintf(&b, "\n  hint: %s", hint)
	}
	return b.String()
}

// hint maps the error_code (which mirrors the HTTP status) and a few well-known Telegram
// descriptions to the concrete next action. Keyed remedies are the difference between a CLI
// you can debug and one you can't. It falls back to the HTTP StatusCode when the body carried
// no usable error_code (e.g. an empty-body 5xx from a proxy), so a hint is still produced.
func (e *APIError) hint() string {
	if h := e.hintForCode(e.Code); h != "" {
		return h
	}
	if e.StatusCode != e.Code {
		return e.hintForCode(e.StatusCode)
	}
	return ""
}

func (e *APIError) hintForCode(code int) string {
	d := strings.ToLower(e.Description)
	switch {
	case code == 401 || strings.Contains(d, "unauthorized"):
		return "invalid bot token — run `tgctl auth login` (get a token from @BotFather)"
	case code == 403 || strings.Contains(d, "forbidden"):
		return "the bot lacks rights here — it must be a member/admin of the chat, and the user must have started it"
	case code == 404:
		return "method or path not found — check the method name (`tgctl api <method>`) or your --base-url"
	case code == 409 || strings.Contains(d, "terminated by other getupdates") || strings.Contains(d, "can't use getupdates"):
		return "conflict — a webhook is set or another getUpdates is running; run `tgctl webhook delete` or stop the other poller"
	case code == 429:
		if e.Parameters != nil && e.Parameters.RetryAfter > 0 {
			return fmt.Sprintf("rate limited — wait %ds and retry (tgctl backs off automatically; lower --rps for steady load)", e.Parameters.RetryAfter)
		}
		return "rate limited — slow down (lower --rps; tgctl already honors retry_after)"
	case code == 400 && strings.Contains(d, "chat not found"):
		return "chat not found — verify the chat id/@username; the bot must have interacted with the chat first"
	case code == 400 && strings.Contains(d, "message to edit not found"):
		return "no such message — list recent updates with `tgctl updates get` to find a valid message id"
	case code == 400 && strings.Contains(d, "can't parse"):
		return "markup parse error — check --parse-mode (MarkdownV2 needs special characters escaped); try plain text"
	case code == 400:
		return "bad request — re-check the parameters; `--dry-run` prints the exact request"
	case code >= 500:
		return "Telegram server error — usually transient; retry shortly"
	}
	return ""
}

// RetryAfter returns the server-requested backoff seconds, or 0 if none.
func (e *APIError) RetryAfter() int {
	if e.Parameters != nil {
		return e.Parameters.RetryAfter
	}
	return 0
}

// IsStatus reports whether the error is an APIError with the given Telegram code.
func IsStatus(err error, code int) bool {
	ae, ok := err.(*APIError)
	return ok && ae.Code == code
}
