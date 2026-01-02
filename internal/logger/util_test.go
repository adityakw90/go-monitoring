package logger

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestLogger_ConvertFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  map[string]interface{}
		want    []zap.Field
		wantErr bool
	}{
		{
			name:    "nil fields",
			fields:  nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertFields(tt.fields)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
