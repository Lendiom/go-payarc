package ach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) GetAll(limit, page uint) (int, []payarc.ACHCharge, error) {
	reqUrl := fmt.Sprintf("%s?include=customer&limit=%d&page=%d", s.client.Url.String(), limit, page)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	r, err := s.client.HttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer r.Body.Close()

	var res payarc.ACHChargesResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return 0, nil, err
	}

	return res.Metadata.Pagination.Total, res.Charges, nil
}

func (s *Service) GetByID(id string) (*payarc.ACHCharge, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", s.client.Url.String(), id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	r, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Failed to get a charge. Result is:")
		log.Println(string(body))

		r.Body = io.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to get the charge: %+v", errMsg)

		if errMsg.Error != "" {
			return nil, fmt.Errorf("failed to get the charge: %s", errMsg.Error)
		}

		return nil, fmt.Errorf("failed to get the charge: %s", errMsg.Message)
	}

	var res payarc.ACHChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}
