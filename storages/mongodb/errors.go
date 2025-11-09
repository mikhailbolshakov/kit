package mongodb

import "github.com/mikhailbolshakov/kit"

const (
	ErrCodeReadCertFile = "MNG-001"
	ErrCodeAppendCert   = "MNG-002"
	ErrCodeConnection   = "MNG-003"
)

var (
	ErrReadCertFile = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeReadCertFile, "certificate file reading error").Wrap(cause).Err()
	}
	ErrAppendCert = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeAppendCert, "appending certificate error").Wrap(cause).Err()
	}
	ErrConnection = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeConnection, "connection error").Wrap(cause).Err()
	}
)
