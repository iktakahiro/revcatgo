// Package revcatgo models RevenueCat webhook payloads and Subscriber API v1 responses.
//
// # Webhooks
//
// Use [ParseWebhookRequest] to read, verify, and decode an HTTP webhook request.
// Alternatively, use [VerifyWebhookSignature] with the raw request body before
// decoding it using [encoding/json]. A positive tolerance checks the signature
// timestamp; zero disables the timestamp check.
//
// Webhook JSON bodies decode into [WebhookEvent].
// Event.Type and Event.Environment provide typed values for the event type and
// environment; use String to compare them with the corresponding constants and
// IsKnown to check whether RevenueCat documents the value. Unknown values are
// preserved during decoding so new event types and enum values do not break
// existing integrations.
//
// # Subscriber API
//
// Decode Subscriber API v1 JSON responses into [SubscriberResponse] using
// [encoding/json].
package revcatgo
