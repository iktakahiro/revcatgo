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
		{"VIRTUAL_CURRENCY_TRANSACTION", "VIRTUAL_CURRENCY_TRANSACTION", nil},
		{"EXPERIMENT_ENROLLMENT", "EXPERIMENT_ENROLLMENT", nil},
		{"PAYWALL_COMPONENT_INTERACTED", "PAYWALL_COMPONENT_INTERACTED", nil},
		{"PRICE_INCREASE_CONSENT_REQUIRED", "PRICE_INCREASE_CONSENT_REQUIRED", nil},
		{"INVALID", "", errors.New("eventType value should be one of the following: TEST, INITIAL_PURCHASE, NON_RENEWING_PURCHASE, RENEWAL, PRODUCT_CHANGE, CANCELLATION, UNCANCELLATION, BILLING_ISSUE, SUBSCRIBER_ALIAS, SUBSCRIPTION_PAUSED, TRANSFER, EXPIRATION, SUBSCRIPTION_EXTENDED, TEMPORARY_ENTITLEMENT_GRANT, REFUND_REVERSED, INVOICE_ISSUANCE, VIRTUAL_CURRENCY_TRANSACTION, EXPERIMENT_ENROLLMENT, PURCHASE_REDEEMED, PAYWALL_IMPRESSION, PAYWALL_CLOSE, PAYWALL_CANCEL, PAYWALL_EXIT_OFFER, PAYWALL_COMPONENT_INTERACTED, PRICE_INCREASE_CONSENT_REQUIRED, PRICE_INCREASE_CONSENT_APPROVED, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newEventType(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.NoError(t, err)
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
		{`"MY_CUSTOM_PAYWALL_EVENT"`, "MY_CUSTOM_PAYWALL_EVENT", nil},
		{`""`, "", errors.New("")},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var e eventType
		b := []byte(c.in)
		err := json.Unmarshal(b, &e)

		assert.Equal(t, c.expected, e.String())
		if c.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestEventTypeIsKnown(t *testing.T) {
	var known eventType
	assert.NoError(t, json.Unmarshal([]byte(`"PURCHASE_REDEEMED"`), &known))
	assert.True(t, known.IsKnown())

	var custom eventType
	assert.NoError(t, json.Unmarshal([]byte(`"MY_CUSTOM_PAYWALL_EVENT"`), &custom))
	assert.False(t, custom.IsKnown())
}
