package revcatgo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestNewEvent(t *testing.T) {
	b := []byte(`{"store":"APP_STORE", "type":"INITIAL_PURCHASE"}`)

	var event Event
	err := json.Unmarshal(b, &event)
	require.NoError(t, err)

	assert.Equal(t, "APP_STORE", event.Store.String())
	assert.Equal(t, "INITIAL_PURCHASE", event.Type.String())

	b = []byte(`{"store":1}`)
	err = json.Unmarshal(b, &event)
	require.Error(t, err)

	b = []byte(`{"store":null}`)
	err = json.Unmarshal(b, &event)
	assert.EqualError(t, err, "store is a required field")
}

func TestEventIsExpired(t *testing.T) {
	// 2020-12-10:10:00
	expirationAt, _ := newMilliseconds(null.IntFrom(1607595010000))
	event := &Event{ExpirationAt: *expirationAt}

	// 2020-12-15:10:00
	futureTime := time.Unix(1608027010000/1000, 0)
	// 2020-12-05:10:00
	pastTime := time.Unix(1607163010000/1000, 0)

	assert.True(t, event.IsExpired(0, &futureTime))
	assert.False(t, event.IsExpired(0, &pastTime))
	assert.False(t, event.IsExpired(time.Hour*168, &futureTime))
	assert.True(t, event.IsExpired(time.Hour*48, &futureTime))
	assert.False(t, (&Event{}).IsExpired(0, &futureTime))
}

const initialPurchaseRawJSON = `
{
   "product_id": "my.subscription.sandbox",
   "app_id": "app-123",
   "event_timestamp_ms": 1605256336738,
   "original_app_user_id": "$RCAnonymousID:0000000000000000000000000000000b",
   "expiration_at_ms": 1605256730251,
   "presented_offering_id": "default",
   "environment": "SANDBOX",
   "app_user_id": "$RCAnonymousID:0000000000000000000000000000000b",
   "subscriber_attributes": {
     "$gpsAdId": {
       "updated_at_ms": 1605236044746,
       "value": "0000-0b0a-4d05-b3fc-000000000000"
     },
     "$ip": {
       "value": "127.0.0.1",
       "updated_at_ms": 1605236044746
     }
   },
   "type": "INITIAL_PURCHASE",
   "currency": "USD",
   "entitlement_id": null,
   "purchased_at_ms": 1605256330385,
   "original_transaction_id": "GPA.0000-4204-5621-00000",
   "entitlement_ids": [
     "premium"
   ],
   "price_in_purchased_currency": 550,
   "tax_percentage": 5.5,
   "commission_percentage": 30.0,
   "takehome_percentage": 0.7,
   "store": "PLAY_STORE",
   "price": 5.233,
   "transaction_id": "GPA.0000-4204-5621-00000",
   "period_type": "NORMAL",
   "id": "00A23FAE-0DB8-42E2-A8DC-00000BCDF0D6",
   "aliases": [
     "$RCAnonymousID:0000000000000000000000000000000b"
   ],
   "experiments": [
     {
       "experiment_id": "exp_a",
       "experiment_variant": "treatment",
       "enrolled_at_ms": 1605256000000
     }
   ],
   "renewal_number": 1,
   "metadata": {"campaign": "summer"},
   "discount_percentage": 10.5,
   "discount_amount": 1.25,
   "discount_identifier": "launch-offer",
   "quantity": 2
}
`

func TestUnmarshalInitialPurchaseEvent(t *testing.T) {
	b := []byte(initialPurchaseRawJSON)

	var event Event
	err := json.Unmarshal(b, &event)
	require.NoError(t, err)

	assert.Equal(t, "my.subscription.sandbox", event.ProductID)
	assert.Equal(t, "PLAY_STORE", event.Store.String())
	assert.Equal(t, "INITIAL_PURCHASE", event.Type.String())
	assert.Equal(t, "SANDBOX", event.Environment.String())
	assert.Equal(t, "$RCAnonymousID:0000000000000000000000000000000b", event.AppUserID)
	assert.InDelta(t, 550, event.PriceInPurchasedCurrency, 0.001)
	assert.InDelta(t, 5.5, event.TaxPercentage, 0.001)
	assert.Equal(t, "app-123", event.AppID)
	assert.Len(t, event.Experiments, 1)
	assert.Equal(t, "exp_a", event.Experiments[0].ID)
	assert.Equal(t, int64(1605256000000), event.Experiments[0].EnrolledAt.Int64())
	assert.Equal(t, "summer", event.Metadata["campaign"])
	assert.Equal(t, 10.5, event.DiscountPercentage.Float64)
	assert.Equal(t, int64(2), event.Quantity.Int64)
	assert.Equal(t, int64(1605256730251), event.ExpirationAt.Int64())
	assert.True(t, event.HasEntitlementID("premium"))
	assert.False(t, event.HasEntitlementID("invalid_entitlement_id"))
}

