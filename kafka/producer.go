package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/algmib/kit"
)

const (
	defaultRetryTimes   = 3
	defaultRetryTimeout = time.Second
)

// Producer allows sending message to broker
type Producer interface {
	// Send sends a message to broker
	Send(ctx context.Context, key string, payload interface{}) error
	// SendMany sends bulk of messages to broker
	SendMany(ctx context.Context, messages ...*Message) error
}

type producerImpl struct {
	topic           *TopicConfig
	writer          *kafka.Writer
	logger          kit.CLoggerFunc
	cancellationCtx context.Context
	retryTimes      int
	retryTimeout    time.Duration
}

func (p *producerImpl) l() kit.CLogger {
	return p.logger().Cmp("kafka-producer")
}

func newProducer(ctx context.Context, logger kit.CLoggerFunc, topic *TopicConfig, cfg *ProducerConfig, urls []string, transport *kafka.Transport) Producer {

	// populate writer params
	writer := &kafka.Writer{
		Addr:        kafka.TCP(urls...),
		Topic:       topic.Topic,
		ErrorLogger: kafka.LoggerFunc(logger().Mth("producer").F(kit.KV{"topic": topic.Topic}).PrintfErr),
		Transport:   transport,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				logger().Mth("producer-completion").F(kit.KV{"topic": topic.Topic}).E(ErrKafkaProduceMsg(ctx, err)).Err()
			}
		},
		Balancer: &kafka.Hash{}, // set up hash balancer to guaranty the messages with the same key are always in the same partition
		Async:    cfg.Async,
	}
	if cfg.BatchSize != nil {
		writer.BatchSize = *cfg.BatchSize
	}
	if cfg.BatchTimeout != nil {
		writer.BatchTimeout = *cfg.BatchTimeout
	}
	if cfg.RequiredAcks != nil {
		writer.RequiredAcks = kafka.RequiredAcks(*cfg.RequiredAcks)
	}
	if cfg.MaxAttempts != nil {
		writer.MaxAttempts = *cfg.MaxAttempts
	}

	r := &producerImpl{
		writer:          writer,
		logger:          logger,
		topic:           topic,
		cancellationCtx: ctx,
	}

	if cfg.RetryTimes != nil {
		r.retryTimes = *cfg.RetryTimes
	} else {
		r.retryTimes = defaultRetryTimes
	}
	if cfg.RetryTimeout != nil {
		r.retryTimeout = *cfg.RetryTimeout
	} else {
		r.retryTimeout = defaultRetryTimeout
	}

	return r
}

func (p *producerImpl) Send(ctx context.Context, key string, payload interface{}) error {
	l := p.l().Mth("publish").F(kit.KV{"topic": p.topic.Topic}).Dbg()

	// prepare message
	ctxRq, err := p.rqCtx(ctx)
	if err != nil {
		return err
	}

	msg := &Message{
		Ctx:     ctxRq,
		Payload: payload,
		Key:     key,
	}

	// write message
	err = p.sendWithRetry(ctx, msg)
	if err != nil {
		return err
	}

	l.Dbg("ok").TrcObj("%+v", msg)

	return nil
}

func (p *producerImpl) SendMany(ctx context.Context, messages ...*Message) error {
	l := p.l().Mth("send-many").F(kit.KV{"topic": p.topic.Topic}).Dbg()

	// prepare message
	ctxRq, err := p.rqCtx(ctx)
	if err != nil {
		return err
	}

	for _, m := range messages {
		m.Ctx = ctxRq
	}

	// write message
	err = p.sendWithRetry(ctx, messages...)
	if err != nil {
		return err
	}

	l.Dbg("ok").TrcObj("%+v", messages)

	return nil
}

func (p *producerImpl) rqCtx(ctx context.Context) (*kit.RequestContext, error) {
	if rCtx, ok := kit.Request(ctx); ok {
		return rCtx, nil
	}
	return nil, ErrKafkaMessageContextInvalid(ctx, p.topic.Topic)
}

func (p *producerImpl) sendWithRetry(ctx context.Context, messages ...*Message) error {
	messagesToSend := make([]kafka.Message, 0, len(messages))
	now := kit.Now()
	for _, msg := range messages {

		m, err := kit.Marshal(msg)
		if err != nil {
			return ErrKafkaMessageMarshal(ctx, err, p.topic.Topic)
		}

		messagesToSend = append(messagesToSend, kafka.Message{
			Key:   []byte(msg.Key),
			Value: m,
			Time:  now,
		})
	}

	// send with retry
	for i := 0; i < p.retryTimes; i++ {
		err := p.writer.WriteMessages(p.cancellationCtx, messagesToSend...)
		if err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) {
				time.Sleep(p.retryTimeout)
				continue
			} else {
				return ErrKafkaMessageWrite(ctx, err, p.topic.Topic)
			}
		}
		break
	}
	return nil
}
