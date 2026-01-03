package logger

import "go.uber.org/zap"

// convertFields converts a map[string]interface{} into a slice of zap.Field,
// producing one zap.Field for each map entry. If the input is nil, convertFields returns nil.
func convertFields(fields map[string]interface{}) []zap.Field {
	if fields == nil {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}