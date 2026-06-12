package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFieldCreation(t *testing.T) {
	lf := LogField{Key: "user_id", Value: "12345"}
	assert.Equal(t, "user_id", lf.Key)
	assert.Equal(t, "12345", lf.Value)
}

func TestLogFieldVariousTypes(t *testing.T) {
	tests := []struct {
		name  string
		field LogField
	}{
		{"string value", LogField{Key: "key1", Value: "string"}},
		{"int value", LogField{Key: "key2", Value: 42}},
		{"bool value", LogField{Key: "key3", Value: true}},
		{"float value", LogField{Key: "key4", Value: 3.14}},
		{"slice value", LogField{Key: "key5", Value: []string{"a", "b"}}},
		{"nil value", LogField{Key: "key6", Value: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.field.Key, tt.field.Key) // just verify no panic
			switch tt.name {
			case "string value":
				assert.Equal(t, "string", tt.field.Value)
			case "int value":
				assert.Equal(t, 42, tt.field.Value)
			case "bool value":
				assert.Equal(t, true, tt.field.Value)
			case "float value":
				assert.Equal(t, 3.14, tt.field.Value)
			case "nil value":
				assert.Nil(t, tt.field.Value)
			}
		})
	}
}

func TestLogFieldZeroValue(t *testing.T) {
	lf := LogField{}
	assert.Empty(t, lf.Key)
	assert.Nil(t, lf.Value)
}
