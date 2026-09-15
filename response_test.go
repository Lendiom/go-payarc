package payarc

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func responseWith(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestCheckResponse(t *testing.T) {
	type args struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name string
		args args

		wantErr      bool
		wantNotFound bool
		wantMessage  string
		wantDetail   string
		wantBody     string
		wantErrText  string
	}{
		{
			name: "a 200 passes and leaves the body readable",
			args: args{statusCode: http.StatusOK, body: `{"data":{"id":"chg_1"}}`},
		},
		{
			name: "a 201 passes",
			args: args{statusCode: http.StatusCreated, body: `{"data":{"id":"chg_1"}}`},
		},
		{
			// The case the whole helper exists for: PayArc answers an unknown charge id this
			// way, and the body decodes cleanly into the success type as a zero-valued Charge.
			name:         "a 404 with a payarc message is reported as not found",
			args:         args{statusCode: http.StatusNotFound, body: `{"message":"The route v1/charges/nope could not be found."}`},
			wantErr:      true,
			wantNotFound: true,
			wantMessage:  "The route v1/charges/nope could not be found.",
			wantErrText:  "get charge failed with status 404: The route v1/charges/nope could not be found.",
		},
		{
			name:        "a 401 carries payarc's message but is not a not-found",
			args:        args{statusCode: http.StatusUnauthorized, body: `{"message":"Unauthenticated."}`},
			wantErr:     true,
			wantMessage: "Unauthenticated.",
			wantErrText: "get charge failed with status 401: Unauthenticated.",
		},
		{
			name:        "the error field is used when message is empty",
			args:        args{statusCode: http.StatusBadRequest, body: `{"error":"invalid parameters"}`},
			wantErr:     true,
			wantDetail:  "invalid parameters",
			wantErrText: "get charge failed with status 400: invalid parameters",
		},
		{
			// A gateway or CDN page rather than PayArc itself; keeping it separates an
			// upstream outage from a PayArc business rejection.
			name:        "a non-payarc body is kept raw",
			args:        args{statusCode: http.StatusBadGateway, body: "<html>502 Bad Gateway</html>"},
			wantErr:     true,
			wantBody:    "<html>502 Bad Gateway</html>",
			wantErrText: "get charge failed with status 502: <html>502 Bad Gateway</html>",
		},
		{
			name:        "an empty error body still names the status",
			args:        args{statusCode: http.StatusInternalServerError, body: `{}`},
			wantErr:     true,
			wantErrText: "get charge failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := responseWith(tt.args.statusCode, tt.args.body)

			err := CheckResponse(resp, "get charge")

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("CheckResponse returned error: %v", err)
				}

				// The success path must be able to decode from the body as if nothing read it.
				got, readErr := io.ReadAll(resp.Body)
				if readErr != nil {
					t.Fatalf("read body after CheckResponse: %v", readErr)
				}

				if string(got) != tt.args.body {
					t.Errorf("body after CheckResponse = %q, want %q", got, tt.args.body)
				}

				return
			}

			if err == nil {
				t.Fatal("CheckResponse returned nil, want an error")
			}

			if err.Error() != tt.wantErrText {
				t.Errorf("error text = %q, want %q", err.Error(), tt.wantErrText)
			}

			if errors.Is(err, ErrNotFound) != tt.wantNotFound {
				t.Errorf("errors.Is(err, ErrNotFound) = %v, want %v", errors.Is(err, ErrNotFound), tt.wantNotFound)
			}

			var statusErr *StatusError
			if !errors.As(err, &statusErr) {
				t.Fatalf("error is %T, want *StatusError", err)
			}

			if statusErr.StatusCode != tt.args.statusCode {
				t.Errorf("StatusCode = %d, want %d", statusErr.StatusCode, tt.args.statusCode)
			}

			if statusErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", statusErr.Message, tt.wantMessage)
			}

			if statusErr.Detail != tt.wantDetail {
				t.Errorf("Detail = %q, want %q", statusErr.Detail, tt.wantDetail)
			}

			if statusErr.Body != tt.wantBody {
				t.Errorf("Body = %q, want %q", statusErr.Body, tt.wantBody)
			}
		})
	}
}

// TestCheckResponse_FieldErrors keeps PayArc's per-field validation errors reachable, since
// they are the only place the specific rejected field is named.
func TestCheckResponse_FieldErrors(t *testing.T) {
	resp := responseWith(http.StatusBadRequest, `{"message":"the given data was invalid.","errors":{"card_holder_name":["The card holder name format is invalid."]}}`)

	err := CheckResponse(resp, "create card")

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("error is %T, want *StatusError", err)
	}

	got := statusErr.Fields["card_holder_name"]
	if len(got) != 1 || got[0] != "The card holder name format is invalid." {
		t.Errorf("Fields[card_holder_name] = %v, want the format message", got)
	}
}
