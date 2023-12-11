package payarc

type RequestErrorErrors map[string][]string

type RequestError struct {
	Message string             `json:"message"`
	Errors  RequestErrorErrors `json:"errors,omitempty"`
}
