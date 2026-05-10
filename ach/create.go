package ach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) Create(input CreateAchChargeInput) (*ACHChargeResult, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.client.Url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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
			slog.Error("payarc ach charge create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "ach.create"),
				slog.Int("status_code", r.StatusCode),
				slog.String("payload", string(data)),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, err
		}

		slog.Error("payarc ach charge create failed",
			slog.String("component", "go-payarc"),
			slog.String("op", "ach.create"),
			slog.Int("status_code", r.StatusCode),
			slog.String("payload", string(data)),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		)

		switch errMsg.Message {
		case "Unauthorized SEC type":
			return nil, payarc.ErrUnauthorizedSECType
		}

		return nil, fmt.Errorf("create charge failed: %s", errMsg.Message)
	}

	var res CreateACHChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		slog.Error("payarc ach charge create succeeded; response body did not parse",
			slog.String("component", "go-payarc"),
			slog.String("op", "ach.create"),
			slog.String("body", string(body)),
			slog.Any("error", err),
		)

		return nil, err
	}

	return &res.Charge, nil
}
