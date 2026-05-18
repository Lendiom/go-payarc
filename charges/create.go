package charges

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

		r.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			slog.Error("payarc charge create failed; response body did not parse",
				slog.String("component", "go-payarc"),
				slog.String("op", "charges.create"),
				slog.Int("status_code", r.StatusCode),
				slog.String("body", string(body)),
				slog.Any("error", err),
			)
			return nil, err
		}

		// Business declines (insufficient funds, expired card, do-not-honor,
		// etc.) are the expected outcome of attempting a payment, not server
		// failures — log them at WARN so they stop polluting ERROR dashboards
		// and on-call channels. Anything we don't recognize stays at ERROR.
		sentinel, declined := payarc.ClassifyChargeDecline(errMsg.Message)

		attrs := []any{
			slog.String("component", "go-payarc"),
			slog.String("op", "charges.create"),
			slog.Int("status_code", r.StatusCode),
			slog.String("body", string(body)),
			slog.String("payarc_message", errMsg.Message),
			slog.String("payarc_error", errMsg.Error),
			slog.Any("payarc_field_errors", errMsg.Errors),
		}

		if declined {
			slog.Warn("payarc declined charge", attrs...)
			return nil, sentinel
		}

		slog.Error("payarc charge create failed", attrs...)

		return nil, fmt.Errorf("create charge failed: %s", errMsg.Message)
	}

	var res CreateChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}
