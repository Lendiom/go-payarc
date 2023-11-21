package customers

import "github.com/Lendiom/go-payarc/charges"

type Customers struct {
	Customers []CustomerData `json:"data"`
}

type Customer struct {
	Customer CustomerData `json:"data"`
}

type CustomerData struct {
	Object           string             `json:"object"`
	Id               string             `json:"customer_id"`
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
	Cards            CardData           `json:"card"`
	Charges          charges.ChargeData `json:"charge"`
}

type CardData struct {
	Cards []Card `json:"data"`
}

type Card struct {
	Object    string `json:"object"`
	Id        string `json:"id"`
	Last4     string `json:"last4digit"`
	ExpMonth  string `json:"exp_month"`
	ExpYear   string `json:"exp_year"`
	IsDefault int    `json:"is_default"`
}

type CustomerInput struct {
	Name             *string `form:"name,omitempty"`
	Email            *string `form:"email,omitempty"`
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
