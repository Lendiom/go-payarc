package customers

import (
	"go-payarc/client"
)

type CustomerService struct {
	client *client.Client
}

func NewCustomerService(apiKey string) (*CustomerService, error) {
	client, err := client.NewClient(apiKey)

	if err != nil {
		panic(err)
	}

	client.Url.Path = "v1/customers"

	return &CustomerService{
		client: client,
	}, nil
}
