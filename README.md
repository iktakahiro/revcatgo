# revcatgo

[![GoDev][godev-image]][godev-url]

![Run test](https://github.com/iktakahiro/revcatgo/workflows/Run%20test/badge.svg?branch=main)

A helper library for integrating server-side apps with the [RevenueCat](https://www.revenuecat.com) webhook service.

```bash
go get "github.com/iktakahiro/revcatgo@v1.1.0"
```

## Receiving webhooks

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

## RevenueCat webhooks specifications

Refer to the official documentation:

* <https://docs.revenuecat.com/docs/webhooks>

[godev-image]: https://pkg.go.dev/badge/github.com/iktakahiro/revcatgo
[godev-url]: https://pkg.go.dev/github.com/iktakahiro/revcatgo
