package revcatgo

// Virtual currency transaction sources documented at https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields.
const (
	// VirtualCurrencySourceInAppPurchase identifies an in-app purchase.
	VirtualCurrencySourceInAppPurchase = "in_app_purchase"
	// VirtualCurrencySourceAdminAPI identifies an admin API adjustment.
	VirtualCurrencySourceAdminAPI = "admin_api"
	// VirtualCurrencySourceServerAPI identifies a server API adjustment.
	VirtualCurrencySourceServerAPI = "server_api"
	// VirtualCurrencySourceClientSDK identifies a client SDK adjustment.
	VirtualCurrencySourceClientSDK = "client_sdk"
	// VirtualCurrencySourceExpiration identifies an expiration.
	VirtualCurrencySourceExpiration = "expiration"
	// VirtualCurrencySourceAdReward identifies an ad reward.
	VirtualCurrencySourceAdReward = "ad_reward"
)

// Redemption outcomes documented at https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields.
const (
	// RedemptionOutcomeAlias indicates that the identities were aliased.
	RedemptionOutcomeAlias = "alias"
	// RedemptionOutcomeTransfer indicates that the purchase was transferred.
	RedemptionOutcomeTransfer = "transfer"
	// RedemptionOutcomeRedeemerOwns indicates that the redeemer already owns the purchase.
	RedemptionOutcomeRedeemerOwns = "redeemer_owns"
)

// IsKnownVirtualCurrencySource reports whether RevenueCat documents source.
// Unknown values remain decodable. See https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields.
func IsKnownVirtualCurrencySource(source string) bool {
	switch source {
	case VirtualCurrencySourceInAppPurchase, VirtualCurrencySourceAdminAPI, VirtualCurrencySourceServerAPI,
		VirtualCurrencySourceClientSDK, VirtualCurrencySourceExpiration, VirtualCurrencySourceAdReward:
		return true
	default:
		return false
	}
}

// IsKnownRedemptionOutcome reports whether RevenueCat documents outcome.
// Unknown values remain decodable. See https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields.
func IsKnownRedemptionOutcome(outcome string) bool {
	switch outcome {
	case RedemptionOutcomeAlias, RedemptionOutcomeTransfer, RedemptionOutcomeRedeemerOwns:
		return true
	default:
		return false
	}
}
