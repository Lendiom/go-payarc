package charges

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Lendiom/go-payarc"
	"github.com/Lendiom/go-payarc/client"
)

// levelRecorder is a slog.Handler that captures the level and message of
// every record routed through it. It exists so tests can assert the
// classification chose WARN vs ERROR without inspecting log output buffers.
type levelRecorder struct {
	records []slog.Record
}

func (lr *levelRecorder) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (lr *levelRecorder) Handle(_ context.Context, r slog.Record) error {
	lr.records = append(lr.records, r)
	return nil
}
func (lr *levelRecorder) WithAttrs([]slog.Attr) slog.Handler { return lr }
func (lr *levelRecorder) WithGroup(string) slog.Handler      { return lr }

func newServiceForTest(t *testing.T, statusCode int, body string) (*Service, func()) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))

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

	return svc, srv.Close
}

func TestCreate_BusinessDeclineLogsAtWarnAndReturnsSentinel(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        string
		wantErr     error
		wantLogLvl  slog.Level
		wantLogMsg  string
	}{
		{
			name:       "Insufficient Funds → WARN + ErrInsufficientFunds",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"status":"error","code":0,"message":"Insufficient Funds","status_code":422,"data":{}}`,
			wantErr:    payarc.ErrInsufficientFunds,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined charge",
		},
		{
			name:       "Expired Card → WARN + ErrExpiredCard",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"status":"error","message":"Expired Card","status_code":422}`,
			wantErr:    payarc.ErrExpiredCard,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined charge",
		},
		{
			name:       "Do Not Honor → WARN + ErrDoNotHonor",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"status":"error","message":"Do Not Honor","status_code":422}`,
			wantErr:    payarc.ErrDoNotHonor,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined charge",
		},
		{
			name:       "Re-enter transaction → WARN + ErrReEnterTransaction",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"status":"error","message":"Re-enter transaction","status_code":422,"data":{"failure_code":"D0092","host_response_code":"19"}}`,
			wantErr:    payarc.ErrReEnterTransaction,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined charge",
		},
		{
			name:       "Unknown message → ERROR (alerting path stays intact)",
			statusCode: http.StatusInternalServerError,
			body:       `{"status":"error","message":"Some New Failure Mode","status_code":500}`,
			wantErr:    nil, // wrapped fmt.Errorf, not a sentinel
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc charge create failed",
		},
		{
			name:       "Unparseable body → ERROR (the parse-failure log line)",
			statusCode: http.StatusBadGateway,
			body:       `<html>upstream down</html>`,
			wantErr:    nil,
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc charge create failed; response body did not parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := newServiceForTest(t, tt.statusCode, tt.body)
			defer cleanup()

			recorder := &levelRecorder{}
			origLogger := slog.Default()
			slog.SetDefault(slog.New(recorder))
			defer slog.SetDefault(origLogger)

			_, err := svc.Create(ChargeInput{Amount: 100, CustomerID: "cust_1"})
			if err == nil {
				t.Fatalf("expected error")
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want sentinel %v", err, tt.wantErr)
			}

			if len(recorder.records) != 1 {
				t.Fatalf("expected exactly 1 log record, got %d", len(recorder.records))
			}

			rec := recorder.records[0]
			if rec.Level != tt.wantLogLvl {
				t.Errorf("log level = %v, want %v (msg=%q)", rec.Level, tt.wantLogLvl, rec.Message)
			}

			if rec.Message != tt.wantLogMsg {
				t.Errorf("log message = %q, want %q", rec.Message, tt.wantLogMsg)
			}
		})
	}
}
