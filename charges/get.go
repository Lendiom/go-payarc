package charges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *ChargeService) GetAll() ([]ChargeData, error) {
	req, err := http.NewRequest("GET", s.client.Url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var charges Charges
	if err := json.NewDecoder(res.Body).Decode(&charges); err != nil {
		return nil, err
	}

	return charges.Charges, nil
}

func (s *ChargeService) GetByID(id string) (*ChargeData, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.client.Url.String(), id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var charge Charge
	if err := json.NewDecoder(res.Body).Decode(&charge); err != nil {
		return nil, err
	}

	return &charge.Charge, nil
}
