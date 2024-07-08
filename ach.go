package payarc

import "time"

type ACHFlowType string

var (
	ACHFlowTypeDebit  ACHFlowType = "debit"
	ACHFlowTypeCredit ACHFlowType = "credit"
)

type ACHChargeStatus string

var (
	ACHChargeStatusSettled  ACHChargeStatus = "settled"
	ACHChargeStatusPending  ACHChargeStatus = "validated"
	ACHChargeStatusRejected ACHChargeStatus = "rejected"
)

type ACHCharge struct {
	Object              string              `json:"object"`
	ID                  string              `json:"id"`
	Amount              int                 `json:"amount"`
	Status              ACHChargeStatus     `json:"status"`
	Type                ACHFlowType         `json:"type"`
	AuthorizationID     string              `json:"authorization_id"`
	ValidationCode      string              `json:"validation_code"`
	Successful          Boolean             `json:"successful"`
	ResponseMessage     string              `json:"response_message"`
	SecCode             string              `json:"sec_code"`
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
