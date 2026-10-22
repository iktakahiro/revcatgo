package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCancelReason(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"UNSUBSCRIBE", "UNSUBSCRIBE", nil},
		{"BILLING_ERROR", "BILLING_ERROR", nil},
		{"INVALID", "", errors.New("cancelReason value should be one of the following: UNSUBSCRIBE,BILLING_ERROR,DEVELOPER_INITIATED,PRICE_INCREASE,CUSTOMER_SUPPORT,SUBSCRIPTION_PAUSED,UNKNOWN, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newCancelReason(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestCancelReasonUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		known    bool
		wantErr  bool
	}{
		{`"UNSUBSCRIBE"`, "UNSUBSCRIBE", true, false},
		{`"BILLING_ERROR"`, "BILLING_ERROR", true, false},
		{`"NEW_REASON"`, "NEW_REASON", false, false},
		{`"unsubscribe"`, "unsubscribe", false, false},
		{`""`, "", false, false},
		{`null`, "", false, false},
		{`1`, "", false, true},
	}

	for _, c := range cases {
		var cr cancelReason
		err := json.Unmarshal([]byte(c.in), &cr)
		if c.wantErr {
			require.Error(t, err, c.in)
			continue
		}
		require.NoError(t, err, c.in)
		assert.Equal(t, c.expected, cr.String(), c.in)
		assert.Equal(t, c.known, cr.IsKnown(), c.in)
	}
}

func TestCancelReasonIsKnown(t *testing.T) {
	known, err := newCancelReason(CancelReasonUnsubscribe)
	require.NoError(t, err)
	assert.True(t, known.IsKnown())

	var unknown cancelReason
	require.NoError(t, json.Unmarshal([]byte(`"NEW_REASON"`), &unknown))
	assert.False(t, unknown.IsKnown())
	assert.Equal(t, "NEW_REASON", unknown.String())

	var zero cancelReason
	assert.False(t, zero.IsKnown())
}

func TestCancelReasonMarshalRoundTrip(t *testing.T) {
	for _, in := range []string{
		`"UNSUBSCRIBE"`,
		`"NEW_REASON"`,
		`null`,
	} {
		var cr cancelReason
		require.NoError(t, json.Unmarshal([]byte(in), &cr))

		out, err := json.Marshal(cr)
		require.NoError(t, err)
		assert.JSONEq(t, in, string(out))
	}
}

func TestCancelReasonNullString(t *testing.T) {
	c, err := newCancelReason(CancelReasonCustomerSupport)
	require.NoError(t, err)

	nullString := c.NullString()
	assert.True(t, nullString.Valid)
	assert.Equal(t, CancelReasonCustomerSupport, nullString.ValueOrZero())
}
