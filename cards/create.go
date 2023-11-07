package cards

import (
	"encoding/json"
	"fmt"
	"go-payarc/tokens"
	"go-payarc/utils"
	"net/http"
	"path"
	"strings"
)

func (s *CardService) Create(id string, input tokens.TokenInput) (*Card, error) {
	token, err := s.tokenService.Create(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(token)
	}

	fmt.Println(token)

	data := utils.GenerateFormPayload(token)

	s.client.Url.Path = path.Join(s.client.Url.Path, id)

	req, err := http.NewRequest("PATCH", s.client.Url.String(), strings.NewReader(data.Encode()))
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
	var cardData CardData
	if err := json.NewDecoder(res.Body).Decode(&cardData); err != nil {
		return nil, err
	}

	return &cardData.Cards[0], nil
}
