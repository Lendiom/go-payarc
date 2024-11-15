package deposits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lendiom/go-payarc"
)

func (s *Service) GetReport(month time.Month, year uint) (*payarc.DepositReport, error) {
	//https://${ baseUrl }/v1/batch/reports?month=7&year=2024
	url := fmt.Sprintf("%s/batch/reports?month=%d&year=%d", s.client.Url.String(), month, year)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var report payarc.DepositReport
	if err := json.NewDecoder(res.Body).Decode(&report); err != nil {
		return nil, err
	}

	return &report, nil
}

func (s *Service) GetReportDetails(batchReferenceNumber string) (*payarc.DepositBatchDetailsResponse, error) {
	//https://${ baseUrl }/v1/deposit/reports/details?reference_number=xxx
	url := fmt.Sprintf("%s/deposit/reports/details?reference_number=%s", s.client.Url.String(), batchReferenceNumber)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.client.ApiKey))
	req.Header.Add("Accept", "application/json")

	res, err := s.client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result payarc.DepositBatchDetailsResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