const cancellationRawJSON = `
{
   "price": 0,
   "entitlement_ids": [
     "premium"
   ],
   "currency": "USD",
   "store": "PLAY_STORE",
   "product_id": "my.subscription.sandbox",
   "app_user_id": "$RCAnonymousID:0000000000000000000000000000000b",
   "original_transaction_id": "GPA.0000-4204-5621-00000",
   "type": "CANCELLATION",
   "presented_offering_id": "default",
   "purchased_at_ms": 1605258350251,
   "original_app_user_id": "$RCAnonymousID:0000000000000000000000000000000b",
   "id": "00A23FAE-0DB8-42E2-A8DC-00000BCDF0D6",
   "environment": "SANDBOX",
   "transaction_id": "GPA.0000-4204-5621-00000",
   "period_type": "NORMAL",
   "price_in_purchased_currency": 0,
   "subscriber_attributes": {
     "$gpsAdId": {
       "updated_at_ms": 1605236044746,
       "value": "0000-0b0a-4d05-b3fc-000000000000"
     },
     "$ip": {
       "value": "127.0.0.1",
       "updated_at_ms": 1605236044746
     }
   },
   "aliases": [
     "$RCAnonymousID:0000000000000000000000000000000b"
   ],
   "cancel_reason": "BILLING_ERROR",
   "entitlement_id": null,
   "event_timestamp_ms": 1605258535218,
   "takehome_percentage": 0.7,
   "expiration_at_ms": 1605258533373
}
`

func TestUnmarshalCancellationEvent(t *testing.T) {
	b := []byte(cancellationRawJSON)

	var event Event
	err := json.Unmarshal(b, &event)
	require.NoError(t, err)

	assert.Equal(t, "BILLING_ERROR", event.CancelReason.String())
}

const virtualCurrencyTransactionJSON = `
{
  "type": "VIRTUAL_CURRENCY_TRANSACTION",
  "app_user_id": "user123",
  "app_id": "app-virtual",
  "product_id": "coins.pack",
  "purchase_environment": "PRODUCTION",
  "source": "in_app_purchase",
  "virtual_currency_transaction_id": "vtxn_001",
  "updated_balance": 80,
  "adjustments": [
    {
      "amount": 100,
      "currency": {
        "code": "coins",
        "name": "Coins",
        "description": "In-app coins"
      }
    },
    {
      "amount": -20,
      "currency": {
        "code": "coins",
        "name": "Coins",
        "description": "In-app coins"
      }
    }
  ],
  "product_display_name": "Extra Coins Pack"
}
`

func TestUnmarshalVirtualCurrencyTransactionEvent(t *testing.T) {
	b := []byte(virtualCurrencyTransactionJSON)

	var event Event
	err := json.Unmarshal(b, &event)
	require.NoError(t, err)

	assert.Equal(t, "VIRTUAL_CURRENCY_TRANSACTION", event.Type.String())
	assert.Equal(t, "app-virtual", event.AppID)
	assert.Equal(t, "vtxn_001", event.VirtualTransactionID)
	assert.Equal(t, "Extra Coins Pack", event.ProductDisplayName)
	assert.Equal(t, "PRODUCTION", event.PurchaseEnvironment.String())
	assert.Equal(t, "in_app_purchase", event.Source)
	assert.Equal(t, int64(80), event.UpdatedBalance.Int64)
	assert.Len(t, event.Adjustments, 2)
	assert.Equal(t, 100, event.Adjustments[0].Amount)
	assert.Equal(t, "coins", event.Adjustments[0].Currency.Code)
}

