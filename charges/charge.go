package charges

import "github.com/Lendiom/go-payarc"

type ChargeInput struct {
	Amount            int64          `form:"amount"`
	Capture           payarc.Boolean `form:"capture"`
	CustomerID        string         `form:"customer_id"`
	CardID            string         `form:"card_id,omitempty"`
	ExternalOrderID   *int           `form:"external_order_id,omitempty"`
	ChargeDescription string         `form:"charge_description,omitempty"`

	Currency                 payarc.Currency `form:"currency"`
	StatementDescription     *string         `form:"statement_description,omitempty"`
	Metadata                 *string         `form:"metadata,omitempty"` //TODO: convert this to map[string]string
	DoNotSendEmailToCustomer payarc.YesOrNo  `form:"do_not_send_email_to_customer"`
	DoNotSendSmsToCustomer   payarc.YesOrNo  `form:"do_not_send_sms_to_customer"`
}

type CreateChargeResponse struct {
	Charge ChargeResult `json:"data"`
}

type ChargeResult struct {
	Object               string `json:"object"`
	ID                   string `json:"id"`
	Type                 string `json:"type"`
	ChargeDescription    string `json:"charge_description"`
	StatementDescription string `json:"statement_description"`
	ExternalOrderID      int    `json:"external_order_id"`

	Amount         int    `json:"amount"`
	AmountApproved string `json:"amount_approved"`
	AmountCaptured int    `json:"amount_captured"`
	AmountRefunded int    `json:"amount_refunded"`
	AmountVoided   int    `json:"amount_voided"`

	ApplicationFeeAmount int `json:"application_fee_amount"`
	TipAmount            int `json:"tip_amount"`
	PayArcFees           int `json:"payarc_fees"`
	NetAmount            int `json:"net_amount"`
	Surcharge            int `json:"surcharge"`

	Captured    string              `json:"captured"` //Captured is a string 1 or 0, which is true or false respectively
	IsRefunded  payarc.Boolean      `json:"is_refunded"`
	Status      payarc.ChargeStatus `json:"status"`
	UnderReview bool                `json:"under_review"`

	CardLevel      payarc.ChargeCardLevel `json:"card_level"`
	AuthCode       string                 `json:"auth_code"`
	FailureCode    string                 `json:"failure_code"`
	FailureMessage string                 `json:"failure_message"`

	DoNotSendEmailToCustomer payarc.Boolean `json:"do_not_send_email_to_customer"`
	DoNotSendSmsToCustomer   payarc.Boolean `json:"do_not_send_sms_to_customer"`

	KountDetails        string `json:"kount_details"`
	KountStatus         string `json:"kount_status"`
	TsysResponseCode    string `json:"tsys_response_code"`
	HostResponseCode    string `json:"host_response_code"`
	HostResponseMessage string `json:"host_response_message"`
	HostReferenceNumber string `json:"host_reference_number"`

	CreatedBy string `json:"created_by"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`

	Card payarc.CardResponse `json:"card"`
}
