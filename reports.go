package payarc

type DepositReport struct {
	DateReports map[string]DepositReportDate `json:"data"`
	GrandTotal  DepositReportGrandTotal      `json:"grand_total"`
}

type DepositReportGrandTotal struct {
	TotalAmount      int64 `json:"total_amount"`
	TotalTransaction int64 `json:"total_transaction"`
}

type DepositReportDate struct {
	Data   []DepositReportDateData `json:"row_data"`
	Totals DepositReportDateTotals `json:"row_totals"`
}

type DepositReportDateData struct {
	MerchantAccountNumber string      `json:"Merchant_Account_Number"`
	SettlementDate        DateTime    `json:"Settlement_Date"`
	AdTotalSale           interface{} `json:"ad_total_sale"`
	AdTotalRefunds        interface{} `json:"ad_total_refunds"`
	AdNetAmt              interface{} `json:"ad_net_amt"`
	Amounts               int64       `json:"Amounts"`       //Amounts is in cents
	TotalRefunds          int64       `json:"total_refunds"` //TotalRefunds is in cents
	TotalNetAmt           int64       `json:"total_net_amt"` //TotalNetAmt is in cents
	RjTotalSale           interface{} `json:"rj_total_sale"`
	RjTotalRefunds        interface{} `json:"rj_total_refunds"`
	RjNetAmt              interface{} `json:"rj_net_amt"`
	Transactions          int64       `json:"Transactions"`
	BatchReferenceNumber  int64       `json:"Batch_Reference_Number"`
	RejectRecord          interface{} `json:"reject_record"`
}

type DepositReportDateTotals struct {
	Amounts      int64 `json:"Amounts"`
	Transactions int64 `json:"Transactions"`
}

type DepositBatchDetailsResponse struct {
	Data []DepositBatchDetails `json:"data"`
	Meta Metadata              `json:"meta"`
}

type DepositBatchDetails struct {
	Object                  string   `json:"object"`
	ID                      string   `json:"id"`
	MerchantAccountNumber   string   `json:"merchant_account_number"`
	DBA                     string   `json:"dba"`
	TransactionAmount       string   `json:"transaction_amount"`
	TransactionType         string   `json:"transaction_type"`
	TransactionDate         DateTime `json:"transaction_date"`
	CardType                string   `json:"card_type"`
	CardholderAccountNumber string   `json:"cardholder_account_number"`
	POSEntryMode            string   `json:"pos_entry_mode"`
	AuthorizationNumber     string   `json:"authorization_number"`
	BatchDate               DateTime `json:"batch_date"`
	DebitCreditIndicator    string   `json:"debit_credit_indicator"`
}
