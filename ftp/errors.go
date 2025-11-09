package ftp

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeFtpUrlParse      = "FTP-001"
	ErrCodeFtpServerConnect = "FTP-002"
	ErrCodeFtpServerLogin   = "FTP-003"
	ErrCodeFtpRead          = "FTP-003"
	ErrCodeFtpLoad          = "FTP-004"
)

var (
	ErrFtpUrlParse = func(ctx context.Context, cause error, url string) error {
		return kit.NewAppErrBuilder(ErrCodeFtpUrlParse, "ftp url: parsing failed (%s)", url).Wrap(cause).Err()
	}
	ErrFtpServerConnect = func(ctx context.Context, cause error, url string) error {
		return kit.NewAppErrBuilder(ErrCodeFtpServerConnect, "ftp url: server connect (%s)", url).Wrap(cause).Err()
	}
	ErrFtpServerLogin = func(ctx context.Context, cause error, url string) error {
		return kit.NewAppErrBuilder(ErrCodeFtpServerLogin, "ftp url: login (%s)", url).Wrap(cause).Err()
	}
	ErrFtpRead = func(ctx context.Context, cause error, url string) error {
		return kit.NewAppErrBuilder(ErrCodeFtpRead, "ftp url: read file (%s)", url).Wrap(cause).Err()
	}
	ErrFtpLoad = func(ctx context.Context, cause error, url string) error {
		return kit.NewAppErrBuilder(ErrCodeFtpLoad, "ftp url: load file (%s)", url).Wrap(cause).Err()
	}
)
