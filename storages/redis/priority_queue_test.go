//go:build integration

package redis

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"testing"
)

type priorityQueueTestSuite struct {
	kit.Suite
	config *Config
}

func (s *priorityQueueTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
	s.config = &Config{
		Host: "localhost",
		Port: "6379",
		Ttl:  0,
	}
}

func TestPriorityQueueSuite(t *testing.T) {
	suite.Run(t, new(priorityQueueTestSuite))
}

func (s *priorityQueueTestSuite) Test_PopPushMultipleDifferentObjects() {

	cl, err := Open(s.Ctx, s.config, s.L)
	s.NoError(err)
	defer cl.Close()

	type QItem struct {
		S string   `json:"s"`
		I int      `json:"i"`
		F *float64 `json:"f"`
	}

	queue := kit.NewRandString()

	q := NewPriorityQueue[*QItem](queue, cl)

	// push
	s.NoError(q.Push(s.Ctx, &QItem{"1", 1, kit.Float64Ptr(0.1)}, 5))
	s.NoError(q.Push(s.Ctx, &QItem{"2", 1, kit.Float64Ptr(0.1)}, 4))
	s.NoError(q.Push(s.Ctx, &QItem{"3", 1, kit.Float64Ptr(0.1)}, 4))
	s.NoError(q.Push(s.Ctx, &QItem{"4", 1, kit.Float64Ptr(0.1)}, 3))
	s.NoError(q.Push(s.Ctx, &QItem{"5", 1, kit.Float64Ptr(0.1)}, 1))

	// pop
	items, err := q.Pop(s.Ctx, 2)
	s.NoError(err)
	s.Len(items, 2)
	kit.ForAll(items, func(i *QItem) {
		s.NotEmpty(i)
		s.NotEmpty(i.I)
		s.NotEmpty(i.S)
	})
	s.Equal(items[0].S, "5")
	s.Equal(items[1].S, "4")

	// pop & remove
	items, err = q.PopManyAndRemove(s.Ctx, 5)
	s.NoError(err)
	s.Len(items, 5)

	// check empty
	items, err = q.Pop(s.Ctx, 2)
	s.NoError(err)
	s.Empty(items)

}

func (s *priorityQueueTestSuite) Test_PopPushMultipleSameObjects() {

	cl, err := Open(s.Ctx, s.config, s.L)
	s.NoError(err)
	defer cl.Close()

	type QItem struct {
		S string   `json:"s"`
		I int      `json:"i"`
		F *float64 `json:"f"`
	}

	queue := kit.NewRandString()

	q := NewPriorityQueue[*QItem](queue, cl)

	// push
	s.NoError(q.Push(s.Ctx, &QItem{"1", 1, kit.Float64Ptr(0.1)}, 5))
	s.NoError(q.Push(s.Ctx, &QItem{"1", 1, kit.Float64Ptr(0.1)}, 5))
	s.NoError(q.Push(s.Ctx, &QItem{"1", 1, kit.Float64Ptr(0.1)}, 5))

	// pop
	items, err := q.Pop(s.Ctx, 3)
	s.NoError(err)
	s.Len(items, 3)

}

func (s *priorityQueueTestSuite) Test_PopPushSimpleType() {

	cl, err := Open(s.Ctx, s.config, s.L)
	s.NoError(err)
	defer cl.Close()

	queue := kit.NewRandString()

	q := NewPriorityQueue[string](queue, cl)

	// push
	s.NoError(q.Push(s.Ctx, "1", 2))
	s.NoError(q.Push(s.Ctx, "2", 1))

	// pop
	items, err := q.Pop(s.Ctx, 3)
	s.NoError(err)
	s.Len(items, 2)
	s.Equal(items[0], "2")
	s.Equal(items[1], "1")

}

func (s *priorityQueueTestSuite) Test_PopAndRemoveWhenEmpty() {

	cl, err := Open(s.Ctx, s.config, s.L)
	s.NoError(err)
	defer cl.Close()

	type QItem struct {
		S string   `json:"s"`
		I int      `json:"i"`
		F *float64 `json:"f"`
	}

	queue := kit.NewRandString()

	q := NewPriorityQueue[*QItem](queue, cl)

	// push
	s.NoError(q.Push(s.Ctx, &QItem{"1", 1, kit.Float64Ptr(0.1)}, 5))
	s.NoError(q.Push(s.Ctx, &QItem{"2", 1, kit.Float64Ptr(0.1)}, 4))

	// pop & remove
	items, err := q.PopManyAndRemove(s.Ctx, 5)
	s.NoError(err)
	s.Len(items, 2)

	// check empty
	items, err = q.PopManyAndRemove(s.Ctx, 5)
	s.NoError(err)
	s.Len(items, 0)

}
