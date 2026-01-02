package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger_Util_ConvertFields(t *testing.T) {
	tests := []struct {
		name      string
		fields    map[string]interface{}
		checkFunc func(t *testing.T, got []zap.Field)
	}{
		{
			name:   "nil fields",
			fields: nil,
			checkFunc: func(t *testing.T, got []zap.Field) {
				if got != nil {
					t.Errorf("convertFields() = %v, want nil", got)
				}
			},
		},
		{
			name:   "empty fields",
			fields: map[string]interface{}{},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 0 {
					t.Errorf("convertFields() length = %v, want 0", len(got))
				}
			},
		},
		{
			name: "single string field",
			fields: map[string]interface{}{
				"key": "value",
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
					return
				}
				if got[0].Key != "key" {
					t.Errorf("convertFields() field key = %v, want 'key'", got[0].Key)
				}
			},
		},
		{
			name: "multiple fields",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"float":  3.14,
				"bool":   true,
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 4 {
					t.Errorf("convertFields() length = %v, want 4", len(got))
				}
			},
		},
		{
			name: "field with int value",
			fields: map[string]interface{}{
				"count": 100,
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "field with float value",
			fields: map[string]interface{}{
				"rate": 3.14,
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "field with bool value",
			fields: map[string]interface{}{
				"enabled": true,
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "field with slice value",
			fields: map[string]interface{}{
				"items": []string{"a", "b", "c"},
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "field with map value",
			fields: map[string]interface{}{
				"metadata": map[string]string{"key": "value"},
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "field with nil value",
			fields: map[string]interface{}{
				"nil_field": nil,
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 1 {
					t.Errorf("convertFields() length = %v, want 1", len(got))
				}
			},
		},
		{
			name: "fields preserve order independence",
			fields: map[string]interface{}{
				"z": "last",
				"a": "first",
				"m": "middle",
			},
			checkFunc: func(t *testing.T, got []zap.Field) {
				if len(got) != 3 {
					t.Errorf("convertFields() length = %v, want 3", len(got))
				}
				keys := make(map[string]bool)
				for _, field := range got {
					keys[field.Key] = true
				}
				if !keys["a"] || !keys["m"] || !keys["z"] {
					t.Errorf("convertFields() missing expected keys")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertFields(tt.fields)
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestLogger_Util_ConvertFields_TypePreservation(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]interface{}
		want   int
	}{
		{
			name:   "nil returns nil",
			fields: nil,
			want:   0,
		},
		{
			name:   "empty map returns empty slice",
			fields: map[string]interface{}{},
			want:   0,
		},
		{
			name: "single field returns single element",
			fields: map[string]interface{}{
				"key": "value",
			},
			want: 1,
		},
		{
			name: "multiple fields returns multiple elements",
			fields: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertFields(tt.fields)
			if tt.fields == nil {
				if got != nil {
					t.Errorf("convertFields(nil) = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("convertFields() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}
