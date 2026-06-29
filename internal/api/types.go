package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Flexible JSON scalar types. Telegram is fairly disciplined about its JSON, but IDs are
// large int64s (chat_id / user_id can exceed 2^53) and a few fields arrive as either a
// number or a string depending on the method. Decoding through these types prevents the
// most common real-world breakages (precision loss, type drift) — see GOAL.md §2.

// ID unmarshals from a JSON string OR number and always marshals as a string. Telegram
// chat/user/message IDs are int64 that overflow JS's 2^53 safe integer, so we never let
// them round-trip through float64 and we render them consistently in tables.
type ID string

// UnmarshalJSON accepts `123`, `"123"`, or null.
func (id *ID) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*id = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*id = ID(s)
		return nil
	}
	// Validate it is a JSON number (reject 1e9 style that would lose the integer form),
	// then keep the exact textual integer.
	if !isJSONInteger(b) {
		return fmt.Errorf("id: not an integer: %s", b)
	}
	*id = ID(b)
	return nil
}

// MarshalJSON always emits a quoted string so tables and JSON output are stable.
func (id ID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }

func (id ID) String() string { return string(id) }

// Int accepts a JSON number or a numeric string and stores an exact int64. It decodes the
// integer form before any float path so values above 2^53 keep full precision, and it
// rejects NaN/Inf and non-integer numbers.
type Int int64

func (n *Int) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*n = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("int: %w", err)
		}
		*n = Int(v)
		return nil
	}
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		// Not a bare integer — reject floats/NaN/Inf explicitly rather than truncating.
		return fmt.Errorf("int: not an integer: %s", b)
	}
	*n = Int(v)
	return nil
}

func (n Int) MarshalJSON() ([]byte, error) { return []byte(strconv.FormatInt(int64(n), 10)), nil }

func (n Int) Int64() int64 { return int64(n) }

// Bool accepts a real JSON bool or the strings "true"/"false"/"1"/"0"/"yes"/"no".
type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*b = false
		return nil
	}
	switch string(data) {
	case "true", `"true"`, `"1"`, `"yes"`:
		*b = true
		return nil
	case "false", `"false"`, `"0"`, `"no"`:
		*b = false
		return nil
	}
	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("bool: %w", err)
	}
	*b = Bool(v)
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) { return json.Marshal(bool(b)) }

// isJSONInteger reports whether b is a JSON number with no fraction/exponent and finite.
func isJSONInteger(b []byte) bool {
	s := string(b)
	if s == "" {
		return false
	}
	// json.Number keeps the textual form; ParseInt confirms it is a pure integer.
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return true
	}
	// Reject anything non-integer (floats, NaN, Inf, 1e3) — they have no business as an ID.
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}
	return false
}
