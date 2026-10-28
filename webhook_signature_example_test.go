package revcatgo_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/iktakahiro/revcatgo"
)

func ExampleVerifyWebhookSignature_expired() {
	const (
		body      = `{"api_version":"1.0","event":{"id":"evt_123"}}`
		signature = "t=1700000000,v1=450a93f016917ae1cb87a06c780abccb13f8ee581865151b824759080a87da91"
	)

	// The signature is valid, but its timestamp is far older than the one-minute tolerance.
	err := revcatgo.VerifyWebhookSignature(signature, []byte(body), "whsec_test", time.Minute)
	if errors.Is(err, revcatgo.ErrWebhookSignatureExpired) {
		fmt.Println("rejected: possible replay")
	}
	// Output: rejected: possible replay
}
