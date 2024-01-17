package payarc

import "time"

type ACHFlowType string

var (
	ACHFlowTypeDebit  ACHFlowType = "debit"
	ACHFlowTypeCredit ACHFlowType = "credit"
)

type ACHChargeStatus string

var (
	ACHChargeStatusValidated ACHChargeStatus = "status"
)

type ACHCharge struct {
	Object              string              `json:"object"`
	ID                  string              `json:"id"`
	Amount              string              `json:"amount"`
	Status              ACHChargeStatus     `json:"status"`
	Type                ACHFlowType         `json:"type"`
	AuthorizationID     int                 `json:"authorization_id"`
	ValidationCode      string              `json:"validation_code"`
	Successful          bool                `json:"successful"`
	ResponseMessage     string              `json:"response_message"`
	CreatedBy           string              `json:"created_by"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	RetriedACHChargeID  any                 `json:"retried_achcharge_id"`
	BankAccountResponse BankAccountResponse `json:"bank_account"`
}

type ACHChargesResponse struct {
	Charges  []ACHCharge `json:"data"`
	Metadata Metadata    `json:"meta"`
}

type ACHChargeResponse struct {
	Charge   ACHCharge `json:"data"`
	Metadata Metadata  `json:"meta"`
}
