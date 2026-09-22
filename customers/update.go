package customers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/utils"
)

func (s *Service) Update(id string, input CustomerInput) (*payarc.Customer, error) {
	if err := payarc.RequireParam("customer id", id); err != nil {
		return nil, err
	}

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

	if err := payarc.CheckResponse(res, "update customer"); err != nil {
		return nil, err
	}

	var customer payarc.CustomerResponse
	if err := json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer.Data, nil
}

func (s *Service) UpdateDefaultCard(customerID, defaultCardID string) error {
	if err := payarc.RequireParam("default card id", defaultCardID); err != nil {
		return err
	}

	return s.patchDefaultCard(customerID, defaultCardID, "set the default card")
}

// ClearDefaultCard unsets the customer's default card by sending an empty
// default_card_id.
//
// A customer whose chosen payment method is a bank account has no card to point
// at — default_card_id only accepts a card id — so without this the field keeps
// naming whichever card was default last, which is stale and misleading. There
// is no equivalent field for bank accounts, so clearing is the only way to say
// "no default card".
func (s *Service) ClearDefaultCard(customerID string) error {
	return s.patchDefaultCard(customerID, "", "clear the default card")
}

// patchDefaultCard PATCHes default_card_id on the customer. An empty value is
// meaningful here (it clears the field), so unlike the exported callers this
// does not require one; the caller decides whether empty is valid.
func (s *Service) patchDefaultCard(customerID, defaultCardID, action string) error {
	if err := payarc.RequireParam("customer id", customerID); err != nil {
		return err
	}

	payload := strings.NewReader(url.Values{"default_card_id": {defaultCardID}}.Encode())
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

	if err := payarc.CheckResponse(res, action); err != nil {
		return err
	}

	return nil
}
