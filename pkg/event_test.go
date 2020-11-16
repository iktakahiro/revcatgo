package revcatgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	b := []byte(`{"store":"APP_STORE", "type":"INITIAL_PURCHASE"}`)

	var event Event
	err := json.Unmarshal(b, &event)
	assert.Nil(t, err)

	assert.Equal(t, "APP_STORE", event.Store.String())
	assert.Equal(t, "INITIAL_PURCHASE", event.Type.String())

	b = []byte(`{"store":1}`)
	err = json.Unmarshal(b, &event)
	assert.Error(t, err)

	b = []byte(`{"store":null}`)
	err = json.Unmarshal(b, &event)
	assert.EqualError(t, err, "store is a required field")
}

const initialPurchaseRawJSON = `
{
      "product_id": "my.subscription.sandbox",
      "event_timestamp_ms": 1605256336738,
      "original_app_user_id": "$RCAnonymousID:0000000000000000000000000000000b",
      "expiration_at_ms": 1605256730251,
      "presented_offering_id": "test",
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
      "currency": "JPY",
      "entitlement_id": null,
      "purchased_at_ms": 1605256330385,
      "original_transaction_id": "GPA.0000-4204-5621-00000",
      "entitlement_ids": [
        "premium"
      ],
      "price_in_purchased_currency": 550,
      "takehome_percentage": 0.7,
      "store": "PLAY_STORE",
      "price": 5.233,
      "transaction_id": "GPA.0000-4204-5621-00000",
      "period_type": "NORMAL",
      "id": "00A23FAE-0DB8-42E2-A8DC-00000BCDF0D6",
      "aliases": [
        "$RCAnonymousID:0000000000000000000000000000000b"
      ]
    }
`

func TestUnmarshalInitialPurchaseEvent(t *testing.T) {
	b := []byte(initialPurchaseRawJSON)

	var event Event
	err := json.Unmarshal(b, &event)
	assert.Nil(t, err)

	assert.Equal(t, "my.subscription.sandbox", event.ProductID)
	assert.Equal(t, "PLAY_STORE", event.Store.String())
	assert.Equal(t, "INITIAL_PURCHASE", event.Type.String())
	assert.Equal(t, "SANDBOX", event.Environment.String())
	assert.Equal(t, "$RCAnonymousID:0000000000000000000000000000000b", event.AppUserID)
	assert.Equal(t, "$RCAnonymousID:0000000000000000000000000000000b", event.GetUserID())
	assert.Equal(t, float32(550), event.PriceInPurchasedCurrency)
	assert.Equal(t, int64(1605256730251), event.ExpirationAt.Int64())
}
