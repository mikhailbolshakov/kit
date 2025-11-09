package kafka

import "time"

// ProducerConfig specifies producer config params
// use builder rather than manual population
type ProducerConfig struct {
	BatchSize    *int
	BatchTimeout *time.Duration
	RequiredAcks *int
	MaxAttempts  *int
	Async        bool
	RetryTimes   *int
	RetryTimeout *time.Duration
}

type ProducerConfigBuilder interface {
	// BatchTimeout sets batch timeout value (default: 1s)
	BatchTimeout(to time.Duration) ProducerConfigBuilder
	// BatchSize sets batch size (default: 100)
	BatchSize(size int) ProducerConfigBuilder
	// Async if true, WriteMessages call will never block but errors aren't returned (default: false)
	Async(v bool) ProducerConfigBuilder
	// Retry sets retry params (default: time=3, timeout = 1s)
	Retry(time int, timeout time.Duration) ProducerConfigBuilder
	// RequiredAcks sets required acks value (0(None) default, 1(Single), -1(All))
	RequiredAcks(v int) ProducerConfigBuilder
	// Build builds config
	Build() *ProducerConfig
}

type producerConfigBuilder struct {
	cfg *ProducerConfig
}

func NewProducerCfgBuilder() ProducerConfigBuilder {
	return &producerConfigBuilder{
		cfg: &ProducerConfig{},
	}
}

func (p *producerConfigBuilder) Retry(times int, timeout time.Duration) ProducerConfigBuilder {
	p.cfg.RetryTimes = &times
	p.cfg.RetryTimeout = &timeout
	return p
}

func (p *producerConfigBuilder) RequiredAcks(v int) ProducerConfigBuilder {
	p.cfg.RequiredAcks = &v
	return p
}

func (p *producerConfigBuilder) BatchTimeout(to time.Duration) ProducerConfigBuilder {
	p.cfg.BatchTimeout = &to
	return p
}

func (p *producerConfigBuilder) BatchSize(size int) ProducerConfigBuilder {
	p.cfg.BatchSize = &size
	return p
}

func (p *producerConfigBuilder) Async(v bool) ProducerConfigBuilder {
	p.cfg.Async = v
	return p
}

func (p *producerConfigBuilder) Build() *ProducerConfig {
	return p.cfg
}
