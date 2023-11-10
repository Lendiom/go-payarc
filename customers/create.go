package customers

import (
	"encoding/json"
	"fmt"
	"go-payarc/utils"
	"net/http"
	"path"
	"strings"
)

func (s *CustomerService) Create(input CustomerInput) (*Customer, error) {
	data := utils.GenerateFormPayload(input)

	req, err := http.NewRequest("POST", s.client.Url.String(), strings.NewReader(data.Encode()))
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
	var customerData SingleCustomerLookup
	if err := json.NewDecoder(res.Body).Decode(&customerData); err != nil {
		return nil, err
	}

	return &customerData.Customer, err
}

func (s *CustomerService) CreateCard(id string, input TokenInput) (*Customer, error) {
	token, err := s.createToken(input)
	if err != nil {
		return nil, err
	}

	data := utils.GenerateFormPayload(token)

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
	var customerData SingleCustomerLookup
	if err := json.NewDecoder(res.Body).Decode(&customerData); err != nil {
		return nil, err
	}

	return &customerData.Customer, nil
}

func (s *CustomerService) createToken(input TokenInput) (*Token, error) {
	data := utils.GenerateFormPayload(input)

	url := *s.client.Url
	url.Path = "v1/tokens"
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(data.Encode()))
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
	var tokenData TokenData
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	return &tokenData.Token, nil
}
