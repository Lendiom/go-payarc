package payarc

type WebhookEventType string

var (
	WebhookEventTypeChargeCreated   WebhookEventType = "Charges Created"
	WebhookEventTypeChargeCaptured  WebhookEventType = "Charge Captured"
	WebhookEventTypeChargeVoided    WebhookEventType = "Charge Voided"
	WebhookEventTypeChargeFailed    WebhookEventType = "Charge Failure"
	WebhookEventTypeCustomerCreated WebhookEventType = "Customers Created"
	WebhookEventTypeCustomerUpdated WebhookEventType = "Customers Updated"
	WebhookEventTypeCustomerDeleted WebhookEventType = "Customers Deleted"
	WebhookEventTypeCardCreated     WebhookEventType = "Card Created"
	WebhookEventTypeCardUpdated     WebhookEventType = "Card Updated"
	WebhookEventTypeCardDeleted     WebhookEventType = "Card Deleted"
)

type WebhookEvent struct {
	RequestPayload string           `json:"request_payload"`
	ApiResponse    string           `json:"api_response"`
	Type           WebhookEventType `json:"event_type"`
}

type WebhookEventChargeResponse struct {
	Headers        map[string]string `json:"headers"`
	ChargeResponse ChargeResponse    `json:"original"`
}
