package customers

import (
	"encoding/json"
	"fmt"
	"go-payarc/utils"
	"net/http"
	"path"
	"strings"
)

func (s *CustomerService) Update(id string, input CustomerInput) (*CustomerData, error) {
	data := utils.GenerateFormPayload(input)

	s.client.Url.Path = path.Join(s.client.Url.Path, id)
	req, err := http.NewRequest("PATCH", s.client.Url.String(), strings.NewReader(data.Encode()))
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
	var customer Customer
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Customer, err
}

func (s *CustomerService) UpdateDefaultCard(customerId, defaultCardId string) (*CustomerData, error) {
	s.client.Url.Path = path.Join(s.client.Url.Path, customerId)

	payload := strings.NewReader(fmt.Sprintf("default_card_id=%s", defaultCardId))
	req, err := http.NewRequest("PATCH", s.client.Url.String(), payload)
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
	var customer Customer
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Customer, err
}
