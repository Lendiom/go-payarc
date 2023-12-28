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

type VoidInput struct {
	Reason      payarc.RefundReason `form:"reason,omitempty"`
	Description string              `form:"void_description,omitempty"`
}

func (s *Service) Void(chargeID string, input VoidInput) (*payarc.Charge, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	url := s.client.Url
	url.Path = fmt.Sprintf("v1/charges/%s/void", chargeID)

	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(data.Encode()))
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

	if r.StatusCode > http.StatusIMUsed || r.StatusCode < http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Failed to void a charge. Result is:")
		log.Println(string(body))

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to void the charge: %+v", errMsg)

		return nil, fmt.Errorf("void charge failed: %s", errMsg.Message)
	}

	var res payarc.ChargeResponse
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Charge, nil
}

type RefundInput struct {
	Amount                   int64               `form:"amount,omitempty"`
	Reason                   payarc.RefundReason `form:"reason,omitempty"`
	Description              string              `form:"description,omitempty"`
	DoNotSendEmailToCustomer payarc.YesOrNo      `form:"do_not_send_email_to_customer"`
	DoNotSendSmsToCustomer   payarc.YesOrNo      `form:"do_not_send_sms_to_customer"`
}

func (s *Service) Refund(chargeID string, input RefundInput) (*payarc.Refund, error) {
	data, err := utils.GenerateFormPayload(input)
	if err != nil {
		return nil, err
	}

	url := s.client.Url
	url.Path = fmt.Sprintf("v1/charges/%s/refunds", chargeID)

	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(data.Encode()))
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

	if r.StatusCode > http.StatusIMUsed || r.StatusCode < http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Failed to create a refund. Result is:")
		log.Println(string(body))

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		var errMsg payarc.RequestError
		if err := json.NewDecoder(r.Body).Decode(&errMsg); err != nil {
			return nil, err
		}

		log.Printf("Failed to refund the charge: %+v", errMsg)

		return nil, fmt.Errorf("refund charge failed: %s", errMsg.Message)
	}

	var refundRes payarc.RefundResponse
	if err := json.NewDecoder(r.Body).Decode(&refundRes); err != nil {
		return nil, err
	}

	return &refundRes.Data, nil
}
