package revcatgo

import (
	"gopkg.in/guregu/null.v4"
	"time"
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
	PurchaseDate           time.Time `json:"purchase_dat"`
}

type Subscription struct {
	ExpiresDate            time.Time
	GracePeriodExpiresDate null.Time
	PurchaseDate           time.Time
	OriginalPurchaseDate   time.Time
	PeriodType             periodType
	Store                  store
	IsSandBox              bool
	UnsubscribeDetectedAt  null.Time
	BillingIssueDetectedAt null.Time
}

type NonSubscription struct {
	ID           string    `json:"id"`
	Store        store     `json:"store"`
	PurchaseDate time.Time `json:"purchase_date"`
	IsSandBox    bool      `json:"is_sandbox"`
}
