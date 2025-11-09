//go:build integration

package redis

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"testing"
	"time"
)

type redisTestSuite struct {
	kit.Suite
}

func (s *redisTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(redisTestSuite))
}

var (
	config = &Config{
		Host: "localhost",
		Port: "6379",
		Ttl:  0,
	}
)

func (s *redisTestSuite) Test_Range() {

	cl, err := Open(s.Ctx, config, s.L)
	s.NoError(err)
	defer cl.Close()

	key := kit.NewRandString()
	jsons, err := cl.Instance.LRange(s.Ctx, key, 0, -1).Result()
	s.NoError(err)
	fmt.Println(jsons)

	pipe := cl.Instance.Pipeline()
	pipe.Expire(s.Ctx, key, time.Second*10)

	s.NoError(cl.Instance.RPush(s.Ctx, key, "1").Err())
	s.NoError(cl.Instance.RPush(s.Ctx, key, "2").Err())
	s.NoError(cl.Instance.RPush(s.Ctx, key, "3").Err())

	_, err = pipe.Exec(s.Ctx)
	s.NoError(err)

	jsons, err = cl.Instance.LRange(s.Ctx, key, 0, -1).Result()
	s.NoError(err)
	s.Equal(3, len(jsons))
}

func (s *redisTestSuite) Test_Distributed_Lock() {
	cl, err := Open(s.Ctx, config, s.L)
	s.NoError(err)
	defer cl.Close()

	key, unlockId := kit.NewRandString(), kit.NewRandString()

	// apply lock
	locked, err := cl.Lock(s.Ctx, key, unlockId, time.Second*10)
	s.NoError(err)
	s.True(locked)

	// apply lock again
	locked, err = cl.Lock(s.Ctx, key, unlockId, time.Second*10)
	s.NoError(err)
	s.False(locked)

	// try to lock with another unlockId
	locked, err = cl.Lock(s.Ctx, key, kit.NewRandString(), time.Second*10)
	s.NoError(err)
	s.False(locked)

	// try to unlock with another unlock ID
	unlocked, err := cl.UnLock(s.Ctx, key, kit.NewRandString())
	s.NoError(err)
	s.False(unlocked)

	// try to unlock with another unlock ID
	unlocked, err = cl.UnLock(s.Ctx, key, unlockId)
	s.NoError(err)
	s.True(unlocked)

}

func (s *redisTestSuite) Test_Json() {

	s.T().Skip("Redis 8 support")

	cl, err := Open(s.Ctx, config, s.L)
	s.NoError(err)
	defer cl.Close()

	type Some struct {
		A      string `json:"a"`
		B      int    `json:"b"`
		Nested struct {
			Slice []string          `json:"slice"`
			Map   map[string]string `json:"map"`
		} `json:"nested"`
	}

	key := kit.NewRandString()
	some := Some{
		A: "a",
		B: 1,
	}
	some.Nested.Slice = []string{"a", "b", "c"}
	some.Nested.Map = map[string]string{"a": "b"}

	// set json
	someJs, _ := kit.JsonEncode(some)
	s.NoError(cl.Instance.JSONSet(s.Ctx, key, ".", someJs).Err())

	// get json
	rsTxt, err := cl.Instance.JSONGet(s.Ctx, key, ".").Result()
	s.NoError(err)
	s.NotEmpty(rsTxt)
	resSome, err := kit.JsonDecode[Some]([]byte(rsTxt))
	s.NoError(err)
	s.NotNil(resSome)
	s.Equal(some, *resSome)

	// get attribute
	rsAttrTxt, err := cl.Instance.JSONGet(s.Ctx, key, "$.a").Result()
	s.NoError(err)
	resAttr, err := kit.JsonDecodePlainSlice[string]([]byte(rsAttrTxt))
	s.NoError(err)
	s.NotEmpty(resAttr)
	s.Equal(some.A, resAttr[0])

	// update
	some.B = 2
	some.A = "another"
	some.Nested.Slice = []string{"a", "b", "c", "d"}
	s.NoError(cl.Instance.JSONSet(s.Ctx, key, "$.b", some.B).Err())
	s.NoError(cl.Instance.JSONSet(s.Ctx, key, "$.a", fmt.Sprintf("\"%s\"", some.A)).Err())
	s.NoError(cl.Instance.JSONSet(s.Ctx, key, "$.nested.slice", some.Nested.Slice).Err())

	// get json
	rsTxt, err = cl.Instance.JSONGet(s.Ctx, key, ".").Result()
	s.NoError(err)
	s.NotEmpty(rsTxt)
	resSome, err = kit.JsonDecode[Some]([]byte(rsTxt))
	s.NoError(err)
	s.NotNil(resSome)
	s.Equal(some, *resSome)

}
