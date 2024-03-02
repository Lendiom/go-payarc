package customers

import (
	"github.com/Lendiom/go-payarc"
)

type CustomersResponse struct {
	Customers []payarc.Customer `json:"data"`
}

type CustomerResponse struct {
	Customer payarc.Customer `json:"data"`
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
