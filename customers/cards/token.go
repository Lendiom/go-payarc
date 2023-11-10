package cards

type TokenData struct {
	Token Token `json:"data"`
}

type Token struct {
	Id string `json:"id" form:"token_id,omitempty"`
}

type TokenInput struct {
	CardSource     string  `form:"card_source,omitempty"`
	CardNumber     string  `form:"card_number,omitempty"`
	ExpMonth       string  `form:"exp_month,omitempty"`
	ExpYear        string  `form:"exp_year,omitempty"`
	Ccv            string  `form:"cvv,omitempty"`
	CardHolderName *string `form:"card_holder_name,omitempty"`
	Address1       *string `form:"address_1,omitempty"`
	Address2       *string `form:"address_2,omitempty"`
	City           *string `form:"city,omitempty"`
	State          *string `form:"state,omitempty"`
	Zip            *string `form:"zip,omitempty"`
	Country        *string `form:"country,omitempty"`
	AuthorizeCard  *bool   `form:"authorize_card,omitempty"`
}
