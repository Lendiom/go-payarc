package deposits

import "github.com/Lendiom/go-payarc"

type Report struct {
	DateReports map[string]ReportDate `json:"data"`
	GrandTotal  ReportGrandTotal      `json:"grand_total"`
}

type ReportGrandTotal struct {
	TotalAmount      int `json:"total_amount"`
	TotalTransaction int `json:"total_transaction"`
}

type ReportDate struct {
	Data   []ReportDateData `json:"row_data"`
	Totals ReportDateTotals `json:"row_totals"`
}

type ReportDateData struct {
	MerchantAccountNumber string          `json:"Merchant_Account_Number"`
	SettlementDate        payarc.DateTime `json:"Settlement_Date"`
	AdTotalSale           interface{}     `json:"ad_total_sale"`
	AdTotalRefunds        interface{}     `json:"ad_total_refunds"`
	AdNetAmt              interface{}     `json:"ad_net_amt"`
	Amounts               int             `json:"Amounts"`       //Amounts is in cents
	TotalRefunds          int             `json:"total_refunds"` //TotalRefunds is in cents
	TotalNetAmt           int             `json:"total_net_amt"` //TotalNetAmt is in cents
	RjTotalSale           interface{}     `json:"rj_total_sale"`
	RjTotalRefunds        interface{}     `json:"rj_total_refunds"`
	RjNetAmt              interface{}     `json:"rj_net_amt"`
	Transactions          int             `json:"Transactions"`
	BatchReferenceNumber  int64           `json:"Batch_Reference_Number"`
	RejectRecord          interface{}     `json:"reject_record"`
}

type ReportDateTotals struct {
	Amounts      int `json:"Amounts"`
	Transactions int `json:"Transactions"`
}
