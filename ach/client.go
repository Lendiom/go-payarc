package ach

import (
	"github.com/Lendiom/go-payarc/client"
)

type Service struct {
	client client.Client
}

func NewChargeService(apiKey string, environment client.PayArcEnvironment) (*Service, error) {
	client, err := client.NewClient(apiKey, environment)
	if err != nil {
		return nil, err
	}

	client.Url.Path = "v1/achcharges"

	return &Service{
		client: client,
	}, nil
}
