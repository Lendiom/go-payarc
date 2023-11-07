package tokens

import (
	"encoding/json"
	"fmt"
	"go-payarc/utils"
	"net/http"
	"strings"
)

func (s *TokenService) Create(input TokenInput) (*Token, error) {
	data := utils.GenerateFormPayload(input)

	req, err := http.NewRequest("POST", s.Client.Url.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.Client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var tokenData TokenData
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	return &tokenData.Token[0], nil
}
