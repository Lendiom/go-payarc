package payarc

import "errors"

var (
	ErrInsufficientFunds             = errors.New("insufficient funds, charge failed")
	ErrInvalidFromAccount            = errors.New("invalid from account, charge failed") //ErrInvalidFromAccount is similar to "Do Not Honor"
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
)

type RequestErrorErrors map[string][]string

type RequestError struct {
	Message string             `json:"message"`
	Error   string             `json:"error"`
	Errors  RequestErrorErrors `json:"errors,omitempty"`
}
