package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-sanitize/sanitize"
)

var (
	validate *validator.Validate
	san      *sanitize.Sanitizer
)

func init() {
	validate = validator.New()
	san, _ = sanitize.New()
}
