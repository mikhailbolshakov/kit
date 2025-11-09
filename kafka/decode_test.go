package kafka

import (
	"testing"

	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
)

type decodeTestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *decodeTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestDecodeSuite(t *testing.T) {
	suite.Run(t, new(decodeTestSuite))
}

func (s *decodeTestSuite) Test_WhenRawPayload() {
	rId := kit.NewId()
	rqCtx := kit.NewRequestCtx().WithRequestId(rId)
	rawPl := map[string]any{
		"key1": "val1",
		"key2": "val2",
	}
	msg := &MessageT[map[string]any]{
		Ctx:     rqCtx,
		Key:     kit.NewRandString(),
		Payload: rawPl,
	}
	msgBytes, _ := kit.Marshal(msg)
	decoded, ctx, err := Decode[map[string]any](s.Ctx, msgBytes)
	s.Nil(err)
	s.NotEmpty(ctx)
	actRq, _ := kit.Request(ctx)
	s.NotEmpty(actRq)
	s.Equal(actRq.Rid, rId)
	s.Equal(decoded, msg.Payload)
}

func (s *decodeTestSuite) Test_WhenStructUnmarshalledToRawPayload() {
	rId := kit.NewId()
	rqCtx := kit.NewRequestCtx().WithRequestId(rId)
	pl := struct {
		Key string `json:"key"`
	}{
		Key: "123",
	}
	msg := &Message{
		Ctx:     rqCtx,
		Key:     kit.NewRandString(),
		Payload: pl,
	}
	msgBytes, _ := kit.Marshal(msg)
	decoded, ctx, err := Decode[map[string]any](s.Ctx, msgBytes)
	s.Nil(err)
	s.NotEmpty(ctx)
	actRq, _ := kit.Request(ctx)
	s.NotEmpty(actRq)
	s.Equal(actRq.Rid, rId)
	s.Equal(decoded["key"].(string), pl.Key)
}

func (s *decodeTestSuite) Test_WhenStructUnmarshalledToStruct() {
	rqCtx, _ := kit.Request(s.Ctx)
	type plt struct {
		Key string `json:"key"`
	}
	pl := &plt{
		Key: "123",
	}
	msg := &Message{
		Ctx:     rqCtx,
		Key:     kit.NewRandString(),
		Payload: pl,
	}
	msgBytes, _ := kit.Marshal(msg)
	decoded, ctx, err := Decode[plt](s.Ctx, msgBytes)
	s.Nil(err)
	s.NotEmpty(ctx)
	s.Equal(decoded.Key, pl.Key)
}
