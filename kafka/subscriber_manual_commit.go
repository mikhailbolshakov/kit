package kafka

import (
	"context"
	stdErr "errors"
	"io"
	"time"

	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/goroutine"
	"github.com/segmentio/kafka-go"
)

const (
	commitMessageMaxRetryCount      = 5
	commitMessageRetryBackoffStepMs = 500
	handleMessageMaxRetryCount      = 5
	handleMessageRetryBackoffStepMs = 500
)

type SubscriberManualCommitConfig struct {
	CommitMessageMaxRetryCount      int
	CommitMessageRetryBackoffStepMs int
	HandleMessageMaxRetryCount      int
	HandleMessageRetryBackoffStepMs int
}

type subscriberManualCommit struct {
	logger          kit.CLoggerFunc
	readerCfg       *kafka.ReaderConfig
	handlers        []HandlerFn
	manualCommitCfg *SubscriberManualCommitConfig
	dlqProducer     Producer
	workers         int
}

func newSubscriberManualCommitStrategy(logger kit.CLoggerFunc,
	readerCfg *kafka.ReaderConfig,
	handlers []HandlerFn,
	manualCommitCfg *SubscriberManualCommitConfig,
	dlqProducer Producer,
	workers int) subscriberStrategy {

	if manualCommitCfg == nil {
		manualCommitCfg = &SubscriberManualCommitConfig{
			CommitMessageMaxRetryCount:      commitMessageMaxRetryCount,
			CommitMessageRetryBackoffStepMs: commitMessageRetryBackoffStepMs,
			HandleMessageMaxRetryCount:      handleMessageMaxRetryCount,
			HandleMessageRetryBackoffStepMs: handleMessageRetryBackoffStepMs,
		}
	}

	if manualCommitCfg.CommitMessageMaxRetryCount < 0 {
		manualCommitCfg.CommitMessageMaxRetryCount = commitMessageMaxRetryCount
	}
	if manualCommitCfg.CommitMessageRetryBackoffStepMs < 0 {
		manualCommitCfg.CommitMessageRetryBackoffStepMs = commitMessageRetryBackoffStepMs
	}
	if manualCommitCfg.HandleMessageMaxRetryCount < 0 {
		manualCommitCfg.HandleMessageMaxRetryCount = handleMessageMaxRetryCount
	}
	if manualCommitCfg.HandleMessageRetryBackoffStepMs < 0 {
		manualCommitCfg.HandleMessageRetryBackoffStepMs = handleMessageRetryBackoffStepMs
	}

	return &subscriberManualCommit{
		logger:          logger,
		readerCfg:       readerCfg,
		handlers:        handlers,
		manualCommitCfg: manualCommitCfg,
		dlqProducer:     dlqProducer,
		workers:         workers,
	}
}

func (s *subscriberManualCommit) l() kit.CLogger {
	return s.logger().Cmp("kafka-manual-sub")
}

func (s *subscriberManualCommit) start(ctx context.Context, topic string) {
	s.l().C(ctx).Mth("start").F(kit.KV{"topic": topic}).Dbg()

	reader := kafka.NewReader(*s.readerCfg)

	// start goroutine to fetch messages
	goroutine.New().
		WithLogger(s.l().Mth("fetch")).
		WithRetry(goroutine.Unrestricted).
		Go(ctx,
			func() {

				// close reader (may take some time)
				defer func() { _ = reader.Close() }()

				// run workers
				workersChannels := make([]chan kafka.Message, s.workers)
				for i := 0; i < s.workers; i++ {
					workersChannels[i] = make(chan kafka.Message, workersChanCapacity)
					s.subscriberWorker(ctx, reader, topic, i, workersChannels[i])
				}

				// close all worker channels
				defer kit.ForAll(workersChannels, func(c chan kafka.Message) { close(c) })

				l := s.l().C(ctx).Mth("fetch").F(kit.KV{"topic": topic}).Dbg("started")
				for {

					// check if context is already cancelled
					if ctx.Err() != nil {
						l.Dbg("stopped")
						return
					}

					// read message
					m, err := reader.FetchMessage(ctx)
					if err != nil {

						// reader has been closed, restart
						if stdErr.Is(err, io.EOF) || stdErr.Is(err, io.ErrUnexpectedEOF) {
							l.Dbg("EOF -> restart")
							time.AfterFunc(waitPeriodBeforeReaderRestart, func() { s.start(ctx, topic) })
							return
						}

						s.l().Mth("fetch").F(kit.KV{"topic": topic}).E(ErrKafkaFetchMessage(err)).Err("fetch")
						continue
					}

					l.DbgF("key: %s", string(m.Key)).TrcF("%s", string(m.Value))

					// send a message to the channel to be processed by workers
					if len(m.Value) != 0 && len(m.Key) != 0 {
						// send message to proper channel
						workersChannels[s.chanIndexByPartition(m.Partition)] <- m
					}

				}
			},
		)

}

