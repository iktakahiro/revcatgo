// Package revcatgo models RevenueCat webhook payloads and Subscriber API v1 responses.
//
// # Webhooks
//
// Use [ParseWebhookRequest] to read an HTTP webhook request (up to 1 MiB,
// [MaxWebhookBodyBytes]), verify its signature with a non-empty secret, and
// decode it into [WebhookEvent]. An empty secret returns an error. Alternatively,
// use [VerifyWebhookSignature] on the raw body before decoding with [encoding/json].
// For both functions, a positive tolerance checks the signature timestamp;
// zero disables the timestamp check.
//
// In a decoded [WebhookEvent], Event.Type and Event.Environment provide typed
// values for the event type and environment; use String to compare them with the
// corresponding constants and IsKnown to check whether RevenueCat documents the
// value. Unknown values are
// preserved during decoding so new event types and enum values do not break
// existing integrations.
//
// # Subscriber API
//
// Decode Subscriber API v1 JSON responses into [SubscriberResponse] using
// [encoding/json].
package revcatgo
