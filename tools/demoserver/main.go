// Command demoserver answers a handful of Bot API methods with invented data, so
// the demo recording can show the CLI actually working — tables, JSON, --jq —
// without a bot, a signup, a keyring or a network.
//
// A recording made only of --dry-run shows the same curl line every time and never
// shows the tool working, which is the point of a demo.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// A stand-in Telegram Bot API: enough canned, invented answers for a recording.
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": result(method)})
	})
	addr := ":8642"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	fmt.Fprintln(os.Stderr, "stub listening on", addr)
	_ = http.ListenAndServe(addr, mux) // #nosec G114 -- a local demo stub
}

func result(method string) any {
	switch method {
	case "getMe":
		return map[string]any{"id": 8673889070, "is_bot": true, "first_name": "Demo Bot",
			"username": "demo_bot", "can_join_groups": true, "supports_inline_queries": false}
	case "sendMessage":
		return map[string]any{"message_id": 4821, "date": 1790000000,
			"chat": map[string]any{"id": -1002000000001, "title": "Demo channel", "type": "channel"},
			"text": "hello from tgctl"}
	case "getUpdates":
		return []any{
			map[string]any{"update_id": 913001, "message": map[string]any{"message_id": 4819,
				"date": 1789999100, "text": "an invented incoming message",
				"from": map[string]any{"id": 42, "first_name": "Ada", "username": "ada"},
				"chat": map[string]any{"id": 42, "type": "private", "first_name": "Ada"}}},
			map[string]any{"update_id": 913002, "message": map[string]any{"message_id": 4820,
				"date": 1789999400, "text": "and another one",
				"from": map[string]any{"id": 43, "first_name": "Linus", "username": "linus"},
				"chat": map[string]any{"id": 43, "type": "private", "first_name": "Linus"}}},
		}
	case "getChat":
		return map[string]any{"id": -1002000000001, "title": "Demo channel", "type": "channel",
			"description": "An invented channel for the demo"}
	}
	return map[string]any{}
}
