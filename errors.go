package payarc

type RequestErrorErrors map[string][]string

type RequestError struct {
	Message string             `json:"message"`
	Error   string             `json:"error"`
	Errors  RequestErrorErrors `json:"errors,omitempty"`
}
