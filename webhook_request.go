package revcatgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MaxWebhookBodyBytes is the maximum request body size accepted by ParseWebhookRequest.
const MaxWebhookBodyBytes = 1 << 20

var (
	// ErrWebhookBodyTooLarge indicates that a webhook request exceeds MaxWebhookBodyBytes.
	ErrWebhookBodyTooLarge = errors.New("RevenueCat webhook body exceeds the size limit")
	// ErrInvalidWebhookPayload indicates that a webhook body could not be decoded as JSON.
	ErrInvalidWebhookPayload = errors.New("invalid RevenueCat webhook payload")
)

// ParseWebhookRequest reads, verifies, and decodes a RevenueCat webhook request,
// limiting its body to MaxWebhookBodyBytes. An empty secret returns
// ErrWebhookSignatureMissingSecret. Tolerance has the same meaning as in
// VerifyWebhookSignature.
// The caller remains responsible for closing the request body.
//
// See https://www.revenuecat.com/docs/integrations/webhooks and
// https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields.
func ParseWebhookRequest(r *http.Request, secret string, tolerance time.Duration) (*WebhookEvent, error) {
	if r == nil || r.Body == nil {
		return nil, errors.New("RevenueCat webhook request and body must not be nil")
	}

	if secret == "" {
		return nil, ErrWebhookSignatureMissingSecret
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, MaxWebhookBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read RevenueCat webhook body: %w", err)
	}
	if len(body) > MaxWebhookBodyBytes {
		return nil, ErrWebhookBodyTooLarge
	}
	if err := VerifyWebhookSignature(r.Header.Get(WebhookSignatureHeader), body, secret, tolerance); err != nil {
		return nil, fmt.Errorf("verify RevenueCat webhook signature: %w", err)
	}

	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidWebhookPayload, err)
	}
	return &event, nil
}
