package charges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) GetAll() ([]payarc.Charge, error) {
	req, err := http.NewRequest(http.MethodGet, s.client.Url.String(), nil)
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

	var res payarc.ChargesResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Charges, nil
}

func (s *Service) GetByID(id string) (*payarc.Charge, error) {
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

	var res payarc.ChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}
