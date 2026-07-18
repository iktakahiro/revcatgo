package revcatgo

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

const (
	// EventTypeTest indicates a RevenueCat test webhook event.
	EventTypeTest = "TEST"
	// EventTypeInitialPurchase indicates the initial purchase of a product.
	EventTypeInitialPurchase = "INITIAL_PURCHASE"
	// EventTypeNonRenewingPurchase indicates a non-renewing purchase event.
	EventTypeNonRenewingPurchase = "NON_RENEWING_PURCHASE"
	// EventTypeRenewal indicates that a subscription renewed.
	EventTypeRenewal = "RENEWAL"
	// EventTypeProductChange indicates the customer changed products.
	EventTypeProductChange = "PRODUCT_CHANGE"
	// EventTypeCancellation indicates that a subscription was canceled.
	EventTypeCancellation = "CANCELLATION"
	// EventTypeUnCancellation indicates that a cancellation was reversed.
	EventTypeUnCancellation = "UNCANCELLATION"
	// EventTypeBillingIssue indicates a billing issue occurred.
	EventTypeBillingIssue = "BILLING_ISSUE"
	// EventTypeSubscriberAlias indicates subscriber aliasing activity.
	EventTypeSubscriberAlias = "SUBSCRIBER_ALIAS"
	// EventTypeSubscriptionPaused indicates the subscription was paused.
	EventTypeSubscriptionPaused = "SUBSCRIPTION_PAUSED"
	// EventTypeTransfer indicates the subscription was transferred.
	EventTypeTransfer = "TRANSFER"
	// EventTypeExpiration indicates the subscription expired.
	EventTypeExpiration = "EXPIRATION"
	// EventTypeSubscriptionExtended indicates the subscription was extended.
	EventTypeSubscriptionExtended = "SUBSCRIPTION_EXTENDED"
	// EventTypeTemporaryEntitlementGrant indicates a temporary entitlement grant.
	EventTypeTemporaryEntitlementGrant = "TEMPORARY_ENTITLEMENT_GRANT"
	// EventTypeRefundReversed indicates a previously issued refund was reversed.
	EventTypeRefundReversed = "REFUND_REVERSED"
	// EventTypeInvoiceIssuance indicates a new invoice has been issued for billing.
	EventTypeInvoiceIssuance = "INVOICE_ISSUANCE"
	// EventTypeVirtualCurrencyTransaction indicates a virtual currency adjustment event.
	EventTypeVirtualCurrencyTransaction = "VIRTUAL_CURRENCY_TRANSACTION"
)

var validEventTypeValues = []string{
	EventTypeTest,
	EventTypeInitialPurchase,
	EventTypeNonRenewingPurchase,
	EventTypeRenewal,
	EventTypeProductChange,
	EventTypeCancellation,
	EventTypeUnCancellation,
	EventTypeBillingIssue,
	EventTypeSubscriberAlias,
	EventTypeSubscriptionPaused,
	EventTypeTransfer,
	EventTypeExpiration,
	EventTypeSubscriptionExtended,
	EventTypeTemporaryEntitlementGrant,
	EventTypeRefundReversed,
	EventTypeInvoiceIssuance,
	EventTypeVirtualCurrencyTransaction,
}

type eventType struct {
	value null.String
}

func newEventType(s string) (*eventType, error) {
	if !contains(validEventTypeValues, s) {
		return &eventType{}, fmt.Errorf("eventType value should be one of the following: %v, got %v", strings.Join(validEventTypeValues, ", "), s)
	}
	return &eventType{value: null.StringFrom(s)}, nil
}

func (e eventType) String() string {
	return e.value.ValueOrZero()
}

// MarshalJSON serializes an event type to JSON.
func (e eventType) MarshalJSON() ([]byte, error) {
	return e.value.MarshalJSON()
}

// UnmarshalJSON deserializes an event type from JSON.
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
