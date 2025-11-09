package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"gitlab.com/algmib/kit/kafka"
	"gitlab.com/algmib/kit/mocks"
	"gitlab.com/algmib/kit/rpc"
)

type rpcServerTestSuite struct {
	kit.Suite
	logger       kit.CLoggerFunc
	callProducer *mocks.KafkaProducer
}

func (s *rpcServerTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func (s *rpcServerTestSuite) SetupTest() {
	s.callProducer = &mocks.KafkaProducer{}
}

func TestRpcServerSuite(t *testing.T) {
	suite.Run(t, new(rpcServerTestSuite))
}

type Body struct {
	Value string `json:"val"`
}

func (s *rpcServerTestSuite) Test_Call_NoResponseRequired_Ok() {
	rpcServer := NewServer(s.logger, s.callProducer, rpc.NewDistributedKeys(), &rpc.Config{}).(*rpcServer)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: false,
		Body:             &Body{Value: kit.NewRandString()},
	}
	var actualMsg *rpc.Message
	rpcServer.RegisterType(msg.Type, func(ctx context.Context, msg *rpc.Message) error {
		actualMsg = msg
		return nil
	}, func() interface{} { return &Body{} })
	rqCtx, _ := kit.Request(s.Ctx)
	kafkaMsg := &kafka.Message{
		Ctx:     rqCtx,
		Key:     msg.Key,
		Payload: msg,
	}
	kafkaMsgBytes, _ := kit.Marshal(kafkaMsg)
	s.Nil(rpcServer.RequestHandler(kafkaMsgBytes))
	s.Equal(0, rpcServer.rqPool.Len())

	if err := <-kit.Await(func() (bool, error) {
		return actualMsg != nil, nil
	}, time.Millisecond*500, time.Second*3); err != nil {
		s.Fatal(err)
	}

	s.Equal(msg.Key, actualMsg.Key)
	s.Equal(msg.RequestId, actualMsg.RequestId)
	s.Equal(msg.Body.(*Body).Value, actualMsg.Body.(*Body).Value)
}

func (s *rpcServerTestSuite) Test_Call_ResponseRequired_Ok() {
	rpcServer := NewServer(s.logger, s.callProducer, rpc.NewDistributedKeys(), &rpc.Config{}).(*rpcServer)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: true,
		Body:             &Body{Value: kit.NewRandString()},
	}
	s.callProducer.On("Send", s.Ctx, msg.Key, msg).Return(nil)
	var actualMsg *rpc.Message
	rpcServer.RegisterType(msg.Type, func(ctx context.Context, msg *rpc.Message) error {
		actualMsg = msg
		s.Nil(rpcServer.Response(s.Ctx, msg))
		return nil
	}, func() interface{} { return &Body{} })
	rqCtx, _ := kit.Request(s.Ctx)
	kafkaMsg := &kafka.Message{
		Ctx:     rqCtx,
		Key:     msg.Key,
		Payload: msg,
	}
	kafkaMsgBytes, _ := kit.Marshal(kafkaMsg)
	s.Nil(rpcServer.RequestHandler(kafkaMsgBytes))

	if err := <-kit.Await(func() (bool, error) {
		return actualMsg != nil && rpcServer.rqPool.Len() == 0, nil
	}, time.Millisecond*500, time.Second*3); err != nil {
		s.Fatal(err)
	}

	s.AssertCalled(&s.callProducer.Mock, "Send", s.Ctx, msg.Key, msg)
}

func (s *rpcServerTestSuite) Test_Call_WhenRequestTimeout_Ok() {
	rpcServer := NewServer(s.logger, s.callProducer, rpc.NewDistributedKeys(),
		&rpc.Config{CallTimeOut: time.Second}).(*rpcServer)
	msg := &rpc.Message{
		Type:             rpc.MessageType(1),
		Key:              kit.NewRandString(),
		RequestId:        kit.NewRandString(),
		ResponseRequired: true,
		Body:             &Body{Value: kit.NewRandString()},
	}
	var actualExpiredMsg *rpc.Message
	rpcServer.SetExpirationCallback(func(ctx context.Context, msg *rpc.Message) error {
		actualExpiredMsg = msg
		return nil
	})
	rpcServer.RegisterType(msg.Type, func(ctx context.Context, msg *rpc.Message) error {
		return nil
	}, func() interface{} { return &Body{} })
	rpcServer.Start(s.Ctx)
	defer rpcServer.Close(s.Ctx)
	rqCtx, _ := kit.Request(s.Ctx)
	kafkaMsg := &kafka.Message{
		Ctx:     rqCtx,
		Key:     msg.Key,
		Payload: msg,
	}
	kafkaMsgBytes, _ := kit.Marshal(kafkaMsg)
	s.Nil(rpcServer.RequestHandler(kafkaMsgBytes))
	if err := <-kit.Await(func() (bool, error) {
		return actualExpiredMsg != nil, nil
	}, time.Millisecond*500, time.Second*3); err != nil {
		s.Fatal(err)
	}
	s.Equal(0, rpcServer.rqPool.Len())
	s.Equal(msg.RequestId, actualExpiredMsg.RequestId)
}
