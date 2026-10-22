package revcatgo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDecodeOfficialWebhookSamples decodes the sample payloads published at
// https://www.revenuecat.com/docs/integrations/webhooks/sample-events.
func TestDecodeOfficialWebhookSamples(t *testing.T) {
	expectedTypes := map[string]string{
		"billing_issue":                EventTypeBillingIssue,
		"cancellation":                 EventTypeCancellation,
		"experiment_enrollment":        EventTypeExperimentEnrollment,
		"expiration":                   EventTypeExpiration,
		"initial_purchase":             EventTypeInitialPurchase,
		"invoice_issuance":             EventTypeInvoiceIssuance,
		"non_renewing_purchase":        EventTypeNonRenewingPurchase,
		"product_change":               EventTypeProductChange,
		"purchase_redeemed":            EventTypePurchaseRedeemed,
		"refund":                       EventTypeCancellation,
		"refund_reversed":              EventTypeRefundReversed,
		"renewal":                      EventTypeRenewal,
		"subscription_extended":        EventTypeSubscriptionExtended,
		"subscription_paused":          EventTypeSubscriptionPaused,
		"temporary_entitlement_grant":  EventTypeTemporaryEntitlementGrant,
		"transfer":                     EventTypeTransfer,
		"trial_cancelled":              EventTypeCancellation,
		"trial_started":                EventTypeInitialPurchase,
		"uncancellation":               EventTypeUnCancellation,
		"virtual_currency_transaction": EventTypeVirtualCurrencyTransaction,
	}

	files, err := filepath.Glob(filepath.Join("testdata", "webhooks", "*.json"))
	require.NoError(t, err)
	require.Len(t, files, len(expectedTypes))

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".json")
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			var webhook WebhookEvent
			require.NoError(t, json.Unmarshal(data, &webhook))

			expectedType, ok := expectedTypes[name]
			require.True(t, ok, "missing expected event type for %s", name)
			assert.Equal(t, "1.0", webhook.APIVersion)
			assert.Equal(t, expectedType, webhook.Event.Type.String())
			assert.True(t, webhook.Event.Type.IsKnown())
		})
	}
}

func TestDecodeOfficialWebhookSampleDetails(t *testing.T) {
	load := func(t *testing.T, name string) Event {
		t.Helper()
		data, err := os.ReadFile(filepath.Join("testdata", "webhooks", name+".json"))
		require.NoError(t, err)

		var webhook WebhookEvent
		require.NoError(t, json.Unmarshal(data, &webhook))

		return webhook.Event
	}

	t.Run("invoice_issuance", func(t *testing.T) {
		event := load(t, "invoice_issuance")
		assert.Empty(t, event.PeriodType.String())
		assert.False(t, event.PeriodType.IsKnown())
		assert.Equal(t, StoreRCBilling, event.Store.String())
		assert.True(t, event.Store.IsKnown())
		assert.False(t, event.ExpirationAt.NullInt().Valid)
		assert.False(t, event.EntitlementID.Valid)
		assert.Empty(t, event.TransactionID)
		assert.Zero(t, event.RenewalNumber)
	})

	t.Run("cancellation", func(t *testing.T) {
		event := load(t, "cancellation")
		assert.Equal(t, "pro", event.EntitlementID.ValueOrZero())
		assert.Equal(t, CancelReasonUnsubscribe, event.CancelReason.String())
		assert.True(t, event.CancelReason.IsKnown())
	})

	t.Run("refund", func(t *testing.T) {
		event := load(t, "refund")
		assert.Equal(t, CancelReasonCustomerSupport, event.CancelReason.String())
		assert.True(t, event.CancelReason.IsKnown())
		assert.Negative(t, event.Price.Float64())
	})

	t.Run("expiration", func(t *testing.T) {
		event := load(t, "expiration")
		assert.Equal(t, CancelReasonUnsubscribe, event.ExpirationReason.String())
		assert.True(t, event.ExpirationReason.IsKnown())
		assert.False(t, event.CancelReason.NullString().Valid)
	})

	t.Run("temporary_entitlement_grant", func(t *testing.T) {
		event := load(t, "temporary_entitlement_grant")
		assert.Empty(t, event.Environment.String())
		assert.False(t, event.Environment.IsKnown())
		assert.False(t, event.Environment.IsProduction())
		assert.False(t, event.Environment.IsSandbox())
	})
}
