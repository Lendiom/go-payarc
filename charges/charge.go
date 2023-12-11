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
