package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// No path out of this process may print a bot token, whatever produced the error (issue #21).
func TestErrorLine_RedactsToken(t *testing.T) {
	const fakeToken = "123456789:AAFfakeFAKEfakeFAKEfakeFAKEfake-03"

	line := errorLine(errors.New(`sendMessage: Post "https://api.telegram.org/bot` + fakeToken + `/sendMessage": read: connection timed out`))

	assert.NotContains(t, line, fakeToken)
	assert.Contains(t, line, "Error: sendMessage: ")
	assert.Contains(t, line, "123456789:<redacted>")
}

func TestErrorLine_LeavesOrdinaryErrorsAlone(t *testing.T) {
	assert.Equal(t, "Error: chat not found", errorLine(errors.New("chat not found")))
}
