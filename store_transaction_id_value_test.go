package revcatgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreTransactionIDUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{`"GPA.6801-7988-0152-76034..5"`, "GPA.6801-7988-0152-76034..5"},
		{`1000000652379790`, "1000000652379790"},
		{`"a42db3af39530cb82b17eaf9c6576393"`, "a42db3af39530cb82b17eaf9c6576393"},
		{`null`, ""},
	}

	for _, c := range cases {
		var s storeTransactionID
		err := json.Unmarshal([]byte(c.in), &s)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, s.String())
	}
}

// TestStoreTransactionIDPreservesRawInput guards against mutating the raw JSON
// buffer while quoting a numeric id, which would corrupt the bytes following
// the token for callers that retain the original response.
func TestStoreTransactionIDPreservesRawInput(t *testing.T) {
	type payload struct {
		ID        storeTransactionID `json:"store_transaction_id"`
		IsSandbox bool               `json:"is_sandbox"`
	}

	raw := []byte(`{"store_transaction_id":1000000652379790,"is_sandbox":true}`)
	original := string(raw)

	var p payload
	err := json.Unmarshal(raw, &p)
	assert.Nil(t, err)
	assert.Equal(t, "1000000652379790", p.ID.String())
	assert.True(t, p.IsSandbox)
	assert.Equal(t, original, string(raw))
}
