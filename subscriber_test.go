package revcatgo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalEntitlementPurchaseDate(t *testing.T) {
	const payload = `{"purchase_date":"2026-07-18T12:34:56Z"}`

	var entitlement Entitlement
	require.NoError(t, json.Unmarshal([]byte(payload), &entitlement))

	assert.Equal(t, time.Date(2026, time.July, 18, 12, 34, 56, 0, time.UTC), entitlement.PurchaseDate)
}
