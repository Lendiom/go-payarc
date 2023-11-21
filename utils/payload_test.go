package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGenerateFormPayload(t *testing.T) {
	type testItem struct {
		Amount     int    `form:"amount,omitempty"`
		CustomerID string `form:"customer_id,omitempty"`
	}

	type args struct {
		val interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    func() url.Values
		wantErr bool
	}{
		{
			name: "Valid Full Object",
			args: args{
				val: testItem{
					Amount:     1000,
					CustomerID: "test-id",
				},
			},
			want: func() url.Values {
				v := url.Values{}
				v.Add("amount", "1000")
				v.Add("customer_id", "test-id")

				return v
			},
			wantErr: false,
		},
		{
			name: "Valid Empty Object",
			args: args{
				val: testItem{},
			},
			want: func() url.Values {
				return url.Values{}
			},
			wantErr: false,
		},
		{
			name: "Valid Partial Object",
			args: args{
				val: testItem{
					Amount: 1000,
				},
			},
			want: func() url.Values {
				v := url.Values{}
				v.Add("amount", "1000")

				return v
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFormPayload(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFormPayload() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want()) {
				t.Errorf("GenerateFormPayload() = %+v, want %+v", got, tt.want())
			}
		})
	}
}
