package revcatgo_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/iktakahiro/revcatgo"
)

func ExampleWebhookEvent() {
	body := []byte(`{
		"api_version": "1.0",
		"event": {
			"type": "INITIAL_PURCHASE",
			"environment": "PRODUCTION",
			"app_user_id": "user_123"
		}
	}`)
	var webhook revcatgo.WebhookEvent
	if err := json.Unmarshal(body, &webhook); err != nil {
		panic(err)
	}

	fmt.Println(webhook.Event.Type.String(), webhook.Event.Environment.String(), webhook.Event.AppUserID)
	// Output: INITIAL_PURCHASE PRODUCTION user_123
}

func ExampleVerifyWebhookSignature() {
	body := []byte(`{"api_version":"1.0","event":{"id":"evt_123"}}`)
	secret := "webhook-secret"
	timestamp := "1700000000"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)
	header := fmt.Sprintf("t=%s,v1=%x", timestamp, mac.Sum(nil))

	err := revcatgo.VerifyWebhookSignature(header, body, secret, 0)
	fmt.Println(err)
	// Output: <nil>
}

func ExampleSubscriberResponse() {
	body := []byte(`{
		"request_date": "2024-01-01T00:00:00Z",
		"subscriber": {
			"original_app_user_id": "user_123",
			"entitlements": {
				"premium": {
					"product_identifier": "premium_monthly",
					"purchase_date": "2024-01-01T00:00:00Z",
					"expires_date": "2024-02-01T00:00:00Z"
				}
			}
		}
	}`)
	var response revcatgo.SubscriberResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	fmt.Println(response.Subscriber.Entitlements["premium"].ProductIdentifier)
	// Output: premium_monthly
}
