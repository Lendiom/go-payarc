package payarc

type CustomerResponse struct {
	Data Customer `json:"data"`
}

type Customer struct {
	Object           string          `json:"object"`
	ID               string          `json:"customer_id"`
	Name             *string         `json:"name"`
	Email            string          `json:"email"`
	Description      *string         `json:"description"`
	PaymentOverdue   int             `json:"payment_overdue"`
	SendEmailAddress *string         `json:"send_email_address"`
	CcEmailAddress   *string         `json:"cc_email_address"`
	SourceID         *string         `json:"source_id"`
	Address1         *string         `json:"address_1"`
	Address2         *string         `json:"address_2"`
	City             *string         `json:"city"`
	State            *string         `json:"state"`
	Zip              *string         `json:"zip"`
	Phone            *string         `json:"phone"`
	Country          *string         `json:"country"`
	BankAccounts     BanksResponse   `json:"bank_account"`
	Cards            CardsResponse   `json:"card"`
	Charges          ChargesResponse `json:"charge"`
}
