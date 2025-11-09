package kafka

import (
	"context"

	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeKafkaFetchMessage                                = "KF-001"
	ErrCodeKafkaNotInitialized                              = "KF-003"
	ErrCodeKafkaInvalidConfig                               = "KF-004"
	ErrCodeKafkaMessageContextInvalid                       = "KF-005"
	ErrCodeKafkaMessageMarshal                              = "KF-006"
	ErrCodeKafkaProducerTopicEmpty                          = "KF-008"
	ErrCodeKafkaSubTopicEmpty                               = "KF-009"
	ErrCodeKafkaSubNoHandlers                               = "KF-010"
	ErrCodeKafkaConnection                                  = "KF-011"
	ErrCodeKafkaCreateTopics                                = "KF-012"
	ErrCodeKafkaMessageWrite                                = "KF-013"
	ErrCodeKafkaDecodeMsgUnmarshal                          = "KF-014"
	ErrCodeKafkaMsgUnmarshalPayload                         = "KF-015"
	ErrCodeKafkaProduceMsg                                  = "KF-016"
	ErrCodeKafkaSaslNotSupportedType                        = "KF-017"
	ErrCodeKafkaSaslGetMechanism                            = "KF-018"
	ErrCodeKafkaManualCommitRetryCountExceeded              = "KF-019"
	ErrCodeKafkaHandleMessageManualCommitRetryCountExceeded = "KF-020"
	ErrCodeKafkaSubscriberConfigInvalid                     = "KF-021"
)

var (
	ErrKafkaFetchMessage = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaFetchMessage, "").Wrap(cause).Err()
	}
	ErrKafkaManualCommitRetryCountExceeded = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaManualCommitRetryCountExceeded, "manual commit: retry exceeded").Err()
	}
	ErrKafkaHandleMessageManualCommitRetryCountExceeded = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaHandleMessageManualCommitRetryCountExceeded, "handle message: retry exceeded").Err()
	}
	ErrKafkaNotInitialized = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaNotInitialized, "not initialized").C(ctx).Err()
	}
	ErrKafkaProducerTopicEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaProducerTopicEmpty, "topic empty").C(ctx).Err()
	}
	ErrKafkaInvalidConfig = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaInvalidConfig, "config invalid").C(ctx).Err()
	}
	ErrKafkaMessageContextInvalid = func(ctx context.Context, topic string) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaMessageContextInvalid, "message context invalid").F(kit.KV{"topic": topic}).C(ctx).Err()
	}
	ErrKafkaMessageMarshal = func(ctx context.Context, cause error, topic string) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaMessageMarshal, "").Wrap(cause).F(kit.KV{"topic": topic}).C(ctx).Err()
	}
	ErrKafkaMessageWrite = func(ctx context.Context, cause error, topic string) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaMessageWrite, "").Wrap(cause).F(kit.KV{"topic": topic}).C(ctx).Err()
	}
	ErrKafkaConnection = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaConnection, "").Wrap(cause).C(ctx).Err()
	}
	ErrKafkaCreateTopics = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaCreateTopics, "").Wrap(cause).C(ctx).Err()
	}
	ErrKafkaDecodeMsgUnmarshal = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaDecodeMsgUnmarshal, "").Wrap(cause).C(ctx).Err()
	}
	ErrKafkaMsgUnmarshalPayload = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaMsgUnmarshalPayload, "").Wrap(cause).C(ctx).Err()
	}
	ErrKafkaSubTopicEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaSubTopicEmpty, "topic empty").C(ctx).Err()
	}
	ErrKafkaSubNoHandlers = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaSubNoHandlers, "no handlers specified").C(ctx).Err()
	}
	ErrKafkaProduceMsg = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaProduceMsg, "").Wrap(cause).C(ctx).Err()
	}
	ErrKafkaSaslNotSupportedType = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaSaslNotSupportedType, "not supported sasl type").C(ctx).Err()
	}
	ErrKafkaSubscriberConfigInvalid = func(ctx context.Context, reason string) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaSubscriberConfigInvalid, "subscriber invalid config: %s", reason).C(ctx).Err()
	}
	ErrKafkaSaslGetMechanism = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeKafkaSaslGetMechanism, "sasl mechanism").C(ctx).Err()
	}
)
