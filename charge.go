package payarc

type ChargeStatus string

var (
	ChargeStatusSubmittedForSettlement ChargeStatus = "submitted_for_settlement"
	ChargeStatusSettled                ChargeStatus = "settled"
	ChargeStatusVoid                   ChargeStatus = "void"
)

type ChargeResponse struct {
	Charge   Charge   `json:"data,omitempty"`
	Metadata Metadata `json:"meta,omitempty"`
}

type ChargesResponse struct {
	Charges  []Charge `json:"data,omitempty"`
	Metadata Metadata `json:"meta,omitempty"`
}

type Charge struct {
	Object               string `json:"object"`
	ID                   string `json:"id"`
	Type                 string `json:"type"`
	ChargeDescription    string `json:"charge_description,omitempty"`
	StatementDescription string `json:"statement_description,omitempty"`
	ExternalOrderID      int    `json:"external_order_id,omitempty"`

	Amount         int `json:"amount,omitempty"`
	AmountApproved int `json:"amount_approved,omitempty"`
	AmountCaptured int `json:"amount_captured,omitempty"`
	AmountRefunded int `json:"amount_refunded,omitempty"`
	AmountVoided   int `json:"amount_voided,omitempty"`

	ApplicationFeeAmount int `json:"application_fee_amount,omitempty"`
	TipAmount            int `json:"tip_amount,omitempty"`
	PayArcFees           int `json:"payarc_fees,omitempty"`
	NetAmount            int `json:"net_amount,omitempty"`
	Surcharge            int `json:"surcharge,omitempty"`

	Captured    Boolean      `json:"captured,omitempty"`
	IsRefunded  Boolean      `json:"is_refunded,omitempty"`
	Status      ChargeStatus `json:"status,omitempty"`
	UnderReview Boolean      `json:"under_review,omitempty"`

	CardLevel      ChargeCardLevel `json:"card_level"`
	AuthCode       string          `json:"auth_code"`
	FailureCode    string          `json:"failure_code,omitempty"`
	FailureMessage string          `json:"failure_message,omitempty"`

	DoNotSendEmailToCustomer Boolean `json:"do_not_send_email_to_customer,omitempty"`
	DoNotSendSmsToCustomer   Boolean `json:"do_not_send_sms_to_customer,omitempty"`

	KountDetails        string `json:"kount_details,omitempty"`
	KountStatus         string `json:"kount_status,omitempty"`
	TsysResponseCode    string `json:"tsys_response_code,omitempty"`
	HostResponseCode    string `json:"host_response_code,omitempty"`
	HostResponseMessage string `json:"host_response_message,omitempty"`
	HostReferenceNumber string `json:"host_reference_number,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`

	Card CardResponse `json:"card,omitempty"`

	TransactionMetadata *ChargeTransactionMetadataResponse `json:"transaction_metadata,omitempty"`
}

type ChargeTransactionMetadataResponse struct {
	Data []ChargeTransactionMetadata `json:"data,omitempty"`
}

type ChargeTransactionMetadata struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Signature string `json:"signature,omitempty"`
}
