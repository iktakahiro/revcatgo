package revcatgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

// WebhookSignatureHeader is the HTTP header used for RevenueCat HMAC signatures.
const WebhookSignatureHeader = "X-RevenueCat-Webhook-Signature"

// Sentinel errors returned by VerifyWebhookSignature. Use errors.Is to tell the failure modes apart,
// for example to answer 400 for ErrWebhookSignatureMalformed and 401 for ErrWebhookSignatureMismatch.
var (
	// ErrWebhookSignatureMissingSecret reports that the signing secret passed to
	// VerifyWebhookSignature is empty. It indicates a caller configuration problem.
	ErrWebhookSignatureMissingSecret = errors.New("RevenueCat webhook signing secret must not be empty")

	// ErrWebhookSignatureInvalidTolerance reports that the tolerance passed to
	// VerifyWebhookSignature is negative. It indicates a caller configuration problem.
	ErrWebhookSignatureInvalidTolerance = errors.New("RevenueCat webhook signature tolerance must not be negative")

	// ErrWebhookSignatureMalformed reports that the signature header could not be parsed: the t or v1
	// element is missing, the timestamp is not an integer, or the signature is not a SHA-256 hex digest.
	ErrWebhookSignatureMalformed = errors.New("RevenueCat webhook signature header is malformed")

	// ErrWebhookSignatureMismatch reports that the signature does not match the request body and
	// signing secret. The timestamp is not examined once this error is returned.
	ErrWebhookSignatureMismatch = errors.New("RevenueCat webhook signature is invalid")

	// ErrWebhookSignatureExpired reports that the signature is valid but its timestamp is outside the
	// allowed tolerance, which can indicate a replayed request.
	ErrWebhookSignatureExpired = errors.New("RevenueCat webhook signature timestamp is outside the tolerance")
)

// malformedWebhookSignatureError carries a specific message and optional cause while matching
// ErrWebhookSignatureMalformed through errors.Is.
type malformedWebhookSignatureError struct {
	message string
	cause   error
}

func (e *malformedWebhookSignatureError) Error() string {
	return e.message
}

func (e *malformedWebhookSignatureError) Unwrap() error {
	return e.cause
}

func (e *malformedWebhookSignatureError) Is(target error) bool {
	return target == ErrWebhookSignatureMalformed
}

// VerifyWebhookSignature verifies a RevenueCat webhook HMAC signature against the raw request body.
// A positive tolerance rejects timestamps outside the allowed window. A zero tolerance disables
// timestamp validation while still verifying the signature.
//
// The signature is checked before the timestamp. Every non-nil error matches one of the following
// sentinels via errors.Is:
//
//   - ErrWebhookSignatureMissingSecret: secret is empty.
//   - ErrWebhookSignatureInvalidTolerance: tolerance is negative.
//   - ErrWebhookSignatureMalformed: the header is missing t or v1, or contains an unparsable
//     timestamp or digest.
//   - ErrWebhookSignatureMismatch: the signature does not match the body and secret.
//   - ErrWebhookSignatureExpired: the signature is valid but the timestamp is outside tolerance.
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
		return ErrWebhookSignatureMissingSecret
	}
	if tolerance < 0 {
		return ErrWebhookSignatureInvalidTolerance
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
		return ErrWebhookSignatureMismatch
	}
	if tolerance > 0 {
		signedAt := time.Unix(timestamp, 0)
		if signedAt.Before(now.Add(-tolerance)) || signedAt.After(now.Add(tolerance)) {
			return ErrWebhookSignatureExpired
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
		return "", 0, nil, &malformedWebhookSignatureError{
			message: "RevenueCat webhook signature header must contain t and v1",
		}
	}

	timestamp, err := strconv.ParseInt(timestampText, 10, 64)
	if err != nil {
		return "", 0, nil, &malformedWebhookSignatureError{
			message: "invalid RevenueCat webhook signature timestamp: " + err.Error(),
			cause:   err,
		}
	}
	signature, err := hex.DecodeString(signatureText)
	if err != nil || len(signature) != sha256.Size {
		return "", 0, nil, &malformedWebhookSignatureError{
			message: "RevenueCat webhook signature must be a SHA-256 hex digest",
			cause:   err,
		}
	}

	return timestampText, timestamp, signature, nil
}
