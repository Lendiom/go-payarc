package tokens

import "go-payarc/client"

type TokenService struct {
	Client *client.Client
}

func NewTokenService(apiKey string) (*TokenService, error) {
	client, err := client.NewClient(apiKey)

	if err != nil {
		panic(err)
	}

	client.Url.Path = "v1/tokens"

	return &TokenService{
		Client: client,
	}, nil
}
