package banks

import "github.com/Lendiom/go-payarc"

type BankAccountCreateSecCode string

var (
	BankAccountCreateSecCodeCorporateCashDisbursement BankAccountCreateSecCode = "CCD"
	BankAccountCreateSecCodePrearranged               BankAccountCreateSecCode = "PPD"
	BankAccountCreateSecCodeTelephone                 BankAccountCreateSecCode = "TEL"
	BankAccountCreateSecCodeWeb                       BankAccountCreateSecCode = "WEB"
)

type CreateBankAccountInput struct {
	AccountNumber string                   `form:"account_number"` //AccountNumber must be 3 to 17 characters
	RoutingNumber string                   `form:"routing_number"` //RoutingNumber must be 9 characters
	FirstName     string                   `form:"first_name"`
	LastName      string                   `form:"last_name"`
	AccountType   payarc.BankAccountType   `form:"account_type"`
	SecCode       BankAccountCreateSecCode `form:"sec_code"` //SecCode must be one of the following: ARC, BOC, CCD, POP, PPD, RCK, TEL, WEB
	CustomerID    string                   `form:"customer_id"`
}
