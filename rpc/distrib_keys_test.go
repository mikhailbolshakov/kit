package rpc

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type distributedKeysTestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
	svc    DistributedKeys
}

func (s *distributedKeysTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
	s.svc = NewDistributedKeys()
}

func TestDistributedKeysSuite(t *testing.T) {
	suite.Run(t, new(distributedKeysTestSuite))
}

func (s *distributedKeysTestSuite) Test_CheckWhenEmpty() {
	s.False(s.svc.Check(kit.NewRandString()))
}

func (s *distributedKeysTestSuite) Test_RemoveWhenEmpty() {
	s.svc.Remove(kit.NewRandString())
}

func (s *distributedKeysTestSuite) Test_SetRemoveCheck() {
	key := kit.NewRandString()
	s.svc.Set(key)
	s.True(s.svc.Check(key))
	s.svc.Remove(key)
	s.False(s.svc.Check(key))
}
