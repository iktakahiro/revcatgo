package revcatgo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookFixtureRoundTrips(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "webhooks", "*.json"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			var original WebhookEvent
			require.NoError(t, json.Unmarshal(data, &original))
			encoded, err := json.Marshal(original)
			require.NoError(t, err)
			var decoded WebhookEvent
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestSubscriberAttributeRoundTripShapes(t *testing.T) {
	cases := []struct {
		name string
		json string
	}{
		{"object", `{"value":"customer@example.com","updated_at_ms":1784336523000}`},
		{"flattened", `"customer@example.com"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var original subscriberAttribute
			require.NoError(t, json.Unmarshal([]byte(tc.json), &original))
			encoded, err := json.Marshal(original)
			require.NoError(t, err)
			assert.JSONEq(t, tc.json, string(encoded))
			var decoded subscriberAttribute
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestStoreTransactionIDJSON(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		value   string
		encoded string
	}{
		{"empty", `""`, "", `null`},
		{"string", `"txn_123"`, "txn_123", `"txn_123"`},
		{"number", `1000000652379790`, "1000000652379790", `"1000000652379790"`},
		{"null", `null`, "", `null`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var id storeTransactionID
			require.NoError(t, json.Unmarshal([]byte(tc.input), &id))
			assert.Equal(t, tc.value, id.String())
			encoded, err := json.Marshal(id)
			require.NoError(t, err)
			assert.JSONEq(t, tc.encoded, string(encoded))
		})
	}
}

func TestInvalidJSONValues(t *testing.T) {
	cases := []struct {
		name  string
		input string
		value any
	}{
		{"price string", `"invalid"`, &price{}},
		{"price object", `{}`, &price{}},
		{"milliseconds string", `"invalid"`, &milliseconds{}},
		{"milliseconds object", `{}`, &milliseconds{}},
		{"milliseconds below threshold", `1605526568`, &milliseconds{}},
		{"subscriber transaction ID", `{"subscriptions":{"product":{"store_transaction_id":true}}}`, &Subscriber{}},
		{"subscriber attribute object", `{"subscriber_attributes":{"key":{"updated_at_ms":"invalid"}}}`, &Subscriber{}},
		{"subscriber attribute flattened", `{"subscriber_attributes":{"key":true}}`, &Subscriber{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.input), tc.value)
			require.Error(t, err)
		})
	}
}

func TestInvalidJSONValueErrorsWrapTypeErrors(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		value any
	}{
		{"price", `{}`, &price{}},
		{"milliseconds", `{}`, &milliseconds{}},
		{"subscriber transaction ID", `{"subscriptions":{"product":{"store_transaction_id":true}}}`, &Subscriber{}},
		{"subscriber attribute", `{"subscriber_attributes":{"key":true}}`, &Subscriber{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.input), tc.value)
			require.Error(t, err)
			var typeErr *json.UnmarshalTypeError
			require.ErrorAs(t, err, &typeErr)
		})
	}
}
