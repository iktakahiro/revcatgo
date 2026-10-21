package revcatgo

import (
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type cancelReason struct {
	value null.String
}

const (
	// CancelReasonUnsubscribe indicates that the user intentionally unsubscribed.
	CancelReasonUnsubscribe = "UNSUBSCRIBE"
	// CancelReasonBillingError indicates the subscription was canceled due to billing errors.
	CancelReasonBillingError = "BILLING_ERROR"
	// CancelReasonDeveloperInitiated indicates the developer initiated the cancellation.
	CancelReasonDeveloperInitiated = "DEVELOPER_INITIATED"
	// CancelReasonPriceIncrease indicates the user rejected a price increase.
	CancelReasonPriceIncrease = "PRICE_INCREASE"
	// CancelReasonCustomerSupport indicates cancellation via customer support.
	CancelReasonCustomerSupport = "CUSTOMER_SUPPORT"
	// CancelReasonSubscriptionPaused indicates the subscription was paused.
	CancelReasonSubscriptionPaused = "SUBSCRIPTION_PAUSED"
	// CancelReasonUnknown indicates the cancellation reason was not provided.
	CancelReasonUnknown = "UNKNOWN"
)

var validCancelReasonValues = []string{
	CancelReasonUnsubscribe,
	CancelReasonBillingError,
	CancelReasonDeveloperInitiated,
	CancelReasonPriceIncrease,
	CancelReasonCustomerSupport,
	CancelReasonSubscriptionPaused,
	CancelReasonUnknown,
}

func newCancelReason(v string) (*cancelReason, error) {
	if !contains(validCancelReasonValues, v) {
		return &cancelReason{}, fmt.Errorf("cancelReason value should be one of the following: %v, got %v", strings.Join(validCancelReasonValues, ","), v)
	}
	return &cancelReason{value: null.StringFrom(v)}, nil
}

func (c cancelReason) String() string {
	return c.value.ValueOrZero()
}

func (c cancelReason) NullString() null.String {
	return c.value
}

// IsKnown reports whether RevenueCat documents the reason as one of the CancelReason constants.
// Reasons added by RevenueCat later (and null values) return false but remain decodable.
func (c cancelReason) IsKnown() bool {
	return contains(validCancelReasonValues, c.value.ValueOrZero())
}

// MarshalJSON serializes a cancellation reason to JSON.
func (c cancelReason) MarshalJSON() ([]byte, error) {
	return c.value.MarshalJSON()
}

// UnmarshalJSON deserializes a cancellation reason from JSON.
// JSON null and empty strings decode to the zero value, and undocumented values are
// preserved so that new RevenueCat values do not break decoding; see IsKnown.
func (c *cancelReason) UnmarshalJSON(b []byte) error {
	v, err := unmarshalEnumValue(b, false)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of cancel_reason: %w", err)
	}
	c.value = v

	return nil
}
