package customers

import (
	"go-payarc/cards"
	"go-payarc/charges"
)

type CustomerLookup struct {
	Customers []Customer `json:"data"`
}

type SingleCustomerLookup struct {
	Customer Customer `json:"data"`
}

type Customer struct {
	Object           string             `json:"object"`
	CustomerId       string             `json:"customer_id"`
	Name             *string            `json:"name"`
	Email            string             `json:"email"`
	Description      *string            `json:"description"`
	PaymentOverdue   int                `json:"payment_overdue"`
	SendEmailAddress *string            `json:"send_email_address"`
	CcEmailAddress   *string            `json:"cc_email_address"`
	SourceId         *string            `json:"source_id"`
	Address1         *string            `json:"address_1"`
	Address2         *string            `json:"address_2"`
	City             *string            `json:"city"`
	State            *string            `json:"state"`
	Zip              *string            `json:"zip"`
	Phone            *string            `json:"phone"`
	Country          *string            `json:"country"`
	Cards            cards.CardData     `json:"card"`
	Charges          charges.ChargeData `json:"charge"`
}

type CustomerInput struct {
	Name             string  `form:"name,omitempty"`
	Email            string  `form:"email,omitempty"`
	Description      *string `form:"description,omitempty"`
	SendEmailAddress *string `form:"send_email_address,omitempty"`
	CcEmailAddress   *string `form:"cc_email_address,omitempty"`
	Address1         *string `form:"address_1,omitempty"`
	Address2         *string `form:"address_2,omitempty"`
	City             *string `form:"city,omitempty"`
	State            *string `form:"state,omitempty"`
	Zip              *string `form:"zip,omitempty"`
	Phone            *string `form:"phone,omitempty"`
	Country          *string `form:"country,omitempty"`
}
