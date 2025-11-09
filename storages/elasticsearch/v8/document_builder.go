package v8

import (
	"bytes"
	"gitlab.com/algmib/kit"
	"io"
	"strings"
)

type DataBody interface {
	Reader() io.Reader
}

func StringData(m string) DataBody {
	return &data{
		str: m,
	}
}

func ByteData(b []byte) DataBody {
	return &data{
		bytes: b,
	}
}

func JsonData(a any) DataBody {
	json, _ := kit.JsonEncode(a)
	return ByteData(json)
}

type data struct {
	str   string
	bytes []byte
}

func (m data) Reader() io.Reader {
	if len(m.bytes) != 0 {
		return bytes.NewReader(m.bytes)
	}
	return strings.NewReader(m.str)
}
