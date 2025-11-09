package kafka

import (
	"context"
	"time"

	"github.com/mikhailbolshakov/kit"
	"github.com/segmentio/kafka-go"
)

const (
	subWorkersPerTopic            = 4
	waitPeriodBeforeReaderRestart = time.Second * 30
)

type subscriberStrategy interface {
	start(ctx context.Context, topic string)
}

type subscriber struct {
	readerCfg *kafka.ReaderConfig
	handlers  []HandlerFn
	workers   int
	logger    kit.CLoggerFunc
	strategy  subscriberStrategy
}

func (s *subscriber) l() kit.CLogger {
	return s.logger().Cmp("kafka-sub")
}

func newSubscriber(logger kit.CLoggerFunc, topic *TopicConfig, cfg *SubscriberConfig, urls []string, dialer *kafka.Dialer, handlers ...HandlerFn) *subscriber {

	// setup reader
	readerCfg := &kafka.ReaderConfig{
		Brokers:     urls,
		GroupID:     cfg.GroupId,
		Topic:       topic.Topic,
		Dialer:      dialer,
		ErrorLogger: kafka.LoggerFunc(logger().Mth("subscriber").F(kit.KV{"topic": topic.Topic, "groupId": cfg.GroupId}).PrintfErr),
	}
	if cfg.CommitInterval != nil {
		readerCfg.CommitInterval = *cfg.CommitInterval
	}
	if cfg.BatchTimeout != nil {
		readerCfg.ReadBatchTimeout = *cfg.BatchTimeout
	}
	if cfg.MaxAttempts != nil {
		readerCfg.MaxAttempts = *cfg.MaxAttempts
	}
	if cfg.MaxWait != nil {
		readerCfg.MaxWait = *cfg.MaxWait
	}
	if cfg.JoinGroupBackoff != nil {
		readerCfg.JoinGroupBackoff = *cfg.JoinGroupBackoff
	}
	if cfg.StartOffset != nil {
		readerCfg.StartOffset = *cfg.StartOffset
	} else {
		readerCfg.StartOffset = kafka.LastOffset
	}
	if cfg.Logging {
		readerCfg.Logger = kafka.LoggerFunc(logger().Mth("subscriber").F(kit.KV{"topic": topic.Topic, "groupId": cfg.GroupId}).Printf)
	}

	// subscriber
	sub := &subscriber{
		readerCfg: readerCfg,
		handlers:  handlers,
		workers:   subWorkersPerTopic,
		logger:    logger,
	}

	if cfg.Workers != nil {
		sub.workers = *cfg.Workers
	}

	if sub.manualCommit() {
		sub.strategy = newSubscriberManualCommitStrategy(logger, readerCfg, handlers, cfg.ManualCommit, cfg.DLQProducer, sub.workers)
	} else {
		sub.strategy = newSubscriberAutoCommitStrategy(logger, readerCfg, handlers, sub.workers)
	}

	return sub
}

func (s *subscriber) manualCommit() bool {
	return s.readerCfg.CommitInterval == 0
}

func (s *subscriber) start(ctx context.Context, topic string) {
	s.l().C(ctx).Mth("start").F(kit.KV{"topic": topic}).Dbg()
	s.strategy.start(ctx, topic)
}

func (s *subscriber) close() {}
