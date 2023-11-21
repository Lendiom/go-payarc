package charges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Lendiom/go-payarc/utils"
)

func (s *ChargeService) Create(input ChargeInput) (*ChargeData, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

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

	var charge Charge
	if err := json.NewDecoder(res.Body).Decode(&charge); err != nil {
		return nil, err
	}

	return &charge.Charge, nil
}
