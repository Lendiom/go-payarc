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
	ErrReferToIssuer                 = errors.New("issuing bank declined and asked the cardholder to contact them")
	// ErrReEnterTransaction is ISO 8583 host response code 19 (PayArc failure
	// code D0092 — "Re-enter transaction"). The card network reported a
	// transient processor/network glitch — the charge did NOT post, the card
	// itself is fine, and a retry typically succeeds. Distinct sentinel so
	// callers can surface a "try again" message and avoid disabling the
	// payment method or auto-draft for what is effectively a network hiccup.
	ErrReEnterTransaction = errors.New("payarc asked us to re-enter the transaction")
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
	case "refer to issuer":
		// Standard ISO 8583 decline surfaced by PayArc from host response
		// codes like D2001 ("Refer to card issuers special conditions").
		// The issuing bank declined and wants the cardholder to contact
		// them — there is nothing the merchant or processor can do to make
		// the charge succeed. Functionally close to ErrDoNotHonor but a
		// distinct host code, so keep it as its own sentinel so callers
		// can surface a tailored "contact your card issuer" message.
		return ErrReferToIssuer, true
	case "re-enter transaction":
		// ISO 8583 host response code 19 (PayArc failure code D0092). A
		// transient processor/network glitch — the charge did NOT post,
		// the card itself is not at fault, and a retry typically succeeds.
		// Not strictly a "bank decline," but it IS an expected outcome of
		// attempting a payment (not a server bug), so it belongs at WARN
		// alongside the other declines rather than paging on-call.
		return ErrReEnterTransaction, true
	}

	return nil, false
}

// ClassifyCardTokenDecline maps a known PayArc card-tokenization decline
// message to its sentinel error. It returns (sentinel, true) when PayArc
// rejected the card at tokenization for a reason the cardholder/bank owns
// (invalid card, CVV2 failure, do-not-honor, suspected fraud, etc.) and
// (nil, false) otherwise.
//
// Mirrors ClassifyChargeDecline but for the token-create code path: the set
// of messages PayArc returns during tokenization overlaps with — but is not
// identical to — the charge-decline set. CVV/CVV2 issues, for example, only
// show up at token time, while "insufficient funds" only shows up at charge
// time. Keeping them as separate classifiers keeps the WARN/ERROR decision
// honest at each call site.
//
// Callers use the boolean to decide log level: known declines are expected
// outcomes of adding a card and should be logged at WARN, while unknown
// messages, parse failures, and server-side errors should keep ERROR-level
// logging so they still surface in alerting dashboards.
//
// Matching is case-insensitive against the PayArc `message` field.
func ClassifyCardTokenDecline(message string) (error, bool) {
	switch strings.ToLower(message) {
	case "invalid card":
		return ErrInvalidCard, true
	case "invalid cvv":
		return ErrInvalidCCV, true
	case "cvv2 verification failed":
		return ErrCVV2Failed, true
	case "suspected fraud":
		return ErrSuspectedFraud, true
	case "do not honor":
		return ErrDoNotHonor, true
	case "suspected card":
		return ErrSuspectedCard, true
	case "expired card":
		// PayArc returns "Expired Card" (HTTP 409) from the token-create
		// endpoint when callers pass AuthorizeCard=true and the $0 auth is
		// declined because the card's expiration date has passed. It's a
		// cardholder-owned condition — the customer needs to use a
		// different card — not a server bug, so it belongs at WARN
		// alongside the other known token declines.
		return ErrExpiredCard, true
	}

	return nil, false
}
