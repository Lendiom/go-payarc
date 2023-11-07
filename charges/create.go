package charges

import (
	"encoding/json"
	"fmt"
	"go-payarc/utils"
	"net/http"
	"strings"
)

func (s *ChargeService) Create(input ChargeInput) (*Charge, error) {
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
	var chargeData ChargeData
	if err := json.NewDecoder(res.Body).Decode(&chargeData); err != nil {
		return nil, err
	}

	return &chargeData.Charges[0], nil
}
