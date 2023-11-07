package charges

import (
	"go-payarc/client"
)

type ChargeService struct {
	client *client.Client

	//tokenservice
}

func NewChargeService(apiKey string) (*ChargeService, error) {
	client, err := client.NewClient(apiKey)

	if err != nil {
		panic(err)
	}

	client.Url.Path = "v1/charges"

	return &ChargeService{
		client: client,
	}, nil
}
