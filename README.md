# revcatgo

[![GoDev][godev-image]][godev-url]
![CI](https://github.com/iktakahiro/revcatgo/actions/workflows/ci.yml/badge.svg?branch=main)

revcatgo is a lightweight Go helper library that models [RevenueCat](https://www.revenuecat.com) webhook payloads and subscriber API responses with type-safe value objects and small convenience helpers.

## Features

- Strongly typed wrappers for event type, environment, store, and other enumerations to avoid typo-prone string comparisons.
- Coverage for the current RevenueCat webhook families, including experiment enrollment, purchase redemption, Paywall UI, price increase consent, and virtual currency events.
- Forward-compatible decoding for custom Paywall event names and new RevenueCat event types.
- HMAC-SHA256 webhook signature verification with optional replay protection.
- Utility methods for common webhook workflows such as checking expiration windows, collecting related user identifiers, and inspecting sandbox events.
- Data structures that mirror RevenueCat’s Subscriber API v1 responses, ready for direct decoding with the standard library.
- Works with Go 1.26+ and has a small dependency footprint (only `gopkg.in/guregu/null.v4` for nullable values).

## Installation

```bash
go get github.com/iktakahiro/revcatgo@latest
```

## Quick start

Decode a webhook payload inside your HTTP handler and branch on the event type or environment using the provided helpers:

```go
package webhook

import (
 "encoding/json"
 "log"
 "net/http"
 "time"

 "github.com/iktakahiro/revcatgo"
)

func HandleRevenueCatWebhook(w http.ResponseWriter, r *http.Request) {
 var hook revcatgo.WebhookEvent
 if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 event := hook.Event

 if !event.Environment.IsProduction() {
  log.Printf("ignoring %s event from sandbox", event.Environment.String())
  w.WriteHeader(http.StatusNoContent)
  return
 }

 switch event.Type.String() {
 case revcatgo.EventTypeInitialPurchase:
  log.Printf("new subscriber %q", event.AppUserID)
 case revcatgo.EventTypeCancellation:
  if event.IsExpired(24*time.Hour, nil) {
   log.Printf("subscription cancelled and grace period elapsed for %q", event.AppUserID)
  }
 default:
  log.Printf("received %s for %q", event.Type.String(), event.AppUserID)
 }

 allIDs := event.GetAllRelatedUserID() // includes aliases and original user ids
 log.Printf("related identifiers: %v", allIDs)

 w.WriteHeader(http.StatusNoContent)
}
```

## Verifying webhook signatures

RevenueCat signs the exact request bytes, so verify the raw body before unmarshalling it. A positive tolerance also rejects replayed requests whose timestamp is outside the allowed window.

```go
body, err := io.ReadAll(r.Body)
if err != nil {
 http.Error(w, err.Error(), http.StatusBadRequest)
 return
}

err = revcatgo.VerifyWebhookSignature(
 r.Header.Get(revcatgo.WebhookSignatureHeader),
 body,
 webhookSigningSecret,
 5*time.Minute,
)
if err != nil {
 http.Error(w, "invalid RevenueCat signature", http.StatusUnauthorized)
 return
}

var hook revcatgo.WebhookEvent
if err := json.Unmarshal(body, &hook); err != nil {
 http.Error(w, err.Error(), http.StatusBadRequest)
 return
}
```

## Working with Subscriber API v1 responses

The same structs can unmarshal the JSON returned by `GET /v1/subscribers/{app_user_id}`. `null` values from RevenueCat are preserved through `gopkg.in/guregu/null.v4`.

```go
import (
 "encoding/json"
 "io"
 "time"

 "github.com/iktakahiro/revcatgo"
)

func decodeSubscriberResponse(body io.Reader) (*revcatgo.SubscriberResponse, error) {
 var resp revcatgo.SubscriberResponse
 if err := json.NewDecoder(body).Decode(&resp); err != nil {
  return nil, err
 }

 for productID, entitlement := range resp.Subscriber.Entitlements {
  // RevenueCat returns null for lifetime entitlement expiration dates.
  if !entitlement.ExpiresDate.IsZero() && entitlement.ExpiresDate.Before(time.Now()) {
   continue
  }
  // grant access for productID
 }

 return &resp, nil
}
```

## Development

Common helper targets:

```bash
# Install the pinned Go toolchain and golangci-lint binary
mise install

# Format Go sources via golangci-lint fmt (requires golangci-lint v2+)
make fmt

# Run the configured lint suite
make lint

# Execute unit tests
make test

# Scan dependencies for known vulnerabilities
make vulncheck
```

## References

- RevenueCat webhooks documentation: <https://www.revenuecat.com/docs/integrations/webhooks>
- RevenueCat webhook event types and fields: <https://www.revenuecat.com/docs/integrations/webhooks/event-types-and-fields>
- RevenueCat API v1 reference: <https://www.revenuecat.com/docs/api-v1>

## License

Distributed under the [MIT License](LICENSE).

[godev-image]: https://pkg.go.dev/badge/github.com/iktakahiro/revcatgo
[godev-url]: https://pkg.go.dev/github.com/iktakahiro/revcatgo
