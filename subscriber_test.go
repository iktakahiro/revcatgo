package revcatgo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalSubscriberResponse(t *testing.T) {
	data := []byte(`{
		"request_date":"2026-07-18T01:02:03Z",
		"request_date_ms":1784336523000,
		"subscriber":{
			"entitlements":{
				"pro":{
					"expires_date":null,
					"grace_period_expires_date":null,
					"product_identifier":"lifetime",
					"purchase_date":"2026-01-02T03:04:05Z"
				}
			},
			"first_seen":"2025-01-01T00:00:00Z",
			"last_seen":"2026-07-18T00:00:00Z",
			"management_url":"https://example.com/manage",
			"original_app_user_id":"user-123",
			"original_application_version":"1.0",
			"original_purchase_date":"2025-01-01T00:00:00Z",
			"subscriber_attributes":{
				"$email":{"value":"customer@example.com","updated_at_ms":1784336523000}
			},
			"other_purchases":{},
			"subscriptions":{
				"annual":{
					"auto_resume_date":null,
					"billing_issues_detected_at":"2026-07-17T00:00:00Z",
					"expires_date":"2027-01-01T00:00:00Z",
					"grace_period_expires_date":"2027-01-08T00:00:00Z",
					"is_sandbox":false,
					"original_purchase_date":"2025-01-01T00:00:00Z",
					"ownership_type":"PURCHASED",
					"period_type":"normal",
					"purchase_date":"2026-01-01T00:00:00Z",
					"refunded_at":null,
					"store":"rc_billing",
					"store_transaction_id":"txn_123",
					"unsubscribe_detected_at":null
				},
				"monthly":{
					"expires_date":"2026-08-01T00:00:00Z",
					"original_purchase_date":"2026-07-01T00:00:00Z",
					"period_type":"prepaid",
					"purchase_date":"2026-07-01T00:00:00Z",
					"store":"app_store",
					"store_transaction_id":1000000652379790
				}
			},
			"non_subscriptions":{}
		}
	}`)

	var response SubscriberResponse
	require.NoError(t, json.Unmarshal(data, &response))

	assert.Equal(t, time.Date(2026, time.July, 18, 1, 2, 3, 0, time.UTC), response.RequestDate)
	assert.Equal(t, int64(1784336523000), response.RequestDateAt.Int64())
	assert.Equal(t, "1.0", response.Subscriber.OriginalApplicationVersion.String)
	assert.True(t, response.Subscriber.OriginalPurchaseDate.Valid)
	assert.Equal(t, "customer@example.com", response.Subscriber.SubscriberAttributes["$email"].Value)
	assert.Equal(t, time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC), response.Subscriber.Entitlements["pro"].PurchaseDate)

	annual := response.Subscriber.Subscriptions["annual"]
	assert.Equal(t, "RC_BILLING", annual.Store.String())
	assert.Equal(t, PeriodTypeNormal, annual.PeriodType.String())
	assert.Equal(t, "PURCHASED", annual.OwnershipType)
	assert.True(t, annual.BillingIssueDetectedAt.Valid)
	assert.Equal(t, "txn_123", annual.StoreTransactionID.String())

	monthly := response.Subscriber.Subscriptions["monthly"]
	assert.Equal(t, PeriodTypePrepaid, monthly.PeriodType.String())
	assert.Equal(t, "1000000652379790", monthly.StoreTransactionID.String())
}
