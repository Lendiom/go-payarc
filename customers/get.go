package customers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *CustomerService) GetAll() ([]Customer, error) {
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
	var customerData CustomerLookup
	if err := json.NewDecoder(res.Body).Decode(&customerData); err != nil {
		println("here")
		return nil, err
	}

	return customerData.Customers, nil
}

func (s *CustomerService) GetById(id string) (*Customer, error) {
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
	var customerData SingleCustomerLookup
	if err := json.NewDecoder(res.Body).Decode(&customerData); err != nil {
		println("here")
		return nil, err
	}

	return &customerData.Customer, nil
}
