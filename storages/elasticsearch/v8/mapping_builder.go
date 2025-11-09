package v8

import (
	"bytes"
	"github.com/mikhailbolshakov/kit"
	"io"
	"strings"
)

type MappingBody interface {
	Reader() io.Reader
}

func StringMapping(m string) MappingBody {
	return &mapping{
		str: m,
	}
}

func ByteMapping(b []byte) MappingBody {
	return &mapping{
		bytes: b,
	}
}

func JsonMapping(a any) MappingBody {
	json, _ := kit.JsonEncode(a)
	return ByteMapping(json)
}

type mapping struct {
	str   string
	bytes []byte
}

func (m mapping) Reader() io.Reader {
	if len(m.bytes) != 0 {
		return bytes.NewReader(m.bytes)
	}
	return strings.NewReader(m.str)
}
