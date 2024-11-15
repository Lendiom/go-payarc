package deposits

import "github.com/Lendiom/go-payarc/client"

type Service struct {
	client client.Client
}

func NewDepositService(apiKey string, environment client.PayArcEnvironment) (*Service, error) {
	client, err := client.NewClient(apiKey, environment)
	if err != nil {
		return nil, err
	}

	client.Url.Path = "v1" //Sadly, PayArc does not have a dedicated endpoint, as we use several for deposits

	return &Service{
		client: client,
	}, nil
}
