package revcatgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// subscriberRawJSON is trimmed from a RevenueCat v1 GET /subscribers response.
// It keeps one Play Store subscription (string store_transaction_id) and one
// App Store subscription (numeric store_transaction_id) to cover both shapes.
const subscriberRawJSON = `
{
  "request_date_ms": 1564162810884,
  "subscriber": {
    "entitlements": {},
    "first_seen": "2019-02-21T00:08:41Z",
    "last_seen": "2019-07-26T17:40:10Z",
    "original_app_user_id": "XXX-XXXXX-XXXXX-XX",
    "non_subscriptions": {},
    "subscriptions": {
      "annual": {
        "expires_date": "2019-08-14T21:07:40Z",
        "original_purchase_date": "2019-02-21T00:42:05Z",
        "period_type": "normal",
        "purchase_date": "2019-07-14T20:07:40Z",
        "store": "play_store",
        "store_transaction_id": "GPA.6801-7988-0152-76034..5",
        "is_sandbox": true
      },
      "onemonth": {
        "expires_date": "2019-06-17T22:47:55Z",
        "original_purchase_date": "2019-02-21T00:42:05Z",
        "period_type": "normal",
        "purchase_date": "2019-06-17T22:42:55Z",
        "store": "app_store",
        "store_transaction_id": 1000000652379790,
        "is_sandbox": true
      }
    }
  }
}
`

func TestUnmarshalSubscriberResponse(t *testing.T) {
	var resp SubscriberResponse
	err := json.Unmarshal([]byte(subscriberRawJSON), &resp)
	assert.Nil(t, err)

	annual := resp.Subscriber.Subscriptions["annual"]
	assert.Equal(t, "PLAY_STORE", annual.Store.String())
	assert.Equal(t, "GPA.6801-7988-0152-76034..5", annual.StoreTransactionID.String())

	onemonth := resp.Subscriber.Subscriptions["onemonth"]
	assert.Equal(t, "APP_STORE", onemonth.Store.String())
	assert.Equal(t, "1000000652379790", onemonth.StoreTransactionID.String())
}
