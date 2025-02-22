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
	CustomerID    string `json:"customer_id,omitempty"`     //CustomerID is required if the AccountNumber, RoutingNumber, FirstName, LastName are not provided
	BankAccountID string `json:"bank_account_id,omitempty"` //BankAccountID is required if the AccountNumber and Routing number are not provided

	AccountType payarc.BankAccountType `json:"account_type"`
	Currency    payarc.Currency        `json:"currency"` //Currency is the three letter ISO currency code. Currently on usd is allowed.
	Amount      int64                  `json:"amount"`   //Amount is a positive integer in cents representing how much to charge.
	Type        payarc.ACHFlowType     `json:"type"`     //Type is required
	SecCode     AchCreateChargeSecCode `json:"sec_code"` //SecCode must be one of the following: ARC, BOC, CCD, POP, PPD, RCK, TEL, WEB

	AccountNumber string `json:"account_number,omitempty"` //AccountNumber is required if the CustomerID is not provided
	RoutingNumber string `json:"routing_number,omitempty"` //RoutingNumber is required if the CustomerID is not provided
	FirstName     string `json:"first_name,omitempty"`     //FirstName is required if the CustomerID is not provided
	LastName      string `json:"last_name,omitempty"`      //LastName is required if the CustomerID is not provided
	ReceiptEmail  string `json:"receipt_email,omitempty"`  //ReceiptEmail is optional
	ReceiptPhone  string `json:"receipt_phone,omitempty"`  //ReceiptPhone is optional
	AddressLine1  string `json:"address_line1,omitempty"`  //AddressLine1 is optional
	Zip           string `json:"zip,omitempty"`            //Zip is optional
}

type CreateACHChargeResponse struct {
	Charge ACHChargeResult `json:"data"`
	Meta   payarc.Metadata `json:"meta"`
}

type ACHChargeResult struct {
	Object             string                     `json:"object"`
	ID                 string                     `json:"id"`
	Amount             int64                      `json:"amount"`
	CreatedBy          string                     `json:"created_by"`
	Status             string                     `json:"status"`
	Type               payarc.ACHFlowType         `json:"type"`
	AuthorizationID    int64                      `json:"authorization_id"`
	ValidationCode     int64                      `json:"validation_code"`
	Successful         bool                       `json:"successful"`
	ResponseMessage    string                     `json:"response_message"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	RetriedAchChargeID any                        `json:"retried_achcharge_id"`
	SecCode            AchCreateChargeSecCode     `json:"sec_code"`
	BankAccount        payarc.BankAccountResponse `json:"bank_account"`
}
