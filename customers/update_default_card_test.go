package customers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Lendiom/go-payarc/client"
)

// recordingService returns a Service whose requests are captured rather than
// sent to PayArc, along with pointers to the last method, path and body seen.
func recordingService(t *testing.T) (*Service, *string, *string, *string, func()) {
	t.Helper()

	var method, path, body string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		body = string(raw)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse httptest URL: %v", err)
	}

	svc := &Service{client: client.Client{ApiKey: "test-key", HttpClient: *srv.Client(), Url: *u}}

	return svc, &method, &path, &body, srv.Close
}

// TestClearDefaultCard sends an empty default_card_id, which is how a customer
// is told they have no default card. A payer whose chosen method is a bank
// account has no card to point at, and default_card_id accepts nothing else, so
// clearing is the only way to stop the field naming a stale card.
func TestClearDefaultCard(t *testing.T) {
	svc, method, path, body, cleanup := recordingService(t)
	defer cleanup()

	if err := svc.ClearDefaultCard("cus_123"); err != nil {
		t.Fatalf("ClearDefaultCard: %v", err)
	}

	if *method != http.MethodPatch {
		t.Errorf("method = %s, want PATCH", *method)
	}

	if *path != "/cus_123" {
		t.Errorf("path = %s, want /cus_123", *path)
	}

	if *body != "default_card_id=" {
		t.Errorf("body = %q, want %q", *body, "default_card_id=")
	}
}

// TestClearDefaultCardRequiresCustomerID keeps the one guard that still applies:
// an empty card id is the point of this call, but an empty customer id would
// PATCH the collection endpoint instead of a customer.
func TestClearDefaultCardRequiresCustomerID(t *testing.T) {
	svc, _, _, _, cleanup := recordingService(t)
	defer cleanup()

	if err := svc.ClearDefaultCard("  "); err == nil {
		t.Fatal("expected an error for a blank customer id")
	}
}

// TestUpdateDefaultCardStillRequiresACardID guards the split: sharing the
// transport with ClearDefaultCard must not make the setter accept an empty id,
// which would silently clear the default instead of setting it.
func TestUpdateDefaultCardStillRequiresACardID(t *testing.T) {
	svc, _, _, _, cleanup := recordingService(t)
	defer cleanup()

	if err := svc.UpdateDefaultCard("cus_123", ""); err == nil {
		t.Fatal("expected an error for a blank default card id")
	}
}

// TestUpdateDefaultCardEncodesTheCardID pins that the id is form-encoded rather
// than interpolated, so a value containing a reserved character cannot alter the
// payload.
func TestUpdateDefaultCardEncodesTheCardID(t *testing.T) {
	svc, _, _, body, cleanup := recordingService(t)
	defer cleanup()

	if err := svc.UpdateDefaultCard("cus_123", "card&x=1"); err != nil {
		t.Fatalf("UpdateDefaultCard: %v", err)
	}

	if *body != "default_card_id=card%26x%3D1" {
		t.Errorf("body = %q, want the id percent-encoded", *body)
	}
}
