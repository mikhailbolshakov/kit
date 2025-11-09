package minio

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeErrMinioPutObject        = "S3-001"
	ErrCodeMinioCannotGetObject     = "S3-002"
	ErrCodeMinioCannotGetStatObject = "S3-003"
	ErrCodeMinioObjectNotFound      = "S3-004"
	ErrCodeMinioCreateBucket        = "S3-005"
	ErrCodeMinioRemoveObject        = "S3-006"
	ErrCodeMinioNew                 = "S3-007"
	ErrCodeMinioCopyObject          = "S3-008"
)

var (
	ErrMinioPutObject = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeErrMinioPutObject, "").Wrap(cause).C(ctx).Err()
	}
	ErrMinioCannotGetObject = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeMinioCannotGetObject, "").Wrap(cause).C(ctx).Err()
	}
	ErrMinioCannotGetStatObject = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeMinioCannotGetStatObject, "").Wrap(cause).C(ctx).Err()
	}
	ErrMinioObjectNotFound = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeMinioObjectNotFound, "").C(ctx).Err()
	}
	ErrMinioCreateBucket = func(cause error, ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeMinioCreateBucket, "").Wrap(cause).C(ctx).Err()
	}
	ErrMinioRemoveObject = func(cause error, ctx context.Context, fileId string) error {
		return kit.NewAppErrBuilder(ErrCodeMinioRemoveObject, "").Wrap(cause).C(ctx).F(kit.KV{"fileID ": fileId}).Err()
	}
	ErrMinioCopyObject = func(cause error, ctx context.Context, fileId string) error {
		return kit.NewAppErrBuilder(ErrCodeMinioCopyObject, "").Wrap(cause).C(ctx).F(kit.KV{"fileID ": fileId}).Err()
	}
	ErrMinioNew = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeMinioNew, "").Wrap(cause).Err()
	}
)
