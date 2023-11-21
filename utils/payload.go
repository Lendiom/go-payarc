package utils

import (
	"net/url"

	"github.com/go-playground/form"
)

func GenerateFormPayload(s interface{}) (url.Values, error) {
	encoder := form.NewEncoder()

	return encoder.Encode(&s)
}