func (s *subscriberManualCommit) subscriberWorker(ctx context.Context, reader *kafka.Reader, topic string, workerTag int, receiverChan chan kafka.Message) {

	goroutine.New().
		WithLogger(s.l().Mth("sub-worker")).
		WithRetry(goroutine.Unrestricted).
		Go(ctx,
			func() {
				l := s.l().Mth("worker").F(kit.KV{"tag": workerTag, "topic": topic}).Dbg("started")
				for {
					select {
					case msg, ok := <-receiverChan:

						if !ok {
							l.Dbg("closed")
							return
						}

						l.DbgF("key: %s", string(msg.Key)).TrcF("%s", string(msg.Value))

						// run handler
						if err := s.handleWithRetry(ctx, topic, msg); err != nil {
							s.l().C(ctx).Mth("handler").F(kit.KV{"topic": topic, "key": msg.Key}).E(err).St().Err()

							// try to send to DLQ
							if !s.dlq(ctx, topic, msg) {
								// if DLQ is not configured, just skip the message without committing
								// otherwise commit the message and continue
								continue
							}

						}

						// commit with retry
						if err := s.commitWithRetry(ctx, topic, reader, msg); err != nil {
							s.l().C(ctx).Mth("commit").F(kit.KV{"topic": topic, "key": msg.Key}).E(err).St().Err()
						}

					case <-ctx.Done():
						l.Dbg("stopped")
						return
					}
				}
			},
		)
}

func (s *subscriberManualCommit) commitWithRetry(ctx context.Context, topic string, reader *kafka.Reader, message kafka.Message) error {
	l := s.l().C(ctx).Mth("commit").F(kit.KV{"topic": topic})

	for attempt := 0; attempt < s.manualCommitCfg.CommitMessageMaxRetryCount; attempt++ {
		if err := reader.CommitMessages(ctx, message); err != nil {

			l.E(err).St().ErrF("attempt %d failed", attempt+1)

			if attempt < s.manualCommitCfg.CommitMessageMaxRetryCount-1 {
				// Exponential backoff: 100ms, 200ms, 400ms...
				time.Sleep(time.Duration(s.manualCommitCfg.CommitMessageRetryBackoffStepMs*(1<<attempt)) * time.Millisecond)
				continue
			}

			// retries exceeded
			return ErrKafkaManualCommitRetryCountExceeded(ctx)
		}

	}

	return nil
}

func (s *subscriberManualCommit) handleWithRetry(ctx context.Context, topic string, m kafka.Message) error {
	l := s.l().C(ctx).Mth("handle").F(kit.KV{"topic": topic})

	handlerFn := func() error {
		for _, handler := range s.handlers {
			if err := handler(m.Value); err != nil {
				return err
			}
		}
		return nil
	}

	for attempt := 0; attempt < s.manualCommitCfg.HandleMessageMaxRetryCount; attempt++ {
		if err := handlerFn(); err != nil {

			l.E(err).St().ErrF("attempt %d failed", attempt+1)

			if attempt < s.manualCommitCfg.HandleMessageMaxRetryCount-1 {
				// Exponential backoff: 100ms, 200ms, 400ms...
				time.Sleep(time.Duration(s.manualCommitCfg.HandleMessageRetryBackoffStepMs*(1<<attempt)) * time.Millisecond)
				continue
			}

			// retries exceeded
			return ErrKafkaHandleMessageManualCommitRetryCountExceeded(ctx)
		}

	}

	return nil

}

func (s *subscriberManualCommit) dlq(ctx context.Context, topic string, m kafka.Message) bool {
	l := s.l().C(ctx).Mth("dlq").F(kit.KV{"topic": topic, "key": m.Key})

	if s.dlqProducer == nil {
		return false
	}

	err := s.dlqProducer.Send(ctx, string(m.Key), &DLQMessage{
		Topic:         topic,
		FailedMessage: m.Value,
	})
	if err != nil {
		l.E(err).St().Err()
		return false
	}

	l.Dbg("sent")
	return true

}

func (s *subscriberManualCommit) chanIndexByPartition(partition int) int {
	return partition % s.workers
}
