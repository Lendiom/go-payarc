package charges

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/client"
)

// TestVoid_BlankChargeID is the request half of the production failure. With an empty id the
// path collapsed to v1/charges//void, which PayArc answered "The route v1/charges//void could
// not be found." Nothing should reach the network for an id we can already see is unusable.
func TestVoid_BlankChargeID(t *testing.T) {
	tests := []struct {
		name     string
		chargeID string
	}{
		{name: "empty", chargeID: ""},
		{name: "whitespace only", chargeID: "  "},
		{name: "contains a slash", chargeID: "ch_1/void"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()

			u, err := url.Parse(srv.URL)
			if err != nil {
				t.Fatalf("parse httptest URL: %v", err)
			}

			svc := &Service{client: client.Client{ApiKey: "test-key", HttpClient: *srv.Client(), Url: *u}}

			charge, voidErr := svc.Void(tt.chargeID, VoidInput{})

			if !errors.Is(voidErr, payarc.ErrMissingParameter) {
				t.Errorf("Void error = %v, want it to match payarc.ErrMissingParameter", voidErr)
			}

			if charge != nil {
				t.Errorf("Void returned charge %+v, want nil", charge)
			}

			if called {
				t.Error("Void sent a request; it should have been refused before reaching the network")
			}
		})
	}
}
