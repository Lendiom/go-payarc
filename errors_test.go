package payarc

import (
	"errors"
	"testing"
)

func TestClassifyChargeDecline(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		wantErr     error
		wantMatched bool
	}{
		{
			name:        "Insufficient Funds (canonical casing)",
			message:     "Insufficient Funds",
			wantErr:     ErrInsufficientFunds,
			wantMatched: true,
		},
		{
			name:        "insufficient funds (lowercase) still matches",
			message:     "insufficient funds",
			wantErr:     ErrInsufficientFunds,
			wantMatched: true,
		},
		{
			name:        "INSUFFICIENT FUNDS (uppercase) still matches",
			message:     "INSUFFICIENT FUNDS",
			wantErr:     ErrInsufficientFunds,
			wantMatched: true,
		},
		{
			name:        "Expired Card",
			message:     "Expired Card",
			wantErr:     ErrExpiredCard,
			wantMatched: true,
		},
		{
			name:        "Do Not Honor",
			message:     "Do Not Honor",
			wantErr:     ErrDoNotHonor,
			wantMatched: true,
		},
		{
			name:        "Closed Account",
			message:     "Closed Account",
			wantErr:     ErrClosedAccount,
			wantMatched: true,
		},
		{
			name:        "Customer Requested Stop (long message verbatim)",
			message:     "Customer Requested Stop of All Recurring Payments from Specific Merchant",
			wantErr:     ErrCustomerRequestedStopPayments,
			wantMatched: true,
		},
		{
			name:        "General CardAuth Decline",
			message:     "General CardAuth Decline",
			wantErr:     ErrGeneralCardAuthDecline,
			wantMatched: true,
		},
		{
			name:        "Refer to Issuer (canonical casing)",
			message:     "Refer to Issuer",
			wantErr:     ErrReferToIssuer,
			wantMatched: true,
		},
		{
			name:        "refer to issuer (lowercase) still matches",
			message:     "refer to issuer",
			wantErr:     ErrReferToIssuer,
			wantMatched: true,
		},
		{
			name:        "Unknown message does not match (must stay at ERROR)",
			message:     "Some new failure mode PayArc invented",
			wantErr:     nil,
			wantMatched: false,
		},
		{
			name:        "Empty message does not match",
			message:     "",
			wantErr:     nil,
			wantMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotMatched := ClassifyChargeDecline(tt.message)

			if gotMatched != tt.wantMatched {
				t.Errorf("matched = %v, want %v", gotMatched, tt.wantMatched)
			}

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("err = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
