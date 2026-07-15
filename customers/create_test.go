package customers

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

// validTokenInput returns a TokenInput that passes the local validation
// guards in createToken so the request actually hits the httptest server.
func validTokenInput() TokenInput {
	return TokenInput{
		CardSource: payarc.CardSourceInternet,
		CardNumber: "4111111111111111",
		ExpMonth:   "04",
		ExpYear:    "2029",
		CCV:        "123",
	}
}

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

func TestCreateToken_BusinessDeclineLogsAtWarnAndReturnsSentinel(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    error
		wantLogLvl slog.Level
		wantLogMsg string
	}{
		{
			name:       "Do Not Honor → WARN + ErrDoNotHonor",
			statusCode: http.StatusConflict,
			body:       `{"message":"Do not honor"}`,
			wantErr:    payarc.ErrDoNotHonor,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "CVV2 Verification Failed → WARN + ErrCVV2Failed",
			statusCode: http.StatusConflict,
			body:       `{"message":"CVV2 verification failed"}`,
			wantErr:    payarc.ErrCVV2Failed,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "Invalid Card → WARN + ErrInvalidCard",
			statusCode: http.StatusConflict,
			body:       `{"message":"Invalid card"}`,
			wantErr:    payarc.ErrInvalidCard,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "Invalid CVV → WARN + ErrInvalidCCV",
			statusCode: http.StatusConflict,
			body:       `{"message":"Invalid CVV"}`,
			wantErr:    payarc.ErrInvalidCCV,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "Suspected Fraud → WARN + ErrSuspectedFraud",
			statusCode: http.StatusConflict,
			body:       `{"message":"Suspected fraud"}`,
			wantErr:    payarc.ErrSuspectedFraud,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "Suspected Card → WARN + ErrSuspectedCard",
			statusCode: http.StatusConflict,
			body:       `{"message":"Suspected card"}`,
			wantErr:    payarc.ErrSuspectedCard,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "Expired Card → WARN + ErrExpiredCard",
			statusCode: http.StatusConflict,
			body:       `{"message":"Expired Card"}`,
			wantErr:    payarc.ErrExpiredCard,
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc declined card token",
		},
		{
			name:       "card_holder_name format validation → WARN + per-field message",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"message":"The given data was invalid.","errors":{"card_holder_name":["The card holder name format is invalid."]}}`,
			wantErr:    errors.New("The card holder name format is invalid."),
			wantLogLvl: slog.LevelWarn,
			wantLogMsg: "payarc rejected card_holder_name format",
		},
		{
			name:       "Validation error on a non-card_holder_name field → ERROR (alerting path stays intact)",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"message":"The given data was invalid.","errors":{"card_number":["The card number is required."]}}`,
			wantErr:    errors.New("The given data was invalid."),
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc token create failed",
		},
		{
			name:       "Validation error on card_holder_name plus another field → ERROR (alerting path stays intact)",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"message":"The given data was invalid.","errors":{"card_holder_name":["The card holder name format is invalid."],"card_number":["The card number is required."]}}`,
			wantErr:    errors.New("The card holder name format is invalid."),
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc token create failed",
		},
		{
			name:       "Unknown message → ERROR (alerting path stays intact)",
			statusCode: http.StatusInternalServerError,
			body:       `{"message":"Some New Failure Mode"}`,
			wantErr:    nil, // wrapped errors.New, not a sentinel
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc token create failed",
		},
		{
			name:       "Unparseable body → ERROR (the parse-failure log line)",
			statusCode: http.StatusBadGateway,
			body:       `<html>upstream down</html>`,
			wantErr:    nil,
			wantLogLvl: slog.LevelError,
			wantLogMsg: "payarc token create failed; response body did not parse",
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

			_, err := svc.createToken(validTokenInput())
			if err == nil {
				t.Fatalf("expected error")
			}

			if tt.wantErr != nil {
				// For the validation-error case the returned error is a fresh
				// errors.New of the per-field message, not a sentinel, so fall
				// back to string compare instead of errors.Is when wantErr
				// isn't a sentinel.
				if errors.Is(tt.wantErr, payarc.ErrInvalidCard) ||
					errors.Is(tt.wantErr, payarc.ErrInvalidCCV) ||
					errors.Is(tt.wantErr, payarc.ErrCVV2Failed) ||
					errors.Is(tt.wantErr, payarc.ErrSuspectedFraud) ||
					errors.Is(tt.wantErr, payarc.ErrDoNotHonor) ||
					errors.Is(tt.wantErr, payarc.ErrSuspectedCard) ||
					errors.Is(tt.wantErr, payarc.ErrExpiredCard) {
					if !errors.Is(err, tt.wantErr) {
						t.Errorf("err = %v, want sentinel %v", err, tt.wantErr)
					}
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("err = %q, want %q", err.Error(), tt.wantErr.Error())
				}
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
