package charges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/utils"
)

func (s *Service) Create(input ChargeInput) (*payarc.Charge, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.client.Url.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	if r.StatusCode > http.StatusIMUsed || r.StatusCode < http.StatusOK {
		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to create the charge: %+v", errMsg)

		return nil, fmt.Errorf("create charge failed: %s", errMsg.Message)
	}

	var res payarc.ChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}
