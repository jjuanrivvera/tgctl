package api

import (
	"context"
	"encoding/json"
)

// Recorder observes successful Bot API calls, so a caller can persist a local history of
// everything tgctl sent (issue #5) without internal/api depending on internal/store — the
// client stays a generic HTTP core (GOAL.md §2; AGENTS.md). Record is called once per
// successful, non-dry-run Call/Upload; ctx is the same context the triggering command passed
// in, so an implementation can honor cancellation without ever synthesizing a background one.
//
// Record must not block the caller on its own failures: an implementation is responsible for
// filtering to the methods it cares about and for swallowing/logging its own errors — a broken
// message store must never fail a send.
type Recorder interface {
	Record(ctx context.Context, method string, params map[string]any, result json.RawMessage)
}
