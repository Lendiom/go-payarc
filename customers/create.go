package customers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/utils"
)

var (
	ErrInvalidExpirationMonth = errors.New("card expiration month must be a two digit number")
	ErrInvalidExpirationYear  = errors.New("card expiration year must be a four digit number")
	ErrInvalidCardNumber      = errors.New("card number must be 14 to 19 digits long")
	ErrInvalidCardSource      = errors.New("card source is invalid")
)

func (s *Service) Create(input CustomerInput) (*payarc.Customer, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.client.Url.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusIMUsed || res.StatusCode < http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			slog.Error("payarc customer create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "customers.create"),
				slog.Int("status_code", res.StatusCode),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, err
		}

		slog.Error("payarc customer create failed",
			slog.String("component", "go-payarc"),
			slog.String("op", "customers.create"),
			slog.Int("status_code", res.StatusCode),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		)

		return nil, fmt.Errorf("create customer failed: %s OR %s", errMsg.Message, errMsg.Error)
	}

	var customer payarc.CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Data, err
}

func (s *Service) CreateCard(id string, input TokenInput) (*payarc.Customer, *payarc.Card, error) {
	token, err := s.createToken(input)
	if err != nil {
		return nil, nil, err
	}

	data, err := utils.GenerateFormPayload(token)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", s.client.Url.String(), id), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusIMUsed || res.StatusCode < http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil, err
		}

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			slog.Error("payarc customer card create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "customers.create_card"),
				slog.String("customer_id", id),
				slog.Int("status_code", res.StatusCode),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, nil, err
		}

		slog.Error("payarc customer card create failed",
			slog.String("component", "go-payarc"),
			slog.String("op", "customers.create_card"),
			slog.String("customer_id", id),
			slog.Int("status_code", res.StatusCode),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		)

		return nil, nil, fmt.Errorf("create card failed: %s", errMsg.Message)
	}

	var customer payarc.CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, nil, err
	}

	return &customer.Data, &token.Card.Data, nil
}

func (s *Service) createToken(input TokenInput) (*Token, error) {
	if len(input.ExpMonth) != 2 {
		return nil, ErrInvalidExpirationMonth
	}

	if len(input.ExpYear) != 4 {
		return nil, ErrInvalidExpirationYear
	}

	if cardLen := len(input.CardNumber); cardLen > 19 || cardLen < 14 {
		return nil, ErrInvalidCardNumber
	}

	if !input.CardSource.Valid() {
		return nil, ErrInvalidCardSource
	}

	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	url := s.client.Url
	url.Path = "v1/tokens"

	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusIMUsed || res.StatusCode < http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			slog.Error("payarc token create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "customers.create_token"),
				slog.Int("status_code", res.StatusCode),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, err
		}

		// Business declines (do-not-honor, CVV2 failure, suspected fraud,
		// etc.) are the expected outcome of attempting to tokenize a card,
		// not server failures — log them at WARN so they stop polluting
		// ERROR dashboards and on-call channels. Mirrors the WARN/ERROR
		// split already in place for charges/create.go.
		sentinel, declined := payarc.ClassifyCardTokenDecline(errMsg.Message)

		attrs := []any{
			slog.String("component", "go-payarc"),
			slog.String("op", "customers.create_token"),
			slog.Int("status_code", res.StatusCode),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		}

		if declined {
			slog.Warn("payarc declined card token", attrs...)
			return nil, sentinel
		}

		// "The given data was invalid." with a card_holder_name error is a
		// customer-data problem — PayArc's Laravel-style format check rejects
		// digits and characters outside its accepted set for the cardholder
		// name. That is not a server bug: the customer typed something the
		// processor won't accept, and the caller can surface a "fix your
		// name" message. Log at WARN so this expected outcome stops paging
		// on-call, but keep ERROR (below) for the same "The given data was
		// invalid." shape combined with other/additional fields, since those
		// usually mean the caller sent a malformed request payload (missing
		// card_number, wrong type, etc.) that a human should investigate.
		if strings.EqualFold(errMsg.Message, "the given data was invalid.") &&
			len(errMsg.Errors) == 1 &&
			len(errMsg.Errors["card_holder_name"]) > 0 {
			slog.Warn("payarc rejected card_holder_name format", attrs...)
			return nil, errors.New(errMsg.Errors["card_holder_name"][0])
		}

		slog.Error("payarc token create failed", attrs...)

		// Other "The given data was invalid." shapes (multiple fields, or a
		// non-card_holder_name field) still surface the per-field message so
		// the caller has something actionable to show, but the ERROR log
		// above keeps alerting on-call because it likely means our request
		// payload is broken.
		if strings.EqualFold(errMsg.Message, "the given data was invalid.") {
			if cardHolderNameErrors, ok := errMsg.Errors["card_holder_name"]; ok && len(cardHolderNameErrors) > 0 {
				return nil, errors.New(cardHolderNameErrors[0])
			}
		}

		return nil, errors.New(errMsg.Message)
	}

	var tokenData TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	return &tokenData.Data, nil
}
