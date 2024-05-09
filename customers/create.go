package customers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
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

		log.Println("Failed to create a customer. Result is:")
		log.Println(string(body))

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("create customer failed: %s OR %s", errMsg.Message, errMsg.Error)
	}

	var customer CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Customer, err
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

	s.client.Url.Path = path.Join(s.client.Url.Path, id)

	req, err := http.NewRequest(http.MethodPatch, s.client.Url.String(), strings.NewReader(data.Encode()))
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

		log.Println("Failed to create a card. Result is:")
		log.Println(string(body))

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return nil, nil, err
		}

		return nil, nil, fmt.Errorf("create card failed: %s", errMsg.Message)
	}

	var customer CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, nil, err
	}

	return &customer.Customer, &token.Card.Data, nil
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

		log.Println("Failed to create a token. Result is:")
		log.Println(string(body))

		res.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		switch strings.ToLower(errMsg.Message) {
		case "invalid card":
			return nil, payarc.ErrInvalidCard
		case "invalid cvv":
			return nil, payarc.ErrInvalidCCV
		case "cvv2 verification failed":
			return nil, payarc.ErrCVV2Failed
		case "suspected fraud":
			return nil, payarc.ErrSuspectedFraud
		case "do not honor":
			return nil, payarc.ErrDoNotHonor
		case "suspected card":
			return nil, payarc.ErrSuspectedCard
		case "the given data was invalid.":
			// This error usually has more details in the "errors" field
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
