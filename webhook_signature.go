package revcatgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// WebhookSignatureHeader is the HTTP header used for RevenueCat HMAC signatures.
const WebhookSignatureHeader = "X-RevenueCat-Webhook-Signature"

// VerifyWebhookSignature verifies a RevenueCat webhook HMAC signature against the raw request body.
// A positive tolerance rejects timestamps outside the allowed window. A zero tolerance disables
// timestamp validation while still verifying the signature.
func VerifyWebhookSignature(signatureHeader string, body []byte, secret string, tolerance time.Duration) error {
	return verifyWebhookSignatureAt(signatureHeader, body, secret, tolerance, time.Now())
}

func verifyWebhookSignatureAt(
	signatureHeader string,
	body []byte,
	secret string,
	tolerance time.Duration,
	now time.Time,
) error {
	if secret == "" {
		return errors.New("RevenueCat webhook signing secret must not be empty")
	}
	if tolerance < 0 {
		return errors.New("RevenueCat webhook signature tolerance must not be negative")
	}

	timestampText, timestamp, signature, err := parseWebhookSignature(signatureHeader)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestampText))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(body)
	if !hmac.Equal(mac.Sum(nil), signature) {
		return errors.New("RevenueCat webhook signature is invalid")
	}
	if tolerance > 0 {
		signedAt := time.Unix(timestamp, 0)
		if signedAt.Before(now.Add(-tolerance)) || signedAt.After(now.Add(tolerance)) {
			return errors.New("RevenueCat webhook signature timestamp is outside the tolerance")
		}
	}

	return nil
}

func parseWebhookSignature(signatureHeader string) (string, int64, []byte, error) {
	var timestampText string
	var signatureText string

	for part := range strings.SplitSeq(signatureHeader, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		switch key {
		case "t":
			timestampText = value
		case "v1":
			signatureText = value
		}
	}
	if timestampText == "" || signatureText == "" {
		return "", 0, nil, errors.New("RevenueCat webhook signature header must contain t and v1")
	}

	timestamp, err := strconv.ParseInt(timestampText, 10, 64)
	if err != nil {
		return "", 0, nil, fmt.Errorf("invalid RevenueCat webhook signature timestamp: %w", err)
	}
	signature, err := hex.DecodeString(signatureText)
	if err != nil || len(signature) != sha256.Size {
		return "", 0, nil, errors.New("RevenueCat webhook signature must be a SHA-256 hex digest")
	}

	return timestampText, timestamp, signature, nil
}
