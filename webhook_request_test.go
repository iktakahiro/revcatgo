package revcatgo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWebhookRequest(t *testing.T) {
	sample, err := os.ReadFile("testdata/webhooks/initial_purchase.json")
	require.NoError(t, err)
	sign := func(body string, timestamp int64) string {
		mac := hmac.New(sha256.New, []byte("secret"))
		_, _ = fmt.Fprintf(mac, "%d.%s", timestamp, body)
		return fmt.Sprintf("t=%d,v1=%x", timestamp, mac.Sum(nil))
	}
	signature := sign(string(sample), 1700000000)
	atLimit := `{"event":{"type":"FUTURE_EVENT"}}`
	atLimit += strings.Repeat(" ", MaxWebhookBodyBytes-len(atLimit))

	tests := []struct {
		name, body, secret, signature string
		wantErr                       error
		wantType                      string
		unknown                       bool
		tolerance                     time.Duration
	}{
		{name: "valid signed sample", body: string(sample), secret: "secret", signature: signature, wantType: EventTypeInitialPurchase},
		{name: "empty secret", body: string(sample), wantErr: ErrWebhookSignatureMissingSecret},
		{name: "bad signature", body: string(sample), secret: "wrong", signature: signature, wantErr: ErrWebhookSignatureMismatch},
		{name: "missing signature", body: string(sample), secret: "secret", wantErr: ErrWebhookSignatureMalformed},
		{name: "malformed JSON", body: `{`, secret: "secret", signature: sign(`{`, 1700000000), wantErr: ErrInvalidWebhookPayload},
		{name: "oversized", body: strings.Repeat(" ", MaxWebhookBodyBytes+1), secret: "secret", wantErr: ErrWebhookBodyTooLarge},
		{name: "at limit", body: atLimit, secret: "secret", signature: sign(atLimit, 1700000000), wantType: "FUTURE_EVENT", unknown: true},
		{name: "expired", body: string(sample), secret: "secret", signature: signature, tolerance: time.Minute, wantErr: ErrWebhookSignatureExpired},
		{name: "current", body: string(sample), secret: "secret", signature: sign(string(sample), time.Now().Unix()), tolerance: time.Hour, wantType: EventTypeInitialPurchase},
		{name: "negative tolerance", body: string(sample), secret: "secret", signature: signature, tolerance: -time.Second, wantErr: ErrWebhookSignatureInvalidTolerance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			r.Header.Set(WebhookSignatureHeader, tt.signature)
			got, err := ParseWebhookRequest(r, tt.secret, tt.tolerance)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantType, got.Event.Type.String())
			if tt.unknown {
				assert.False(t, got.Event.Type.IsKnown())
			}
		})
	}
}

func TestParseWebhookRequestInvalidInput(t *testing.T) {
	for _, r := range []*http.Request{nil, {}} {
		got, err := ParseWebhookRequest(r, "", 0)
		require.Error(t, err)
		assert.Nil(t, got)
	}
	readErr := errors.New("read failed")
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(nil))
	r.Body = io.NopCloser(errorReader{readErr})
	_, err := ParseWebhookRequest(r, "secret", 0)
	assert.ErrorIs(t, err, readErr)
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }
