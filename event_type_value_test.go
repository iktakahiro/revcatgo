package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventType(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"INITIAL_PURCHASE", "INITIAL_PURCHASE", nil},
		{"CANCELLATION", "CANCELLATION", nil},
		{"INVALID", "", errors.New("eventType value should be one of the following: TEST, INITIAL_PURCHASE, NON_RENEWING_PURCHASE, RENEWAL, PRODUCT_CHANGE, CANCELLATION, UNCANCELLATION, BILLING_ISSUE, SUBSCRIBER_ALIAS, SUBSCRIPTION_PAUSED, TRANSFER, EXPIRATION, got INVALID")},
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
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{`"INITIAL_PURCHASE"`, "INITIAL_PURCHASE", nil},
		{`"CANCELLATION"`, "CANCELLATION", nil},
		{`"INVALID"`, "", errors.New("")},
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
