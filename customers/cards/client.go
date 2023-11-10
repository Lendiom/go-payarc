package cards

import (
	"go-payarc/client"
)

type CardService struct {
	client client.Client
}

func NewCardService(apiKey string) (*CardService, error) {
	client, err := client.NewClient(apiKey)

	if err != nil {
		panic(err)
	}

	client.Url.Path = "v1/customers"

	return &CardService{
		client: *client,
	}, nil
}
