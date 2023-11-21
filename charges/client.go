package charges

import (
	"github.com/Lendiom/go-payarc/client"
)

type ChargeService struct {
	client *client.Client
}

func NewChargeService(apiKey string, environment client.PayArcEnvironment) (*ChargeService, error) {
	client, err := client.NewClient(apiKey, environment)
	if err != nil {
		return nil, err
	}

	client.Url.Path = "v1/charges"

	return &ChargeService{
		client: client,
	}, nil
}
