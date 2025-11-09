package v8

import (
	"io"
	"strings"
)

type QueryBody interface {
	Reader() io.Reader
}

func MatchAllQuery() QueryBody {
	return matchAllQuery{}
}

type matchAllQuery struct{}

func (q matchAllQuery) Reader() io.Reader {
	return strings.NewReader(`{"query": {"match_all": {}}}`)
}
