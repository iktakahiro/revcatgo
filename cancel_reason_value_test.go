package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCancelReason(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"UNSUBSCRIBE", "UNSUBSCRIBE", nil},
		{"BILLING_ERROR", "BILLING_ERROR", nil},
		{"INVALID", "", errors.New("cancelReason value should be one of the following: UNSUBSCRIBE,BILLING_ERROR,DEVELOPER_INITIATED,PRICE_INCREASE,CUSTOMER_SUPPORT,UNKNOWN, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newCancelReason(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestCancelReasonUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{`"UNSUBSCRIBE"`, "UNSUBSCRIBE", nil},
		{`"BILLING_ERROR"`, "BILLING_ERROR", nil},
		{`"INVALID"`, "", errors.New("")},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var cr cancelReason
		b := []byte(c.in)
		err := json.Unmarshal(b, &cr)

		assert.Equal(t, c.expected, cr.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
