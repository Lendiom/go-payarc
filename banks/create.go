package banks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			slog.Error("payarc bank account create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "banks.create"),
				slog.Int("status_code", r.StatusCode),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, err
		}

		slog.Error("payarc bank account create failed",
			slog.String("component", "go-payarc"),
			slog.String("op", "banks.create"),
			slog.Int("status_code", r.StatusCode),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		)

		msg := errMsg.Message
		if errMsg.Error != "" {
			msg = errMsg.Error
		}

		return nil, fmt.Errorf("create bank account failed: %s", msg)
	}

	var res payarc.BankAccountCreatedResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		slog.Error("payarc bank account create succeeded; response body did not parse",
			slog.String("component", "go-payarc"),
			slog.String("op", "banks.create"),
			slog.String("body", string(body)),
			slog.Any("error", err),
		)

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
