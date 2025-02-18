package ach

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

func (s *Service) Create(input CreateAchChargeInput) (*ACHChargeResult, error) {
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
		log.Println("Payload for ach charge creation is:")
		log.Println(data.Encode())
		log.Println("Failed to create a charge. Result is:")
		log.Println(string(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to create the charge: %+v", errMsg)

		switch errMsg.Message {
		case "Unauthorized SEC type":
			return nil, payarc.ErrUnauthorizedSECType
		}

		return nil, fmt.Errorf("create charge failed: %s", errMsg.Message)
	}

	var res CreateACHChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		log.Println("failed to decode response body. the response was:")
		log.Println(string(body))

		return nil, err
	}

	return &res.Charge, nil
}
