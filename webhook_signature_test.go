package revcatgo

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testWebhookBody      = `{"api_version":"1.0","event":{"id":"evt_123"}}`
	testWebhookSignature = "t=1700000000,v1=450a93f016917ae1cb87a06c780abccb13f8ee581865151b824759080a87da91"
)

func TestVerifyWebhookSignature(t *testing.T) {
	now := time.Unix(1700000000, 0)

	assert.NoError(t, verifyWebhookSignatureAt(
		testWebhookSignature,
		[]byte(testWebhookBody),
		"whsec_test",
		5*time.Minute,
		now,
	))
	assert.NoError(t, VerifyWebhookSignature(
		testWebhookSignature,
		[]byte(testWebhookBody),
		"whsec_test",
		0,
	))
}

func TestVerifyWebhookSignatureRejectsInvalidRequests(t *testing.T) {
	now := time.Unix(1700000000, 0)
	tests := []struct {
		name      string
		header    string
		body      string
		secret    string
		tolerance time.Duration
		wantError string
		wantIs    error
	}{
		{
			name:      "modified body",
			header:    testWebhookSignature,
			body:      testWebhookBody + " ",
			secret:    "whsec_test",
			tolerance: 5 * time.Minute,
			wantError: "signature is invalid",
			wantIs:    ErrWebhookSignatureMismatch,
		},
		{
			name:      "stale timestamp",
			header:    testWebhookSignature,
			body:      testWebhookBody,
			secret:    "whsec_test",
			tolerance: time.Minute,
			wantError: "outside the tolerance",
			wantIs:    ErrWebhookSignatureExpired,
		},
		{
			name:      "missing version",
			header:    "t=1700000000",
			body:      testWebhookBody,
			secret:    "whsec_test",
			tolerance: 0,
			wantError: "must contain t and v1",
			wantIs:    ErrWebhookSignatureMalformed,
		},
		{
			name:      "invalid timestamp",
			header:    "t=soon,v1=450a93f016917ae1cb87a06c780abccb13f8ee581865151b824759080a87da91",
			body:      testWebhookBody,
			secret:    "whsec_test",
			tolerance: 0,
			wantError: "invalid RevenueCat webhook signature timestamp",
			wantIs:    ErrWebhookSignatureMalformed,
		},
		{
			name:      "invalid digest",
			header:    "t=1700000000,v1=not-hex",
			body:      testWebhookBody,
			secret:    "whsec_test",
			tolerance: 0,
			wantError: "must be a SHA-256 hex digest",
			wantIs:    ErrWebhookSignatureMalformed,
		},
		{
			name:      "empty secret",
			header:    testWebhookSignature,
			body:      testWebhookBody,
			secret:    "",
			tolerance: 0,
			wantError: "secret must not be empty",
			wantIs:    ErrWebhookSignatureMissingSecret,
		},
		{
			name:      "negative tolerance",
			header:    testWebhookSignature,
			body:      testWebhookBody,
			secret:    "whsec_test",
			tolerance: -time.Second,
			wantError: "tolerance must not be negative",
			wantIs:    ErrWebhookSignatureInvalidTolerance,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testNow := now
			if test.name == "stale timestamp" {
				testNow = now.Add(2 * time.Minute)
			}
			err := verifyWebhookSignatureAt(
				test.header,
				[]byte(test.body),
				test.secret,
				test.tolerance,
				testNow,
			)
			require.ErrorContains(t, err, test.wantError)
			require.ErrorIs(t, err, test.wantIs)
			for _, other := range allWebhookSignatureErrors {
				if !errors.Is(other, test.wantIs) {
					assert.NotErrorIs(t, err, other)
				}
			}
		})
	}
}

var allWebhookSignatureErrors = []error{
	ErrWebhookSignatureMissingSecret,
	ErrWebhookSignatureInvalidTolerance,
	ErrWebhookSignatureMalformed,
	ErrWebhookSignatureMismatch,
	ErrWebhookSignatureExpired,
}

func TestVerifyWebhookSignatureMalformedTimestampWrapsCause(t *testing.T) {
	err := verifyWebhookSignatureAt(
		"t=soon,v1=450a93f016917ae1cb87a06c780abccb13f8ee581865151b824759080a87da91",
		[]byte(testWebhookBody),
		"whsec_test",
		0,
		time.Unix(1700000000, 0),
	)

	var numErr *strconv.NumError
	require.ErrorIs(t, err, ErrWebhookSignatureMalformed)
	require.ErrorAs(t, err, &numErr)
	assert.ErrorIs(t, err, strconv.ErrSyntax)
}

func TestVerifyWebhookSignatureChecksSignatureBeforeTimestamp(t *testing.T) {
	// A stale timestamp with a wrong body must report a mismatch, not expiry.
	err := verifyWebhookSignatureAt(
		testWebhookSignature,
		[]byte(testWebhookBody+" "),
		"whsec_test",
		time.Minute,
		time.Unix(1700000000, 0).Add(time.Hour),
	)

	require.ErrorIs(t, err, ErrWebhookSignatureMismatch)
	assert.NotErrorIs(t, err, ErrWebhookSignatureExpired)
}

func TestVerifyWebhookSignatureExpiredFutureTimestamp(t *testing.T) {
	err := verifyWebhookSignatureAt(
		testWebhookSignature,
		[]byte(testWebhookBody),
		"whsec_test",
		time.Minute,
		time.Unix(1700000000, 0).Add(-time.Hour),
	)

	assert.ErrorIs(t, err, ErrWebhookSignatureExpired)
}
