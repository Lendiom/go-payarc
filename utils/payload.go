package utils

import (
	"net/url"

	"github.com/go-playground/form"
)

func GenerateFormPayload(val interface{}) (url.Values, error) {
	encoder := form.NewEncoder()

	return encoder.Encode(&val)
}
