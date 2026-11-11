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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWebhookRequest(t *testing.T) {
	sample, err := os.ReadFile("testdata/webhooks/initial_purchase.json")
	require.NoError(t, err)
	mac := hmac.New(sha256.New, []byte("secret"))
	_, _ = mac.Write([]byte("1700000000."))
	_, _ = mac.Write(sample)
	signature := fmt.Sprintf("t=1700000000,v1=%x", mac.Sum(nil))

	tests := []struct {
		name, body, secret, signature string
		wantErr                       error
		wantType                      string
		unknown                       bool
	}{
		{name: "valid signed sample", body: string(sample), secret: "secret", signature: signature, wantType: EventTypeInitialPurchase},
		{name: "empty secret", body: string(sample), wantType: EventTypeInitialPurchase},
		{name: "bad signature", body: string(sample), secret: "wrong", signature: signature, wantErr: ErrWebhookSignatureMismatch},
		{name: "missing signature", body: string(sample), secret: "secret", wantErr: ErrWebhookSignatureMalformed},
		{name: "malformed JSON", body: `{`, wantErr: ErrInvalidWebhookPayload},
		{name: "oversized", body: strings.Repeat(" ", MaxWebhookBodyBytes+1), wantErr: ErrWebhookBodyTooLarge},
		{name: "at limit", body: `{"event":{"type":"FUTURE_EVENT"}}` + strings.Repeat(" ", MaxWebhookBodyBytes-len(`{"event":{"type":"FUTURE_EVENT"}}`)), wantType: "FUTURE_EVENT", unknown: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			r.Header.Set(WebhookSignatureHeader, tt.signature)
			got, err := ParseWebhookRequest(r, tt.secret, 0)
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
	_, err := ParseWebhookRequest(r, "", 0)
	assert.ErrorIs(t, err, readErr)
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }
