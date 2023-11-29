package customers

import (
	"github.com/Lendiom/go-payarc/client"
)

type CustomerService struct {
	client *client.Client
}

func NewCustomerService(apiKey string, environment client.PayArcEnvironment) (*CustomerService, error) {
	client, err := client.NewClient(apiKey, environment)
	if err != nil {
		return nil, err
	}

	client.Url.Path = "v1/customers"

	return &CustomerService{
		client: client,
	}, nil
}
