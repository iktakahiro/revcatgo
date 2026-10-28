package revcatgo

import (
	"encoding/json"
	"testing"
)

func TestKnownWebhookStringValues(t *testing.T) {
	tests := []struct {
		name   string
		known  func(string) bool
		values []string
	}{
		{"source", IsKnownVirtualCurrencySource, []string{VirtualCurrencySourceInAppPurchase, VirtualCurrencySourceAdminAPI, VirtualCurrencySourceServerAPI, VirtualCurrencySourceClientSDK, VirtualCurrencySourceExpiration, VirtualCurrencySourceAdReward}},
		{"outcome", IsKnownRedemptionOutcome, []string{RedemptionOutcomeAlias, RedemptionOutcomeTransfer, RedemptionOutcomeRedeemerOwns}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, value := range tt.values {
				if !tt.known(value) {
					t.Errorf("%q should be known", value)
				}
			}
			for _, value := range []string{"", "future_value", "TRANSFER"} {
				if tt.known(value) {
					t.Errorf("%q should be unknown", value)
				}
			}
		})
	}
}

// Payloads adapted from https://www.revenuecat.com/docs/integrations/webhooks/sample-events.
func TestWebhookStringEnumSamples(t *testing.T) {
	tests := []struct{ name, payload, source, outcome string }{
		{"currency", `{"api_version":"1.0","event":{"type":"VIRTUAL_CURRENCY_TRANSACTION","id":"12345678","app_user_id":"1234567890","source":"in_app_purchase","adjustments":[{"amount":100,"currency":{"code":"CRD","description":"The main currency unit","name":"Credits"}}],"purchase_environment":"PRODUCTION","store":"APP_STORE"}}`, VirtualCurrencySourceInAppPurchase, ""},
		{"redemption", `{"api_version":"1.0","event":{"type":"PURCHASE_REDEEMED","id":"CD489E0E","store":"RC_BILLING","environment":"PRODUCTION","redeemed_from":["anonymous"],"redeemed_by":["redeemer-app-user-id"],"redemption_outcome":"alias","redemption_platform":"ios","product_id":"premium_monthly","entitlement_ids":["premium"]}}`, "", RedemptionOutcomeAlias},
		{"future", `{"api_version":"1.0","event":{"type":"PURCHASE_REDEEMED","source":"future_source","redemption_outcome":"future_outcome"}}`, "future_source", "future_outcome"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload WebhookEvent
			if err := json.Unmarshal([]byte(tt.payload), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.Event.Source != tt.source || payload.Event.RedemptionOutcome != tt.outcome {
				t.Fatalf("unexpected values: source=%q outcome=%q", payload.Event.Source, payload.Event.RedemptionOutcome)
			}
		})
	}
}
