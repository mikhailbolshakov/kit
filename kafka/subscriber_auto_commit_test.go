package kafka

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
)

type subscriberAutoCommitTestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *subscriberAutoCommitTestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestSubscriberAutoCommitSuite(t *testing.T) {
	suite.Run(t, new(subscriberAutoCommitTestSuite))
}

func (s *subscriberAutoCommitTestSuite) Test_IndexByKey() {

	test := func(workers int, keys []string, exp ...int) {
		sub := &subscriberAutoCommit{workers: workers}
		for i, k := range keys {
			s.Equal(exp[i], sub.chanIndexByKey([]byte(k)))
		}
	}

	test(1, []string{"1", "2", "33244", kit.NewRandString(), "AAaaFFFff"}, 0, 0, 0, 0, 0)
	test(2, []string{"1", "1", "2", "2"}, 0, 0, 1, 1)
	test(2, []string{"aaFFaaFF", "bbCCbbCD", "aaFFaaFF", "bbCCbbCD"}, 1, 0, 1, 0)

	randKey := kit.NewRandString()
	sub := &subscriberAutoCommit{workers: 10}
	s.Equal(sub.chanIndexByKey([]byte(randKey)), sub.chanIndexByKey([]byte(randKey)))

}
