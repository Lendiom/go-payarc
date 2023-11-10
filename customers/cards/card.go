package cards

type CardData struct {
	Cards Card `json:"data"`
}

type Card struct {
	CardSource string `json:"card_source"`
	CardNumber int    `json:"card_number"`
	ExpMonth   int    `json:"exp_month"`
	ExpYear    int    `json:"exp_year"`
	Ccv        string `json:"cvv"`
}
