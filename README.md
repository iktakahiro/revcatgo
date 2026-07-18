# revcatgo

[![GoDev][godev-image]][godev-url]
![CI](https://github.com/iktakahiro/revcatgo/actions/workflows/ci.yml/badge.svg?branch=main)

revcatgo is a lightweight Go helper library that models [RevenueCat](https://www.revenuecat.com) webhook payloads and subscriber API responses with type-safe value objects and small convenience helpers.

## Features

- Strongly typed wrappers for event type, environment, store, and other enumerations to avoid typo-prone string comparisons.
- Utility methods for common webhook workflows such as checking expiration windows, collecting related user identifiers, and inspecting sandbox events.
- Data structures that mirror RevenueCat’s subscriber API responses, ready for direct decoding with the standard library.
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

## Working with subscriber API responses

The same structs can unmarshal the JSON returned by the Subscriber API. `null` values from RevenueCat are preserved through `gopkg.in/guregu/null.v4`.

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
  if entitlement.ExpiresDate.Before(time.Now()) {
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

## License

Distributed under the [MIT License](LICENSE).

[godev-image]: https://pkg.go.dev/badge/github.com/iktakahiro/revcatgo
[godev-url]: https://pkg.go.dev/github.com/iktakahiro/revcatgo
