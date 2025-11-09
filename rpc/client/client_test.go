package client

import (
	"context"
	"testing"
	"time"

	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/kafka"
	"github.com/mikhailbolshakov/kit/mocks"
	"github.com/mikhailbolshakov/kit/rpc"
	"github.com/stretchr/testify/suite"
)

type rpcClientTestSuite struct {
	kit.Suite
	logger       kit.CLoggerFunc
	callProducer *mocks.KafkaProducer
}

func (s *rpcClientTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func (s *rpcClientTestSuite) SetupTest() {
	s.callProducer = &mocks.KafkaProducer{}
}

func TestRpcClientSuite(t *testing.T) {
	suite.Run(t, new(rpcClientTestSuite))
}

type Body struct {
	Value string `json:"val"`
}

func (s *rpcClientTestSuite) Test_Call_NoResponseRequired_Ok() {
	rpcCl := NewClient(s.logger, s.callProducer, rpc.NewDistributedKeys(), &rpc.Config{}).(*rpcClient)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: false,
		Body:             &Body{Value: kit.NewRandString()},
	}
	s.callProducer.On("Send", s.Ctx, msg.Key, msg).Return(nil)
	err := rpcCl.Call(s.Ctx, msg, nil)
	s.NoError(err)
	s.AssertCalled(&s.callProducer.Mock, "Send", s.Ctx, msg.Key, msg)
	s.Equal(0, rpcCl.rqPool.Len())
}

func (s *rpcClientTestSuite) Test_Call_WithResponse_Ok() {
	rpcCl := NewClient(s.logger, s.callProducer, rpc.NewDistributedKeys(), &rpc.Config{}).(*rpcClient)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: true,
		Body:             &Body{Value: kit.NewRandString()},
	}
	s.callProducer.On("Send", s.Ctx, msg.Key, msg).Return(nil)
	var actualRsMsg *rpc.Message
	rpcCl.RegisterBodyTypeProvider(msg.Type, func() interface{} { return &Body{} })
	err := rpcCl.Call(s.Ctx, msg, func(ctx context.Context, rqMsg, rsMsg *rpc.Message) error {
		actualRsMsg = rsMsg
		return nil
	})
	s.NoError(err)
	s.AssertCalled(&s.callProducer.Mock, "Send", s.Ctx, msg.Key, msg)
	s.Equal(1, rpcCl.rqPool.Len())
	// call handler
	rsMsg := &rpc.Message{
		Key:       msg.Key,
		RequestId: msg.RequestId,
		Type:      rpc.MessageType(1),
		Body:      &Body{Value: kit.NewRandString()},
	}
	kafkaMsg := &kafka.Message{
		Ctx:     nil,
		Key:     rsMsg.Key,
		Payload: rsMsg,
	}
	kafkaMsgBytes, _ := kit.Marshal(kafkaMsg)
	s.Nil(rpcCl.ResponseHandler(kafkaMsgBytes))

	if err := <-kit.Await(func() (bool, error) {
		return actualRsMsg != nil && rpcCl.rqPool.Len() == 0, nil
	}, time.Millisecond*500, time.Second*3); err != nil {
		s.Fatal(err)
	}

	s.Equal(rsMsg.Key, actualRsMsg.Key)
	s.Equal(rsMsg.RequestId, actualRsMsg.RequestId)
	s.Equal(rsMsg.Body.(*Body).Value, actualRsMsg.Body.(*Body).Value)
}

func (s *rpcClientTestSuite) Test_Call_WhenRequestExpired_Fail() {
	rpcCl := NewClient(s.logger, s.callProducer, rpc.NewDistributedKeys(), &rpc.Config{}).(*rpcClient)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: true,
		Body:             &Body{Value: kit.NewRandString()},
	}
	s.callProducer.On("Send", s.Ctx, msg.Key, msg).Return(nil)
	rpcCl.RegisterBodyTypeProvider(msg.Type, func() interface{} { return &Body{} })
	err := rpcCl.Call(s.Ctx, msg, func(ctx context.Context, rqMsg, rsMsg *rpc.Message) error {
		return nil
	})
	s.NoError(err)
	s.AssertCalled(&s.callProducer.Mock, "Send", s.Ctx, msg.Key, msg)
	s.Equal(1, rpcCl.rqPool.Len())
	// remove request from pool
	rpcCl.rqPool.Remove(msg.RequestId)
	// call handler
	rsMsg := &rpc.Message{
		Key:       msg.Key,
		RequestId: msg.RequestId,
		Type:      rpc.MessageType(1),
		Body:      &Body{Value: kit.NewRandString()},
	}
	rqCtx, _ := kit.Request(s.Ctx)
	kafkaMsg := &kafka.Message{
		Ctx:     rqCtx,
		Key:     rsMsg.Key,
		Payload: rsMsg,
	}
	kafkaMsgBytes, _ := kit.Marshal(kafkaMsg)
	err = rpcCl.ResponseHandler(kafkaMsgBytes)
	s.AssertAppErr(err, rpc.ErrCodeRpcRespNoRequestInPool)
}

func (s *rpcClientTestSuite) Test_Call_WhenRequestTimeout_Ok() {
	rpcCl := NewClient(s.logger, s.callProducer, rpc.NewDistributedKeys(),
		&rpc.Config{CallTimeOut: time.Second}).(*rpcClient)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: true,
		Body:             &Body{Value: kit.NewRandString()},
	}
	var actualExpiredMsg *rpc.Message
	rpcCl.SetExpirationCallback(func(ctx context.Context, msg *rpc.Message) error {
		actualExpiredMsg = msg
		return nil
	})
	s.callProducer.On("Send", s.Ctx, msg.Key, msg).Return(nil)
	rpcCl.RegisterBodyTypeProvider(msg.Type, func() interface{} { return &Body{} })
	rpcCl.Start(s.Ctx)
	defer rpcCl.Close(s.Ctx)
	err := rpcCl.Call(s.Ctx, msg, func(ctx context.Context, rqMsg, rsMsg *rpc.Message) error {
		return nil
	})
	s.NoError(err)
	if err := <-kit.Await(func() (bool, error) {
		return actualExpiredMsg != nil, nil
	}, time.Millisecond*500, time.Second*3); err != nil {
		s.Fatal(err)
	}
	s.Equal(0, rpcCl.rqPool.Len())
	s.Equal(msg.RequestId, actualExpiredMsg.RequestId)
}
