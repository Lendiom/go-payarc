package banks

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

func (s *Service) Create(input CreateBankAccountInput) (*payarc.BankAccountCreated, error) {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	if r.StatusCode > http.StatusIMUsed || r.StatusCode < http.StatusOK {
		log.Println("Failed to create a bank account. Result is:")
		log.Println(string(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to create the bank account: %+v", errMsg)
		msg := errMsg.Message
		if errMsg.Error != "" {
			msg = errMsg.Error
		}

		return nil, fmt.Errorf("create bank account failed: %s", msg)
	}

	var res payarc.BankAccountCreatedResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		log.Println("Failed to decode the bank account create response:", err.Error())
		log.Println(string(body))

		return nil, err
	}

	return &res.BankAccount, nil
}

func (s *Service) Delete(bankID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", s.client.Url.String(), bankID), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
