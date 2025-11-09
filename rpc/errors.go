package rpc

import (
	"context"
	"gitlab.com/algmib/kit"
)

const (
	ErrCodeRpcMsgNoKey            = "RPC-001"
	ErrCodeRpcMsgNoRequestId      = "RPC-002"
	ErrCodeRpcRespNoRequestInPool = "RPC-003"
	ErrCodeRpcRespInvalidBody     = "RPC-004"
	ErrCodeRpcCallNoCb            = "RPC-005"
)

var (
	ErrRpcMsgNoKey = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeRpcMsgNoKey, "key empty").C(ctx).Err()
	}
	ErrRpcMsgNoRequestId = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeRpcMsgNoRequestId, "Request id empty").C(ctx).Err()
	}
	ErrRpcCallNoCb = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeRpcCallNoCb, "callback id empty").C(ctx).Err()
	}
	ErrRpcRespNoRequestInPool = func(ctx context.Context, rqId, key string) error {
		return kit.NewAppErrBuilder(ErrCodeRpcRespNoRequestInPool, "no Request in pool").C(ctx).F(kit.KV{"rqId": rqId, "key": key}).Err()
	}
	ErrRpcRespInvalidBody = func(cause error, ctx context.Context, rqId, key string) error {
		return kit.NewAppErrBuilder(ErrCodeRpcRespInvalidBody, "no Request in pool").Wrap(cause).C(ctx).F(kit.KV{"rqId": rqId, "key": key}).Err()
	}
)
