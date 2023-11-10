package cards

import (
	"encoding/json"
	"fmt"
	"go-payarc/utils"
	"net/http"
	"path"
	"strings"
)

func (s *CardService) Create(id string, input TokenInput) (*Card, error) {
	token, err := s.createToken(input)
	if err != nil {
		fmt.Println(err)
	}

	data := utils.GenerateFormPayload(token)

	s.client.Url.Path = path.Join(s.client.Url.Path, id)

	fmt.Println(s.client.Url.String())
	fmt.Println(strings.NewReader(data.Encode()))

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

	return &cardData.Cards, nil
}

func (s *CardService) createToken(input TokenInput) (*Token, error) {
	data := utils.GenerateFormPayload(input)

	url := *s.client.Url
	url.Path = "v1/tokens"
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(data.Encode()))
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
	var tokenData TokenData
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	return &tokenData.Token, nil
}
