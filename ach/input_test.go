package ach

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Lendiom/go-payarc"
)

func TestCreateAchChargeInput_JSONIncludesCompanyName(t *testing.T) {
	tests := []struct {
		name     string
		input    CreateAchChargeInput
		contains []string
		omits    []string
	}{
		{
			name: "Business Checking includes company_name",
			input: CreateAchChargeInput{
				CustomerID:    "cust_1",
				BankAccountID: "bank_1",
				AccountType:   payarc.ACHAccountTypeBusinessChecking,
				Currency:      payarc.CurrencyUSD,
				Amount:        12345,
				Type:          payarc.ACHFlowTypeDebit,
				SecCode:       AchCreateChargeSecCodeCorporateCashDisbursement,
				CompanyName:   "ACME Corp",
			},
			contains: []string{`"company_name":"ACME Corp"`, `"account_type":"Business Checking"`},
		},
		{
			name: "Personal Checking omits company_name when empty",
			input: CreateAchChargeInput{
				CustomerID:    "cust_1",
				BankAccountID: "bank_1",
				AccountType:   payarc.ACHAccountTypePersonalChecking,
				Currency:      payarc.CurrencyUSD,
				Amount:        12345,
				Type:          payarc.ACHFlowTypeDebit,
				SecCode:       AchCreateChargeSecCodeWeb,
			},
			omits: []string{`"company_name"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			out := string(data)
			for _, want := range tt.contains {
				if !strings.Contains(out, want) {
					t.Errorf("Marshal() = %s; expected to contain %s", out, want)
				}
			}
			for _, omit := range tt.omits {
				if strings.Contains(out, omit) {
					t.Errorf("Marshal() = %s; expected NOT to contain %s", out, omit)
				}
			}
		})
	}
}
