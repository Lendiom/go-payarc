package customers

import (
	"fmt"
	"net/http"
	"path"
)

func (s *CustomerService) Delete(id string) error {
	s.client.Url.Path = path.Join(s.client.Url.Path, id)

	req, err := http.NewRequest("DELETE", s.client.Url.String(), nil)
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
