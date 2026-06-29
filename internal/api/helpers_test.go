package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID_String(t *testing.T) {
	assert.Equal(t, "42", ID("42").String())
}

func TestInt_MarshalAndInt64(t *testing.T) {
	b, _ := json.Marshal(Int(-5))
	assert.Equal(t, "-5", string(b))
	assert.Equal(t, int64(-5), Int(-5).Int64())
}

func TestUser_DisplayName_NoUsername(t *testing.T) {
	u := User{FirstName: "Ada", LastName: "Lovelace"}
	assert.Equal(t, "Ada Lovelace", u.DisplayName())
	assert.Equal(t, "@bot", User{Username: "bot"}.DisplayName())
}

func TestScalarString(t *testing.T) {
	assert.Equal(t, "x", scalarString("x"))
	assert.Equal(t, "", scalarString(nil))
	assert.Equal(t, `{"a":1}`, scalarString(json.RawMessage(`{"a":1}`)))
	assert.Equal(t, "42", scalarString(42))
}

func TestClient_BaseURLAndOptions(t *testing.T) {
	a, _ := NewBotTokenAuth("1:x")
	c := New(a, WithBaseURL("https://example.com"), WithVerbose(true))
	assert.Equal(t, "https://example.com", c.BaseURL())
	assert.True(t, c.Verbose)
}
