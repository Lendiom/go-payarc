package payarc

import "encoding/json"

type WebhookEventType string

var (
	WebhookEventTypeChargeCreated WebhookEventType = "Charges Created"
)

type Event struct {
	RequestPayload json.RawMessage `json:"request_payload"`
	ApiResponse    json.RawMessage `json:"api_response"`
}
