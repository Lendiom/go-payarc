package customers

import "github.com/Lendiom/go-payarc"

type TokenResponse struct {
	Data Token `json:"data"`
}

type Token struct {
	Object             string `json:"object,omitempty" form:"-"`
	ID                 string `json:"id" form:"token_id,omitempty"`
	Used               bool   `json:"used,omitempty" form:"-"`
	TokenizationMethod any    `json:"tokenization_method,omitempty" form:"-"`

	Card payarc.CardResponse `json:"card" form:"-"`

	CreatedAt int `json:"created_at,omitempty" form:"-"`
	UpdatedAt int `json:"updated_at,omitempty" form:"-"`
}

//#region token creation input
type TokenInput struct {
	CardSource     payarc.CardSource `form:"card_source,omitempty"`
	CardNumber     string            `form:"card_number,omitempty"`
	ExpMonth       string            `form:"exp_month,omitempty"`
	ExpYear        string            `form:"exp_year,omitempty"`
	CCV            string            `form:"cvv,omitempty"`
	CardHolderName *string           `form:"card_holder_name,omitempty"`
	Address1       *string           `form:"address_line1,omitempty"`
	Address2       *string           `form:"address_line2,omitempty"`
	City           *string           `form:"city,omitempty"`
	State          *string           `form:"state,omitempty"`
	Zip            *string           `form:"zip,omitempty"`
	Country        *string           `form:"country,omitempty"`
	AuthorizeCard  *payarc.Boolean   `form:"authorize_card,omitempty"`
}

//#endregion token creation input
