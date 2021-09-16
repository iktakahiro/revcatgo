package revcatgo

import (
	"errors"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

const (
	EventTypeTest                = "TEST"
	EventTypeInitialPurchase     = "INITIAL_PURCHASE"
	EventTypeNonRenewingPurchase = "NON_RENEWING_PURCHASE"
	EventTypeRenewal             = "RENEWAL"
	EventTypeProductChange       = "PRODUCT_CHANGE"
	EventTypeChancellation       = "CANCELLATION"
	EventTypeUnChancellation     = "UNCANCELLATION"
	EventTypeBillingIssue        = "BILLING_ISSUE"
	EventTypeSubscriberAlias     = "SUBSCRIBER_ALIAS"
	EventTypeSubscriptionPaused  = "SUBSCRIPTION_PAUSED"
	EventTypeTransfer            = "TRANSFER"
	EventTypeExpiration          = "EXPIRATION"
	// use this when receving undefined eventType
	EventTypeUnknown = "UNKNOWN"
)

var validEventTypeValues = []string{
	EventTypeTest,
	EventTypeInitialPurchase,
	EventTypeNonRenewingPurchase,
	EventTypeRenewal,
	EventTypeProductChange,
	EventTypeChancellation,
	EventTypeUnChancellation,
	EventTypeBillingIssue,
	EventTypeSubscriberAlias,
	EventTypeSubscriptionPaused,
	EventTypeTransfer,
	EventTypeExpiration,
}

type eventType struct {
	value null.String
}

func newEventType(s string) (*eventType, error) {
	if !contains(validEventTypeValues, s) {
		return &eventType{value: null.StringFrom(EventTypeUnknown)}, nil
	}
	return &eventType{value: null.StringFrom(s)}, nil
}

func (e eventType) String() string {
	return e.value.ValueOrZero()
}

// MarshalJSON serializes a store to JSON.
func (e eventType) MarshalJSON() ([]byte, error) {
	return e.value.MarshalJSON()
}

// UnmarshalJSON deserialized a store from JSON
func (e *eventType) UnmarshalJSON(b []byte) error {
	v := &eventType{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of type: %w", err)
	}
	if !v.value.Valid {
		return errors.New("type is a required field")
	}
	_e, err := newEventType(v.value.ValueOrZero())
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of type: %w", err)
	}
	e.value = _e.value

	return nil
}
