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

func TestClassifyCardTokenDecline(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		wantErr     error
		wantMatched bool
	}{
		{
			name:        "Invalid Card",
			message:     "Invalid Card",
			wantErr:     ErrInvalidCard,
			wantMatched: true,
		},
		{
			name:        "Invalid CVV",
			message:     "Invalid CVV",
			wantErr:     ErrInvalidCCV,
			wantMatched: true,
		},
		{
			name:        "CVV2 Verification Failed (canonical casing)",
			message:     "CVV2 Verification Failed",
			wantErr:     ErrCVV2Failed,
			wantMatched: true,
		},
		{
			name:        "cvv2 verification failed (lowercase) still matches",
			message:     "cvv2 verification failed",
			wantErr:     ErrCVV2Failed,
			wantMatched: true,
		},
		{
			name:        "Suspected Fraud",
			message:     "Suspected Fraud",
			wantErr:     ErrSuspectedFraud,
			wantMatched: true,
		},
		{
			name:        "Do Not Honor",
			message:     "Do Not Honor",
			wantErr:     ErrDoNotHonor,
			wantMatched: true,
		},
		{
			name:        "Suspected Card",
			message:     "Suspected Card",
			wantErr:     ErrSuspectedCard,
			wantMatched: true,
		},
		{
			name:        "Charge-only message does not match (Insufficient Funds is not a token decline)",
			message:     "Insufficient Funds",
			wantErr:     nil,
			wantMatched: false,
		},
		{
			name:        "Validation error is not a decline (the given data was invalid.)",
			message:     "The given data was invalid.",
			wantErr:     nil,
			wantMatched: false,
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
			gotErr, gotMatched := ClassifyCardTokenDecline(tt.message)

			if gotMatched != tt.wantMatched {
				t.Errorf("matched = %v, want %v", gotMatched, tt.wantMatched)
			}

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("err = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
