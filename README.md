# revcatgo

[![GoDev][godev-image]][godev-url]

![Run test](https://github.com/iktakahiro/revcatgo/workflows/Run%20test/badge.svg?branch=main)

A helper library for integrating the server side apps and [RevenueCat](https://www.revenuecat.com) webhook service

```bash
go get "github.com/iktakahiro/revcatgo@v0.4.2"
```

## Receive webhooks

```go
func bind(w http.ResponseWriter, r *http.Request) error {
    var webhookEvent revcatgo.WebhookEvent

    err := json.NewDecoder(r.Body).Decode(&webhookEvent)
    if err != nil {
        return err
    }
    fmt.Println(webhookEvent.Type) // e.g. "INITIAL_PURCHASE"
    return nil
}
```

## RevenueCat webhooks specification

See the official document.

* https://docs.revenuecat.com/docs/webhooks

[godev-image]: https://pkg.go.dev/badge/github.com/iktakahiro/revcatgo
[godev-url]: https://pkg.go.dev/github.com/iktakahiro/revcatgo