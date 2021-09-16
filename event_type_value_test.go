package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventType(t *testing.T) {
	var cases = []struct {
		in       string
		expected string
		err      error
	}{
		{"INITIAL_PURCHASE", "INITIAL_PURCHASE", nil},
		{"CANCELLATION", "CANCELLATION", nil},
		{"INVALID", "UNKNOWN", nil},
	}

	for _, c := range cases {
		actual, err := newEventType(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestEventTypeUnMarshal(t *testing.T) {
	var cases = []struct {
		in       string
		expected string
		err      error
	}{
		{`"INITIAL_PURCHASE"`, "INITIAL_PURCHASE", nil},
		{`"CANCELLATION"`, "CANCELLATION", nil},
		{`"INVALID"`, "UNKNOWN", nil},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var e eventType
		b := []byte(c.in)
		err := json.Unmarshal(b, &e)

		assert.Equal(t, c.expected, e.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
