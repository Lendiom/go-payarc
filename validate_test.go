package payarc

import (
	"errors"
	"strings"
	"testing"
)

func TestRequireParam(t *testing.T) {
	type args struct {
		name  string
		value string
	}

	tests := []struct {
		name string
		args args

		wantErr     bool
		wantErrText string
	}{
		{
			name: "an ordinary identifier passes",
			args: args{name: "charge id", value: "ch_abc123"},
		},
		{
			name: "an identifier with surrounding content passes",
			args: args{name: "charge id", value: "eg09AGgggEDEGadD"},
		},
		{
			// Void("") is what produced v1/charges//void in production.
			name:        "an empty value is rejected and names the argument",
			args:        args{name: "charge id", value: ""},
			wantErr:     true,
			wantErrText: "required parameter is missing: charge id",
		},
		{
			name:        "a whitespace-only value is rejected",
			args:        args{name: "customer id", value: "   "},
			wantErr:     true,
			wantErrText: "required parameter is missing: customer id",
		},
		{
			name:        "a tab and newline value is rejected",
			args:        args{name: "card id", value: "\t\n"},
			wantErr:     true,
			wantErrText: "required parameter is missing: card id",
		},
		{
			// Would silently retarget the request at a different endpoint.
			name:        "a value containing a slash is rejected",
			args:        args{name: "customer id", value: "cus_1/cards"},
			wantErr:     true,
			wantErrText: "required parameter is missing: customer id contains a slash, which would change which endpoint is called",
		},
		{
			name:        "a leading slash is rejected",
			args:        args{name: "charge id", value: "/v1/charges"},
			wantErr:     true,
			wantErrText: "required parameter is missing: charge id contains a slash, which would change which endpoint is called",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequireParam(tt.args.name, tt.args.value)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("RequireParam returned error: %v", err)
				}

				return
			}

			if err == nil {
				t.Fatal("RequireParam returned nil, want an error")
			}

			if !errors.Is(err, ErrMissingParameter) {
				t.Errorf("error %v does not match ErrMissingParameter", err)
			}

			if err.Error() != tt.wantErrText {
				t.Errorf("error text = %q, want %q", err.Error(), tt.wantErrText)
			}

			// The argument name has to survive into the message, or a caller staring at a log
			// line cannot tell which of two ids was the blank one.
			if !strings.Contains(err.Error(), tt.args.name) {
				t.Errorf("error %q does not name the argument %q", err, tt.args.name)
			}
		})
	}
}
