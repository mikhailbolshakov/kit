package kafka

import (
	"context"
	stdErr "errors"
	"hash"
	"hash/fnv"
	"io"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/algmib/kit"
	"gitlab.com/algmib/kit/goroutine"
)

const (
	workersChanCapacity = 255
)

type subscriberAutoCommit struct {
	logger    kit.CLoggerFunc
	readerCfg *kafka.ReaderConfig
	handlers  []HandlerFn
	workers   int
}

func newSubscriberAutoCommitStrategy(logger kit.CLoggerFunc,
	readerCfg *kafka.ReaderConfig,
	handlers []HandlerFn,
	workers int) subscriberStrategy {
	return &subscriberAutoCommit{
		logger:    logger,
		readerCfg: readerCfg,
		handlers:  handlers,
		workers:   workers,
	}
}

func (s *subscriberAutoCommit) l() kit.CLogger {
	return s.logger().Cmp("kafka-auto-sub")
}

func (s *subscriberAutoCommit) start(ctx context.Context, topic string) {
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
					s.subscriberWorker(ctx, topic, s.handlers, i, workersChannels[i])
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
					m, err := reader.ReadMessage(ctx)
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
					l.TrcObj("%+v", m)

					// send a message to the channel to be processed by workers
					if len(m.Value) != 0 && len(m.Key) != 0 {
						// send message to proper channel
						workersChannels[s.chanIndexByKey(m.Key)] <- m
					}
				}
			},
		)

}

func (s *subscriberAutoCommit) subscriberWorker(ctx context.Context, topic string, handlers []HandlerFn, workerTag int, receiverChan chan kafka.Message) {

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
						for _, handler := range handlers {
							if err := handler(msg.Value); err != nil {
								s.l().C(ctx).Mth("handler").F(kit.KV{"topic": topic, "key": msg.Key}).E(err).St().Err()
							}
						}

					case <-ctx.Done():
						l.Dbg("stopped")
						return
					}
				}
			},
		)
}

var (
	fnv1aPool = &sync.Pool{
		New: func() interface{} {
			return fnv.New32a()
		},
	}
)

// chanIndexByKey calculates index in channel slice by hashing message key
func (s *subscriberAutoCommit) chanIndexByKey(key []byte) int {

	h := fnv1aPool.Get().(hash.Hash32)
	defer fnv1aPool.Put(h)

	h.Reset()
	_, _ = h.Write(key)

	ind := int32(h.Sum32()) % int32(s.workers)
	if ind < 0 {
		ind = -ind
	}

	return int(ind)
}
