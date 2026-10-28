// Package revcatgo models RevenueCat webhook payloads and Subscriber API v1 responses.
//
// # Webhooks
//
// Use [VerifyWebhookSignature] with the raw request body before decoding it.
// A positive tolerance checks the signature timestamp; zero disables the
// timestamp check.
//
// Decode a webhook JSON body into [WebhookEvent] using [encoding/json].
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
