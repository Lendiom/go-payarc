package utils

import (
	"log"
	"net/url"

	"github.com/go-playground/form"
)

func GenerateFormPayload(s interface{}) url.Values {
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&s)
	if err != nil {
		log.Panic(err)
	}

	return values
}
