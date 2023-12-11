package payarc

type CardSource string

var (
	CardSourceManual   CardSource = "MANUAL"
	CardSourcePhone    CardSource = "PHONE"
	CardSourceMail     CardSource = "MAIL"
	CardSourceInternet CardSource = "INTERNET"
)

func (cs CardSource) Valid() bool {
	switch cs {
	case CardSourceManual, CardSourcePhone, CardSourceMail, CardSourceInternet:
		return true
	default:
		return false
	}
}

type CardBrand string

var (
	CardBrandVisa            CardBrand = "V"
	CardBrandMastercard      CardBrand = "M"
	CardBrandDiscover        CardBrand = "R"
	CardBrandAmericanExpress CardBrand = "X"
)

type CardsResponse struct {
	Cards []Card `json:"data"`
}

type CardResponse struct {
	Data Card `json:"data,omitempty"`
}

type Card struct {
	Object      string     `json:"object,omitempty"`
	ID          string     `json:"id,omitempty"`
	CustomerID  string     `json:"customer_id,omitempty"`
	Brand       CardBrand  `json:"brand,omitempty"`
	First6Digit int        `json:"first6digit,omitempty"`
	Last4Digit  string     `json:"last4digit,omitempty"`
	ExpMonth    string     `json:"exp_month,omitempty"`
	ExpYear     string     `json:"exp_year,omitempty"`
	Fingerprint string     `json:"fingerprint,omitempty"`
	CardSource  CardSource `json:"card_source,omitempty"`

	IsVerified Boolean `json:"is_verified,omitempty"`
	IsDefault  Boolean `json:"is_default,omitempty"`

	HolderName string `json:"card_holder_name,omitempty"`
	Address1   string `json:"address1,omitempty"`
	Address2   string `json:"address2,omitempty"`
	State      string `json:"state,omitempty"`
	City       string `json:"city,omitempty"`
	Zip        string `json:"zip,omitempty"`
	Country    string `json:"country,omitempty"`

	AvsStatus          string  `json:"avs_status,omitempty"`
	CvcStatus          string  `json:"cvc_status,omitempty"`
	AddressCheckPassed Boolean `json:"address_check_passed,omitempty"`
	ZipCheckPassed     Boolean `json:"zip_check_passed,omitempty"`

	CardType    string `json:"card_type,omitempty"`
	BinCountry  string `json:"bin_country,omitempty"`
	BankName    any    `json:"bank_name,omitempty"`
	BankWebsite any    `json:"bank_website,omitempty"`
	BankPhone   any    `json:"bank_phone,omitempty"`

	CreatedAt int `json:"created_at,omitempty"`
	UpdatedAt int `json:"updated_at,omitempty"`
}
