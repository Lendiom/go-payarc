package payarc

import (
	"errors"
	"strings"
)

var (
	ErrInsufficientFunds             = errors.New("insufficient funds, charge failed")
	ErrInvalidFromAccount            = errors.New("invalid from account, charge failed") //ErrInvalidFromAccount is similar to "Do Not Honor"
	ErrGeneralCardAuthDecline        = errors.New("card authorization failed, please contact your bank")
	ErrSuspectedFraud                = errors.New("bank suspects fraud")
	ErrDoNotHonor                    = errors.New("bank said to not honor")
	ErrSuspectedCard                 = errors.New("bank suspects card")
	ErrInvalidCard                   = errors.New("invalid card")
	ErrInvalidCCV                    = errors.New("invalid ccv")
	ErrExpiredCard                   = errors.New("expired card")
	ErrCVV2Failed                    = errors.New("cvv2 verification failed")
	ErrUnauthorizedSECType           = errors.New("unauthorized sec type")
	ErrInvalidData                   = errors.New("invalid data")
	ErrWithdrawalLimitExceeded       = errors.New("withdrawal limit exceeded")
	ErrCustomerRequestedStopPayments = errors.New("customer requested stop payments for this seller")
	ErrClosedAccount                 = errors.New("account is closed")
)

type RequestErrorErrors map[string][]string

type RequestError struct {
	Message string             `json:"message"`
	Error   string             `json:"error"`
	Errors  RequestErrorErrors `json:"errors,omitempty"`
}

// ClassifyChargeDecline maps a known PayArc charge-decline message to its
// sentinel error. It returns (sentinel, true) when the message is a
// recognized business decline the cardholder/bank owns (insufficient funds,
// expired card, do-not-honor, etc.) and (nil, false) otherwise.
//
// Callers use the boolean to decide log level: business declines are expected
// outcomes of attempting a payment and should be logged at WARN, while
// unknown messages, parse failures, and server-side errors should keep
// ERROR-level logging so they still surface in alerting dashboards.
//
// Matching is case-insensitive against the PayArc `message` field.
func ClassifyChargeDecline(message string) (error, bool) {
	switch strings.ToLower(message) {
	case "invalid card":
		return ErrInvalidCard, true
	case "insufficient funds":
		return ErrInsufficientFunds, true
	case "suspected fraud":
		return ErrSuspectedFraud, true
	case "do not honor":
		return ErrDoNotHonor, true
	case "suspected card":
		return ErrSuspectedCard, true
	case "invalid from account":
		return ErrInvalidFromAccount, true
	case "withdrawal limit exceeded":
		return ErrWithdrawalLimitExceeded, true
	case "customer requested stop of all recurring payments from specific merchant":
		return ErrCustomerRequestedStopPayments, true
	case "expired card":
		return ErrExpiredCard, true
	case "general cardauth decline":
		return ErrGeneralCardAuthDecline, true
	case "closed account":
		return ErrClosedAccount, true
	}

	return nil, false
}
