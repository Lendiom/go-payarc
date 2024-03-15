package ach

import (
	"time"

	"github.com/Lendiom/go-payarc"
)

type AchCreateChargeSecCode string

var (
	AchCreateChargeSecCodeCorporateCashDisbursement AchCreateChargeSecCode = "CCD"
	AchCreateChargeSecCodePrearranged               AchCreateChargeSecCode = "PPD"
	AchCreateChargeSecCodeTelephone                 AchCreateChargeSecCode = "TEL"
	AchCreateChargeSecCodeWeb                       AchCreateChargeSecCode = "WEB"
)

type CreateAchChargeInput struct {
	CustomerID string `form:"customer_id"` //CustomerID is required if the AccountNumber, RoutingNumber, FirstName, LastName are not provided

	AccountType payarc.BankAccountType `form:"account_type"`
	Currency    payarc.Currency        `form:"currency"` //Currency is the three letter ISO currency code. Currently on usd is allowed.
	Amount      int64                  `form:"amount"`   //Amount is a positive integer in cents representing how much to charge.
	Type        payarc.ACHFlowType     `form:"type"`     //Type is required
	SecCode     AchCreateChargeSecCode `form:"sec_code"` //SecCode must be one of the following: ARC, BOC, CCD, POP, PPD, RCK, TEL, WEB

	AccountNumber string `form:"account_number,omitempty"` //AccountNumber is required if the CustomerID is not provided
	RoutingNumber string `form:"routing_number,omitempty"` //RoutingNumber is required if the CustomerID is not provided
	FirstName     string `form:"first_name,omitempty"`     //FirstName is required if the CustomerID is not provided
	LastName      string `form:"last_name,omitempty"`      //LastName is required if the CustomerID is not provided
	ReceiptEmail  string `form:"receipt_email,omitempty"`  //ReceiptEmail is optional
	ReceiptPhone  string `form:"receipt_phone,omitempty"`  //ReceiptPhone is optional
	AddressLine1  string `form:"address_line1,omitempty"`  //AddressLine1 is optional
	Zip           string `form:"zip,omitempty"`            //Zip is optional
}

type CreateACHChargeResponse struct {
	Charge ACHChargeResult `json:"data"`
	Meta   payarc.Metadata `json:"meta"`
}

type ACHChargeResult struct {
	Object             string                     `json:"object"`
	ID                 string                     `json:"id"`
	Amount             string                     `json:"amount"`
	CreatedBy          string                     `json:"created_by"`
	Status             string                     `json:"status"`
	Type               payarc.ACHFlowType         `json:"type"`
	AuthorizationID    int                        `json:"authorization_id"`
	ValidationCode     any                        `json:"validation_code"`
	Successful         bool                       `json:"successful"`
	ResponseMessage    string                     `json:"response_message"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	RetriedAchChargeID any                        `json:"retried_achcharge_id"`
	BankAccount        payarc.BankAccountResponse `json:"bank_account"`
}
