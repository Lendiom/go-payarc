package charges

import "go-payarc/cards"

type ChargeData struct {
	Charges []Charge `json:"data"`
}

type Charge struct {
	Object            string         `json:"object"`
	ChargeId          string         `json:"charge_id"`
	RealId            int            `json:"real_id"`
	Amount            int            `json:"amount"`
	AmountRefunded    int            `json:"amount_refunded"`
	AmountCaptured    int            `json:"amount_captured"`
	AmountVoided      int            `json:"amount_voided"`
	PayArcFees        int            `json:"payarc_fees"`
	Type              string         `json:"type"`
	NetAmount         int            `json:"net_amount"`
	Captured          int            `json:"captured"`
	IsRefunded        int            `json:"is_refunded"`
	Status            string         `json:"status"`
	AuthCode          string         `json:"auth_code"`
	FailureCode       *string        `json:"failure_code"`
	FailureMessage    *string        `json:"failure_message"`
	ChargeDescription *string        `json:"charge_description"`
	UnderRevies       int            `json:"under_review"`
	Cards             cards.CardData `json:"card"`
	// INVOICE?
	// Refunds       refunds.Refund `json:"refund"`
}

type ChargeInput struct {
	Amount               int     `form:"amount,omitempty"`
	CustomerId           string  `form:"customer_id,omitempty"`
	Currency             string  `form:"currency,omitempty"`
	StatementDescription *string `form:"statement_description,omitempty"`
	Email                string  `form:"email,omitempty"`
	PhoneNumber          *string `form:"phone_number,omitempty"`
	CardId               string  `form:"card_id,omitempty"`
}
