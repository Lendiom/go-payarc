package customers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/utils"
)

func (s *Service) Update(id string, input CustomerInput) (*payarc.Customer, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", s.client.Url.String(), id), strings.NewReader(data.Encode()))
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

	var customer payarc.CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Data, err
}

func (s *Service) UpdateDefaultCard(customerID, defaultCardID string) error {
	payload := strings.NewReader(fmt.Sprintf("default_card_id=%s", defaultCardID))
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", s.client.Url.String(), customerID), payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
