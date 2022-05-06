package revcatgo

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// WebhookEvent represents a request body of RevenueCat webhook.
// https://docs.revenuecat.com/docs/webhooks
type WebhookEvent struct {
	Event      Event  `json:"event"`
	APIVersion string `json:"api_version"`
}

// Event represents an Event of RevenueCat webhook
type Event struct {
	ID                       string               `json:"id"`
	Type                     eventType            `json:"type"`
	EventTimestampAt         milliseconds         `json:"event_timestamp_ms"`
	AppUserID                string               `json:"app_user_id"`
	Aliases                  []string             `json:"aliases"`
	OriginalAppUserID        string               `json:"original_app_user_id"`
	ProductID                string               `json:"product_id"`
	EntitlementIDs           []string             `json:"entitlement_ids"`
	PeriodType               periodType           `json:"period_type"`
	PurchasedAt              milliseconds         `json:"purchased_at_ms"`
	GracePeriodExpirationAt  milliseconds         `json:"grace_period_expiration_at_ms"`
	ExpirationAt             milliseconds         `json:"expiration_at_ms"`
	AutoResumeAt             milliseconds         `json:"auto_resume_at_ms"`
	Store                    store                `json:"store"`
	Environment              environment          `json:"environment"`
	IsTrialConversion        null.Bool            `json:"is_trial_conversion"`
	CancelReason             cancelReason         `json:"cancel_reason"`
	ExpirationReason         cancelReason         `json:"expiration_reason"`
	NewProductID             string               `json:"new_product_id"`
	PresentedOfferingID      string               `json:"presented_offering_id"`
	Price                    price                `json:"price"`
	Currency                 null.String          `json:"currency"`
	PriceInPurchasedCurrency float32              `json:"price_in_purchased_currency"`
	TakeHomePercentage       float32              `json:"takehome_percentage"`
	SubscriberAttributes     subscriberAttributes `json:"subscriber_attributes"`
	TransactionID            string               `json:"transaction_id"`
	OriginalTransactionID    string               `json:"original_transaction_id"`
}

// IsExpired checks whether a subscription is expired or not.
func (e *Event) IsExpired(grace time.Duration, base *time.Time) bool {
	var b time.Time
	if base == nil {
		b = time.Now()
	} else {
		b = *base
	}
	return e.ExpirationAt.DateTime().Add(grace).Before(b.UTC())
}

// HasEntitlementID checks whether the id exists or not.
func (e *Event) HasEntitlementID(id string) bool {
	return contains(e.EntitlementIDs, id)
}

// GetAllRelatedUserID returns a unique id list of AppUserID, OriginalAppUserID, and Aliases.
func (e *Event) GetAllRelatedUserID() []string {
	m := make(map[string]bool)
	var idList []string
	if e.AppUserID != "" {
		idList = append(idList, e.AppUserID)
	}
	if e.OriginalAppUserID != "" {
		idList = append(idList, e.OriginalAppUserID)
	}
	if len(e.Aliases) > 0 {
		idList = append(idList, e.Aliases...)
	}
	var uniqueIDList []string

	for _, ele := range idList {
		if !m[ele] {
			m[ele] = true
			uniqueIDList = append(uniqueIDList, ele)
		}
	}

	return uniqueIDList
}

// SubscriberAttributes represents a map of SubscriberAttribute.
type subscriberAttributes map[string]subscriberAttribute

// SubscriberAttribute represents attributes of subscriber.
type subscriberAttribute struct {
	Value     string       `json:"value"`
	UpdatedAt milliseconds `json:"updated_at_ms"`
}
