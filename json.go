package kit

import (
	"github.com/goccy/go-json"
	"io"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

func BytesToMapAny(bytes []byte) map[string]any {
	mp := make(map[string]interface{})
	_ = Unmarshal(bytes, &mp)
	return mp
}

func MapAnyToBytes(m map[string]any) []byte {
	bytes, _ := Marshal(m)
	return bytes
}

// JsonEncode encodes type to json bytes
func JsonEncode(v any) ([]byte, error) {
	r, err := Marshal(&v)
	if err != nil {
		return nil, ErrJsonEncode(err)
	}
	return r, nil
}

// JsonDecode decodes type from json bytes
func JsonDecode[T any](payload []byte) (*T, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	var res T
	err := Unmarshal(payload, &res)
	if err != nil {
		return nil, ErrJsonDecode(err)
	}
	return &res, nil
}

// JsonDecodeSlice decodes type from json bytes to slice
func JsonDecodeSlice[T any](payload []byte) ([]*T, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	res, err := JsonDecodePlainSlice[T](payload)
	if err != nil {
		return nil, err
	}
	return ToSlicePtr[T](res), nil
}

// JsonDecodePlainSlice decodes type from json bytes to slice
func JsonDecodePlainSlice[T any](payload []byte) ([]T, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	var res []T
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return nil, ErrJsonDecode(err)
	}
	return res, nil
}
