package cards

import (
	"go-payarc/client"
	"go-payarc/tokens"
)

type CardService struct {
	client *client.Client

	tokenService *tokens.TokenService
}

func NewCardService(apiKey string) (*CardService, error) {
	client, err := client.NewClient(apiKey)

	if err != nil {
		panic(err)
	}

	client.Url.Path = "v1/cards"

	tokenService, err := tokens.NewTokenService(apiKey)
	if err != nil {
		panic(err)
	}
	return &CardService{
		client:       client,
		tokenService: tokenService,
	}, nil
}
