package revcatgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

// SubscriberResponse represents the wrapper returned by RevenueCat's Subscriber API v1.
type SubscriberResponse struct {
	RequestDate   time.Time    `json:"request_date"`
	RequestDateAt milliseconds `json:"request_date_ms"`
	Subscriber    Subscriber   `json:"subscriber"`
}

// Subscriber aggregates subscription data for a single app user.
type Subscriber struct {
	Entitlements               map[string]Entitlement       `json:"entitlements"`
	FirstSeen                  time.Time                    `json:"first_seen"`
	LastSeen                   time.Time                    `json:"last_seen"`
	ManagementURL              null.String                  `json:"management_url"`
	OriginalAppUserID          null.String                  `json:"original_app_user_id"`
	OriginalApplicationVersion null.String                  `json:"original_application_version"`
	OriginalPurchaseDate       null.Time                    `json:"original_purchase_date"`
	SubscriberAttributes       subscriberAttributes         `json:"subscriber_attributes"`
	OtherPurchases             map[string]any               `json:"other_purchases"`
	Subscriptions              map[string]Subscription      `json:"subscriptions"`
	NonSubscription            map[string][]NonSubscription `json:"non_subscriptions"`
}

// Entitlement represents a subscriber entitlement granted for a product.
type Entitlement struct {
	ProductIdentifier      string    `json:"product_identifier"`
	ExpiresDate            time.Time `json:"expires_date"`
	GracePeriodExpiresDate null.Time `json:"grace_period_expires_date"`
	PurchaseDate           time.Time `json:"purchase_date"`
}

// Subscription captures the state of an individual auto-renewing subscription.
type Subscription struct {
	AutoResumeDate         null.Time          `json:"auto_resume_date"`
	ExpiresDate            time.Time          `json:"expires_date"`
	GracePeriodExpiresDate null.Time          `json:"grace_period_expires_date"`
	PurchaseDate           time.Time          `json:"purchase_date"`
	OriginalPurchaseDate   time.Time          `json:"original_purchase_date"`
	OwnershipType          string             `json:"ownership_type"`
	PeriodType             periodType         `json:"period_type"`
	Store                  store              `json:"store"`
	IsSandBox              bool               `json:"is_sandbox"`
	UnsubscribeDetectedAt  null.Time          `json:"unsubscribe_detected_at"`
	BillingIssueDetectedAt null.Time          `json:"billing_issues_detected_at"`
	RefundedAt             null.Time          `json:"refunded_at"`
	StoreTransactionID     storeTransactionID `json:"store_transaction_id"`
}

// NonSubscription describes a one-off, non-renewing purchase.
type NonSubscription struct {
	ID           string    `json:"id"`
	Store        store     `json:"store"`
	PurchaseDate time.Time `json:"purchase_date"`
	IsSandBox    bool      `json:"is_sandbox"`
}

type storeTransactionID struct {
	value string
}

// String returns the store transaction identifier as a string for every store.
func (s storeTransactionID) String() string {
	return s.value
}

func (s storeTransactionID) MarshalJSON() ([]byte, error) {
	if s.value == "" {
		return []byte("null"), nil
	}

	return json.Marshal(s.value)
}

func (s *storeTransactionID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		s.value = ""

		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		s.value = stringValue

		return nil
	}

	var numberValue json.Number
	if err := json.Unmarshal(data, &numberValue); err != nil {
		return fmt.Errorf("failed to unmarshal store_transaction_id: %w", err)
	}
	s.value = numberValue.String()

	return nil
}
