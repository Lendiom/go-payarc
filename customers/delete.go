package customers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) Delete(id string) error {
	if err := payarc.RequireParam("customer id", id); err != nil {
		return err
	}

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

	// A delete that finds nothing already has the outcome the caller wanted, so 404 is
	// success here rather than an error that would block cleaning up our own record of it.
	if err := payarc.CheckResponse(res, "delete customer"); err != nil && !errors.Is(err, payarc.ErrNotFound) {
		return err
	}

	return nil
}

func (s *Service) DeleteCard(customerID, cardID string) error {
	if err := payarc.RequireParam("customer id", customerID); err != nil {
		return err
	}

	if err := payarc.RequireParam("card id", cardID); err != nil {
		return err
	}

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

	// A delete that finds nothing already has the outcome the caller wanted, so 404 is
	// success here rather than an error that would block cleaning up our own record of it.
	if err := payarc.CheckResponse(res, "delete card"); err != nil && !errors.Is(err, payarc.ErrNotFound) {
		return err
	}

	return nil
}
