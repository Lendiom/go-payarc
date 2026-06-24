package payarc

import (
	"encoding/json"
	"testing"
)

// realWebhookApiResponse is a trimmed but faithful copy of a real PayArc
// "Charges Created" api_response. Note "captured" and "amount_approved" are JSON
// strings (not numbers), and there is no transaction_metadata under data.
const realWebhookApiResponse = `{
  "headers": {},
  "original": {
    "data": {
      "object": "Charge",
      "id": "nDXOXWWnRnbBbObL",
      "amount": 39646,
      "amount_approved": "39646",
      "amount_refunded": 0,
      "type": "Sale",
      "captured": "1",
      "is_refunded": 0,
      "status": "submitted_for_settlement",
      "auth_code": "202054",
      "failure_code": null,
      "created_at": 1782270224,
      "card": {"data": {"object": "Card", "id": "19v20vP2N9M1L0M5", "last4digit": "5834"}},
      "splits": {"data": []}
    },
    "meta": {"include": ["review", "transaction_metadata", "splits"]}
  },
  "exception": null
}`

func TestWebhookChargeResponse_ParsesRealPayload(t *testing.T) {
	var resp WebhookChargeResponse
	if err := json.Unmarshal([]byte(realWebhookApiResponse), &resp); err != nil {
		t.Fatalf("failed to unmarshal real webhook api_response: %v", err)
	}

	charge := resp.Original.Data

	if charge.ID != "nDXOXWWnRnbBbObL" {
		t.Errorf("id = %q, want %q", charge.ID, "nDXOXWWnRnbBbObL")
	}

	if !charge.IsCaptured() {
		t.Errorf("IsCaptured() = false, want true (captured=%q)", charge.Captured)
	}

	if !charge.IsCard() {
		t.Error("IsCard() = false, want true")
	}

	if charge.Card.Data.ID != "19v20vP2N9M1L0M5" {
		t.Errorf("card.data.id = %q, want %q", charge.Card.Data.ID, "19v20vP2N9M1L0M5")
	}

	if charge.AuthCode != "202054" {
		t.Errorf("auth_code = %q, want %q", charge.AuthCode, "202054")
	}

	if charge.FailureCode != "" {
		t.Errorf("failure_code = %q, want empty", charge.FailureCode)
	}

	if charge.CreatedAt != 1782270224 {
		t.Errorf("created_at = %d, want %d", charge.CreatedAt, 1782270224)
	}

	if charge.Status != ChargeStatus("submitted_for_settlement") {
		t.Errorf("status = %q, want submitted_for_settlement", charge.Status)
	}
}

// TestCharge_FailsOnWebhookStrings documents WHY WebhookCharge exists: the
// API-shaped Charge cannot unmarshal a webhook payload because captured/amount_approved
// arrive as strings.
func TestCharge_FailsOnWebhookStrings(t *testing.T) {
	var resp WebhookEventChargeResponse
	if err := json.Unmarshal([]byte(realWebhookApiResponse), &resp); err == nil {
		t.Fatal("expected Charge unmarshal to fail on string-encoded captured/amount_approved, but it succeeded")
	}
}