func TestUnmarshalExperimentEnrollmentEvent(t *testing.T) {
	data := []byte(`{
		"type":"EXPERIMENT_ENROLLMENT",
		"id":"evt-experiment",
		"app_user_id":"user-123",
		"experiment_id":"prexpca1234abcd",
		"experiment_variant":"b",
		"offering_id":"experiment_offering_b",
		"experiment_enrolled_at_ms":1658726378679
	}`)

	var event Event
	assert.NoError(t, json.Unmarshal(data, &event))
	assert.Equal(t, EventTypeExperimentEnrollment, event.Type.String())
	assert.Equal(t, "prexpca1234abcd", event.ExperimentID)
	assert.Equal(t, "b", event.ExperimentVariant)
	assert.Equal(t, "experiment_offering_b", event.OfferingID)
	assert.Equal(t, int64(1658726378679), event.ExperimentEnrolledAt.Int64())
}

func TestUnmarshalPurchaseRedeemedEvent(t *testing.T) {
	data := []byte(`{
		"type":"PURCHASE_REDEEMED",
		"store":"RC_BILLING",
		"environment":"PRODUCTION",
		"redeemed_from":["web-user"],
		"redeemed_by":["app-user"],
		"redemption_outcome":"alias",
		"redemption_platform":"ios",
		"workflow_id":"wf_abc123",
		"workflow_step_id":"step_xyz789",
		"trace_id":"trace_abcdef"
	}`)

	var event Event
	assert.NoError(t, json.Unmarshal(data, &event))
	assert.Equal(t, EventTypePurchaseRedeemed, event.Type.String())
	assert.Equal(t, "RC_BILLING", event.Store.String())
	assert.Equal(t, []string{"web-user"}, event.RedeemedFrom)
	assert.Equal(t, []string{"app-user"}, event.RedeemedBy)
	assert.Equal(t, "alias", event.RedemptionOutcome)
	assert.Equal(t, "ios", event.RedemptionPlatform.String)
	assert.Equal(t, "trace_abcdef", event.TraceID)
}

func TestUnmarshalPaywallComponentEvent(t *testing.T) {
	data := []byte(`{
		"type":"PAYWALL_COMPONENT_INTERACTED",
		"event_id":"paywall-event-123",
		"app_user_id":"app-user",
		"platform":"iOS",
		"platform_version":"18.5",
		"sdk_version":"5.80.0",
		"subscriber_attributes":{"$email":"customer@example.com"},
		"paywall_id":"pw_123",
		"paywall_name":"Premium",
		"offering_id":"monthly",
		"session_id":"session-123",
		"display_mode":null,
		"dark_mode":false,
		"locale":"en_IE",
		"environment":"SANDBOX",
		"component_type":"package",
		"component_value":"$rc_monthly",
		"origin_index":0,
		"destination_index":1,
		"destination_package_id":"$rc_annual",
		"destination_product_id":"annual_product"
	}`)

	var event Event
	assert.NoError(t, json.Unmarshal(data, &event))
	assert.Equal(t, EventTypePaywallComponentInteracted, event.Type.String())
	assert.Equal(t, "paywall-event-123", event.PaywallEventID)
	assert.Equal(t, "pw_123", event.PaywallID)
	assert.False(t, event.DisplayMode.Valid)
	assert.True(t, event.DarkMode.Valid)
	assert.False(t, event.DarkMode.Bool)
	assert.Equal(t, "customer@example.com", event.SubscriberAttributes["$email"].Value)
	assert.Equal(t, "package", event.ComponentType)
	assert.Equal(t, int64(1), event.DestinationIndex.Int64)
	assert.Equal(t, "$rc_annual", event.DestinationPackageID)

	roundTrip, err := json.Marshal(event)
	assert.NoError(t, err)
	assert.Contains(t, string(roundTrip), `"$email":"customer@example.com"`)
}

func TestEvent_GetAllRelatedUserID(t *testing.T) {
	e := &Event{AppUserID: "one", OriginalAppUserID: "two", Aliases: []string{"one", "two"}}
	assert.Equal(t, []string{"one", "two"}, e.GetAllRelatedUserID())

	e = &Event{AppUserID: "one", OriginalAppUserID: "one", Aliases: []string{"one"}}
	assert.Equal(t, []string{"one"}, e.GetAllRelatedUserID())

	e = &Event{AppUserID: "one"}
	assert.Equal(t, []string{"one"}, e.GetAllRelatedUserID())
}
