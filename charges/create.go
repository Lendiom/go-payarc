package charges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/utils"
)

func (s *Service) Create(input ChargeInput) (*ChargeResult, error) {
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

	r, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode > http.StatusIMUsed || r.StatusCode < http.StatusOK {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Failed to create a charge. Result is:")
		log.Println(string(body))

		r.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to create the charge: %+v", errMsg)

		switch strings.ToLower(errMsg.Message) {
		case "invalid card":
			return nil, payarc.ErrInvalidCard
		case "insufficient funds":
			return nil, payarc.ErrInsufficientFunds
		case "suspected fraud":
			return nil, payarc.ErrSuspectedFraud
		case "do not honor":
			return nil, payarc.ErrDoNotHonor
		case "suspected card":
			return nil, payarc.ErrSuspectedCard
		case "invalid from account":
			return nil, payarc.ErrInvalidFromAccount
		case "withdrawal limit exceeded":
			return nil, payarc.ErrWithdrawalLimitExceeded
		}

		return nil, fmt.Errorf("create charge failed: %s", errMsg.Message)
	}

	var res CreateChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}
