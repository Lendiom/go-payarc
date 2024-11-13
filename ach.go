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
	Object             string      `json:"object"`
	ID                 string      `json:"id"`
	Type               ACHFlowType `json:"type"`
	SecCode            string      `json:"sec_code"`
	RetriedACHChargeID any         `json:"retried_achcharge_id,omitempty"`

	Status          ACHChargeStatus `json:"status"`
	AuthorizationID string          `json:"authorization_id"`
	ValidationCode  string          `json:"validation_code"`

	Amount          int     `json:"amount"`
	Successful      Boolean `json:"successful"`
	ResponseMessage string  `json:"response_message"`

	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	BankAccount BankAccountResponse `json:"bank_account"`

	Customer *CustomerResponse `json:"customer,omitempty"`
}

type ACHChargesResponse struct {
	Charges  []ACHCharge `json:"data"`
	Metadata Metadata    `json:"meta"`
}

type ACHChargeResponse struct {
	Charge   ACHCharge `json:"data"`
	Metadata Metadata  `json:"meta"`
}
