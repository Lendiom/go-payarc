package payarc

type ChargeStatus string

var (
	ChargeStatusSubmittedForSettlement ChargeStatus = "submitted_for_settlement"
	ChargeStatusSettled                ChargeStatus = "settled"
	ChargeStatusVoid                   ChargeStatus = "void"
)

type ChargeResponse struct {
	Charge   Charge   `json:"data"`
	Metadata Metadata `json:"meta"`
}

type ChargesResponse struct {
	Charges []Charge `json:"data"`
}

type Charge struct {
	Object               string `json:"object"`
	ID                   string `json:"id"`
	Type                 string `json:"type"`
	ChargeDescription    string `json:"charge_description"`
	StatementDescription string `json:"statement_description"`
	ExternalOrderID      int    `json:"external_order_id"`

	Amount         int `json:"amount"`
	AmountApproved int `json:"amount_approved"`
	AmountCaptured int `json:"amount_captured"`
	AmountRefunded int `json:"amount_refunded"`
	AmountVoided   int `json:"amount_voided"`

	ApplicationFeeAmount int `json:"application_fee_amount"`
	TipAmount            int `json:"tip_amount"`
	PayArcFees           int `json:"payarc_fees"`
	NetAmount            int `json:"net_amount"`
	Surcharge            int `json:"surcharge"`

	Captured    Boolean      `json:"captured"`
	IsRefunded  Boolean      `json:"is_refunded"`
	Status      ChargeStatus `json:"status"`
	UnderReview Boolean      `json:"under_review"`

	CardLevel      ChargeCardLevel `json:"card_level"`
	AuthCode       string          `json:"auth_code"`
	FailureCode    string          `json:"failure_code"`
	FailureMessage string          `json:"failure_message"`

	DoNotSendEmailToCustomer Boolean `json:"do_not_send_email_to_customer"`
	DoNotSendSmsToCustomer   Boolean `json:"do_not_send_sms_to_customer"`

	KountDetails        string `json:"kount_details"`
	KountStatus         string `json:"kount_status"`
	TsysResponseCode    string `json:"tsys_response_code"`
	HostResponseCode    string `json:"host_response_code"`
	HostResponseMessage string `json:"host_response_message"`
	HostReferenceNumber string `json:"host_reference_number"`

	CreatedBy string `json:"created_by"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`

	Card CardResponse `json:"card"`
}
