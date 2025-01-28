package customers

import (
	"fmt"
	"net/http"
)

func (s *Service) Delete(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", s.client.Url.String(), id), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (s *Service) DeleteCard(customerID, cardID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/cards/%s", s.client.Url.String(), customerID, cardID), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
