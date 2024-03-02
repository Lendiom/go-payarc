package customers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) GetAll(limit, page uint) ([]payarc.Customer, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?limit=%d&page=%d", s.client.Url.String(), limit, page), nil)
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

	var customers CustomersResponse
	if err := json.NewDecoder(res.Body).Decode(&customers); err != nil {
		return nil, err
	}

	return customers.Customers, nil
}

func (s *Service) GetByID(id string) (*payarc.Customer, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", s.client.Url.String(), id), nil)
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

	var customer CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Customer, nil
}
