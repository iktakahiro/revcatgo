package revcatgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

// WebhookEvent represents a request body of RevenueCat webhook.
// https://www.revenuecat.com/docs/integrations/webhooks
type WebhookEvent struct {
	Event      Event  `json:"event"`
	APIVersion string `json:"api_version"`
}

// Event represents an Event of RevenueCat webhook
type Event struct {
	ID                       string               `json:"id"`
	AppID                    string               `json:"app_id"`
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
	TaxPercentage            float32              `json:"tax_percentage"`
	TakeHomePercentage       float32              `json:"takehome_percentage"`
	CommissionPercentage     float32              `json:"commission_percentage"`
	SubscriberAttributes     subscriberAttributes `json:"subscriber_attributes"`
	Experiments              []experiment         `json:"experiments"`
	TransactionID            string               `json:"transaction_id"`
	OriginalTransactionID    string               `json:"original_transaction_id"`
	IsFamilyShare            bool                 `json:"is_family_share"`
	TransferredFrom          []string             `json:"transferred_from"`
	TransferredTo            []string             `json:"transferred_to"`
	CountryCode              string               `json:"country_code"`
	OfferCode                string               `json:"offer_code"`
	RenewalNumber            int                  `json:"renewal_number"`
	Metadata                 map[string]any       `json:"metadata"`
	DiscountPercentage       null.Float           `json:"discount_percentage"`
	DiscountAmount           null.Float           `json:"discount_amount"`
	DiscountIdentifier       null.String          `json:"discount_identifier"`
	Quantity                 null.Int             `json:"quantity"`
	Adjustments              []virtualAdjustment  `json:"adjustments"`
	ProductDisplayName       string               `json:"product_display_name"`
	PurchaseEnvironment      environment          `json:"purchase_environment"`
	Source                   string               `json:"source"`
	VirtualTransactionID     string               `json:"virtual_currency_transaction_id"`
	UpdatedBalance           null.Int             `json:"updated_balance"`
	AdTransactionID          string               `json:"ad_transaction_id"`
	ExperimentID             string               `json:"experiment_id"`
	ExperimentVariant        string               `json:"experiment_variant"`
	ExperimentEnrolledAt     milliseconds         `json:"experiment_enrolled_at_ms"`
	OfferingID               string               `json:"offering_id"`
	RedeemedFrom             []string             `json:"redeemed_from"`
	RedeemedBy               []string             `json:"redeemed_by"`
	RedemptionOutcome        string               `json:"redemption_outcome"`
	RedemptionPlatform       null.String          `json:"redemption_platform"`
	WorkflowID               string               `json:"workflow_id"`
	WorkflowStepID           string               `json:"workflow_step_id"`
	TraceID                  string               `json:"trace_id"`
	PaywallEventID           string               `json:"event_id"`
	Platform                 string               `json:"platform"`
	PlatformVersion          string               `json:"platform_version"`
	SDKVersion               string               `json:"sdk_version"`
	PaywallID                string               `json:"paywall_id"`
	PaywallName              string               `json:"paywall_name"`
	SessionID                string               `json:"session_id"`
	DisplayMode              null.String          `json:"display_mode"`
	DarkMode                 null.Bool            `json:"dark_mode"`
	Locale                   null.String          `json:"locale"`
	ComponentType            string               `json:"component_type"`
	ComponentValue           string               `json:"component_value"`
	ComponentName            string               `json:"component_name"`
	ComponentURL             string               `json:"component_url"`
	OriginIndex              null.Int             `json:"origin_index"`
	DestinationIndex         null.Int             `json:"destination_index"`
	OriginContextName        string               `json:"origin_context_name"`
	DestinationContextName   string               `json:"destination_context_name"`
	DefaultIndex             null.Int             `json:"default_index"`
	OriginPackageID          string               `json:"origin_package_id"`
	DestinationPackageID     string               `json:"destination_package_id"`
	DefaultPackageID         string               `json:"default_package_id"`
	CurrentPackageID         string               `json:"current_package_id"`
	ResultingPackageID       string               `json:"resulting_package_id"`
	OriginProductID          string               `json:"origin_product_id"`
	DestinationProductID     string               `json:"destination_product_id"`
	DefaultProductID         string               `json:"default_product_id"`
	CurrentProductID         string               `json:"current_product_id"`
	ResultingProductID       string               `json:"resulting_product_id"`
}

// IsExpired checks whether a subscription is expired or not.
func (e *Event) IsExpired(grace time.Duration, base *time.Time) bool {
	if !e.ExpirationAt.NullInt().Valid {
		return false
	}

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
	flattened bool
}

func (s subscriberAttribute) MarshalJSON() ([]byte, error) {
	if s.flattened {
		return json.Marshal(s.Value)
	}

	type attribute subscriberAttribute

	return json.Marshal(attribute(s))
}

func (s *subscriberAttribute) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) > 0 && data[0] == '{' {
		type attribute subscriberAttribute
		var value attribute
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("failed to unmarshal subscriber attribute: %w", err)
		}
		*s = subscriberAttribute(value)

		return nil
	}

	if err := json.Unmarshal(data, &s.Value); err != nil {
		return fmt.Errorf("failed to unmarshal flattened subscriber attribute: %w", err)
	}
	s.flattened = true

	return nil
}

// experiment represents a single experiment enrollment attached to the event.
type experiment struct {
	ID         string       `json:"experiment_id"`
	Variant    string       `json:"experiment_variant"`
	EnrolledAt milliseconds `json:"enrolled_at_ms"`
}

// virtualAdjustment captures adjustments part of virtual currency transactions.
type virtualAdjustment struct {
	Amount   int             `json:"amount"`
	Currency virtualCurrency `json:"currency"`
}

// virtualCurrency describes the virtual currency metadata in adjustments.
type virtualCurrency struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
