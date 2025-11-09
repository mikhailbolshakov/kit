package grpc

import (
	"context"
	"github.com/mikhailbolshakov/kit"
	"google.golang.org/grpc/codes"
)

const (
	ErrCodeGrpcClientDial     = "GRPC-001"
	ErrCodeGrpcInvoke         = "GRPC-002"
	ErrCodeGrpcSrvListen      = "GRPC-003"
	ErrCodeGrpcSrvServe       = "GRPC-004"
	ErrCodeGrpcSrvNotReady    = "GRPC-005"
	ErrCodeGrpcClientError    = "GRPC-006"
	ErrCodeGrpcAuthNoMd       = "GRPC-007"
	ErrCodeGrpcAuthNoHeader   = "GRPC-008"
	ErrCodeGrpcAuthParseToken = "GRPC-009"
)

var (
	ErrGrpcClientDial = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcClientDial, "").Wrap(cause).Err()
	}
	ErrGrpcInvoke = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcInvoke, "").Wrap(cause).Err()
	}
	ErrGrpcSrvListen = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcSrvListen, "").Wrap(cause).Err()
	}
	ErrGrpcSrvServe = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcSrvServe, "").Wrap(cause).Err()
	}
	ErrGrpcSrvNotReady = func(svc string) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcSrvNotReady, "service isn't ready within timeout").F(kit.KV{"svc": svc}).Err()
	}
	ErrGrpcClientError = func(ctx context.Context, cause error, method string) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcClientError, "grpc client error").F(kit.KV{"method": method}).Wrap(cause).Err()
	}
	ErrGrpcAuthNoMd = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcAuthNoMd, "no metadata").GrpcSt(uint32(codes.Unauthenticated)).Err()
	}
	ErrGrpcAuthNoHeader = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcAuthNoHeader, "no header").GrpcSt(uint32(codes.Unauthenticated)).Err()
	}
	ErrGrpcAuthParseToken = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeGrpcAuthParseToken, "no header").Wrap(cause).GrpcSt(uint32(codes.Unauthenticated)).Err()
	}
)
