package revcatgo

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SubscriberResponse struct {
	RequestDateAt milliseconds `json:"request_date_ms"`
	Subscriber    Subscriber   `json:"subscriber"`
}

type Subscriber struct {
	Entitlements      map[string]Entitlement     `json:"entitlements"`
	FirstSeen         time.Time                  `json:"first_seen"`
	LastSeen          time.Time                  `json:"last_seen"`
	ManagementURL     null.String                `json:"management_url"`
	OriginalAppUserID null.String                `json:"original_app_user_id"`
	Subscriptions     map[string]Subscription    `json:"subscriptions"`
	NonSubscription   map[string]NonSubscription `json:"non_subscriptions"`
}

type Entitlement struct {
	ProductIdentifier      string    `json:"product_identifier"`
	ExpiresDate            time.Time `json:"expires_date"`
	GracePeriodExpiresDate null.Time `json:"grace_period_expires_date"`
	PurchaseDate           time.Time `json:"purchase_date"`
}

type Subscription struct {
	ExpiresDate            time.Time  `json:"expires_date"`
	GracePeriodExpiresDate null.Time  `json:"grace_period_expires_date"`
	PurchaseDate           time.Time  `json:"purchase_date"`
	OriginalPurchaseDate   time.Time  `json:"original_purchase_date"`
	PeriodType             periodType `json:"period_type"`
	Store                  store      `json:"store"`
	IsSandBox              bool       `json:"is_sandbox"`
	UnsubscribeDetectedAt  null.Time  `json:"unsubscribe_detected_at"`
	BillingIssueDetectedAt null.Time  `json:"billing_issue_detected_at"`
}

type NonSubscription struct {
	ID           string    `json:"id"`
	Store        store     `json:"store"`
	PurchaseDate time.Time `json:"purchase_date"`
	IsSandBox    bool      `json:"is_sandbox"`
}
