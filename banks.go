package payarc

type BanksResponse struct {
	Banks []BankAccount `json:"data"`
}

type BankAccountType string

var (
	ACHAccountTypePersonalChecking BankAccountType = "Personal Checking"
	ACHAccountTypePersonalSavings  BankAccountType = "Personal Savings"
	ACHAccountTypeBusinessChecking BankAccountType = "Business Checking"
	ACHAccountTypeBusinessSavings  BankAccountType = "Business Savings"
)

type BankAccount struct {
	Object      string          `json:"object"`
	ID          string          `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	CompanyName string          `json:"company_name,omitempty"`
	AccountType BankAccountType `json:"account_type"`
	SecCode     string          `json:"sec_code"`

	RoutingNumber string  `json:"routing_number"`
	AccountNumber string  `json:"account_number"`
	Default       Boolean `json:"is_default"`
}

type BankAccountCreated struct {
	Object      string          `json:"object"`
	ID          string          `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	AccountType BankAccountType `json:"account_type"`
	SecCode     string          `json:"sec_code"`

	RoutingNumber string `json:"routing_number"`
	AccountNumber string `json:"account_number"`
}

type BankAccountResponse struct {
	Data BankAccount `json:"data"`
	Meta Metadata    `json:"meta"`
}

type BankAccountCreatedResponse struct {
	BankAccount BankAccountCreated `json:"data"`
	Meta        Metadata           `json:"meta"`
}
