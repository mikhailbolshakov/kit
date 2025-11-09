package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MapInterfacesToBytesAndBack(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
	}{
		{
			name: "nil maps",
			m:    nil,
		}, {
			name: "empty maps",
			m:    map[string]interface{}{},
		}, {
			name: "map one value",
			m:    map[string]interface{}{"key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := MapAnyToBytes(tt.m)
			m := BytesToMapAny(bytes)
			assert.Equal(t, tt.m, m)
		})
	}
}

func Test_MapInterfacesToBytesNestedTypesAndBack(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
	}{
		{
			name: "diff type values",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
			},
		}, {
			name: "diff type values with map value",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
				"key4": map[string]interface{}{"key4internal1": float64(10), "key4internal2": "value2"}},
		}, {
			name: "diff type values with map value",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
				"key4": map[string]interface{}{
					"key4internal1": float64(10),
					"key4internal2": map[string]interface{}{
						"key4internal1": float64(10),
						"key4internal2": "value2"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := MapAnyToBytes(tt.m)
			m := BytesToMapAny(bytes)
			assertMap(t, tt.m, m)
		})
	}
}
