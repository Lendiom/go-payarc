package charges

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Lendiom/go-payarc/client"
)

// TestGetByID_RequestsTransactionMetadata verifies the detail endpoint is called
// with include=transaction_metadata and that the returned charge's metadata is
// parsed. PayArc omits transaction_metadata from the charge detail response
// unless it is explicitly included.
func TestGetByID_RequestsTransactionMetadata(t *testing.T) {
	const body = `{"data":{"object":"Charge","id":"chg_abc123","amount":35022,"status":"settled","captured":1,"transaction_metadata":{"data":[{"object":"ChargeTransactionMetadata","key":"loanId","value":"640747b60796a70001820d43"}]}}}`

	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query().Get("include")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse httptest URL: %v", err)
	}

	svc := &Service{
		client: client.Client{
			ApiKey:     "test-key",
			HttpClient: *srv.Client(),
			Url:        *u,
		},
	}

	charge, err := svc.GetByID("chg_abc123")
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if gotPath != "/chg_abc123" {
		t.Errorf("request path = %q, want %q", gotPath, "/chg_abc123")
	}

	if gotQuery != "transaction_metadata" {
		t.Errorf("include query = %q, want %q", gotQuery, "transaction_metadata")
	}

	if charge.TransactionMetadata == nil {
		t.Fatal("expected transaction_metadata to be parsed, got nil")
	}

	var loanID string
	for _, m := range charge.TransactionMetadata.Data {
		if m.Key == "loanId" {
			loanID = m.Value
		}
	}

	if loanID != "640747b60796a70001820d43" {
		t.Errorf("metadata loanId = %q, want %q", loanID, "640747b60796a70001820d43")
	}
}
