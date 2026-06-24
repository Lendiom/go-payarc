package payarc

// WebhookCharge is the charge shape PayArc delivers inside a charge webhook's
// api_response ("original.data"). It intentionally differs from Charge:
//
//   - PayArc encodes some fields as JSON strings in webhooks (e.g. "captured":"1",
//     "amount_approved":"39646") that are numbers/booleans in the charge detail and
//     list API responses. Unmarshalling a webhook payload into Charge therefore
//     fails, so this is a separate, tolerant type.
//   - The webhook payload does NOT include transaction_metadata (it is only listed
//     under meta.include as an available relation). The loan/payment metadata
//     Lendiom attaches lives in the event's request_payload, or can be fetched
//     authoritatively via charges.GetByID (which requests include=transaction_metadata).
//
// Only the stable fields a consumer needs to identify and route a charge are
// modeled; authoritative amounts and metadata should come from charges.GetByID.
type WebhookCharge struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Amount      int          `json:"amount"`
	AuthCode    string       `json:"auth_code"`
	FailureCode string       `json:"failure_code"`
	Status      ChargeStatus `json:"status"`
	// Captured is a string ("1"/"0") in webhook payloads. Use IsCaptured for a bool.
	Captured  string            `json:"captured"`
	CreatedAt int64             `json:"created_at"`
	Card      WebhookChargeCard `json:"card"`
}

// IsCaptured reports whether the webhook charge was captured.
func (c WebhookCharge) IsCaptured() bool {
	return c.Captured == "1"
}

// IsCard reports whether the charge was paid with a card (vs ACH). PayArc only
// delivers charge webhooks for card charges; ACH is not sent as a charge event.
func (c WebhookCharge) IsCard() bool {
	return c.Card.Data.ID != ""
}

// WebhookChargeCard wraps the card relation on a webhook charge.
type WebhookChargeCard struct {
	Data WebhookChargeCardData `json:"data"`
}

// WebhookChargeCardData holds the card fields exposed on a webhook charge.
type WebhookChargeCardData struct {
	ID string `json:"id"`
}

// WebhookChargeResponse models the api_response payload of a PayArc charge webhook
// event (WebhookEvent.ApiResponse), e.g. "Charges Created".
type WebhookChargeResponse struct {
	Headers  map[string]string   `json:"headers"`
	Original WebhookChargeResult `json:"original"`
}

// WebhookChargeResult is the "original" envelope wrapping the charge data.
type WebhookChargeResult struct {
	Data WebhookCharge `json:"data"`
}
