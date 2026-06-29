package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_HintsByStatus(t *testing.T) {
	cases := []struct {
		code   int
		desc   string
		params *RespParams
		want   string
	}{
		{401, "Unauthorized", nil, "auth login"},
		{403, "Forbidden: bot was blocked", nil, "member/admin"},
		{404, "Not Found", nil, "method name"},
		{409, "Conflict: terminated by other getUpdates request", nil, "webhook"},
		{429, "Too Many Requests", &RespParams{RetryAfter: 5}, "wait 5s"},
		{400, "Bad Request: chat not found", nil, "chat id"},
		{400, "Bad Request: message to edit not found", nil, "updates get"},
		{500, "Internal Server Error", nil, "transient"},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			e := &APIError{Code: tc.code, Description: tc.desc, Parameters: tc.params, Method: "m"}
			assert.Contains(t, e.Error(), tc.want)
			assert.Contains(t, e.Error(), "m:")
		})
	}
}

func TestAPIError_RetryAfterAndIsStatus(t *testing.T) {
	e := &APIError{Code: 429, Parameters: &RespParams{RetryAfter: 3}}
	assert.Equal(t, 3, e.RetryAfter())
	assert.True(t, IsStatus(e, 429))
	assert.False(t, IsStatus(e, 400))

	none := &APIError{Code: 400}
	assert.Equal(t, 0, none.RetryAfter())
}

func TestAPIError_EmptyDescriptionFallsBack(t *testing.T) {
	e := &APIError{Code: 400}
	assert.Contains(t, e.Error(), "request failed")
}
