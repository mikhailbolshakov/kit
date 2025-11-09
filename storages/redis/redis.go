package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.com/algmib/kit"
	"net"
	"time"
)

const (
	NotFound = redis.Nil
)

type Redis struct {
	Instance *redis.Client
	Ttl      time.Duration
	logger   kit.CLoggerFunc
}

// Config redis config
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Db       int
	Ttl      uint
}

func (r *Redis) l() kit.CLogger {
	return r.logger().Cmp("redis")
}

func Open(ctx context.Context, params *Config, logger kit.CLoggerFunc) (*Redis, error) {

	l := logger().Cmp("redis").Mth("open")

	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(params.Host, params.Port),
		Username: params.Username,
		Password: params.Password,
		DB:       params.Db,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, ErrRedisPingErr(err)
	}

	l.Inf("ok")
	return &Redis{
		Instance: client,
		Ttl:      time.Duration(params.Ttl) * time.Second,
		logger:   logger,
	}, nil
}

func (r *Redis) Close() {
	l := r.l().Mth("close")
	if r.Instance != nil {
		_ = r.Instance.Close()
	}
	l.Inf("ok")
}
