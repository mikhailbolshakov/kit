package kafka

import (
	"context"
	"time"

	"gitlab.com/algmib/kit"
)

// SubscriberConfig specifies subscriber config params
// use builder rather than manual population
type SubscriberConfig struct {
	GroupId          string                        // allows load balancing messages within the same group
	BatchTimeout     *time.Duration                // timeout of batch fetching from kafka (default: 10s)
	MaxWait          *time.Duration                // maximum amount of time to wait for new data to come when fetching batches (default: 10s)
	CommitInterval   *time.Duration                // interval to commit to kafka (default: sync)
	Workers          *int                          // number of workers (default: 4)
	MaxAttempts      *int                          // maximum number of attempts to process a message
	StartOffset      *int64                        // determines from which offset a new group starts to consume (-2:FirstOffset, -1:LastOffset)
	JoinGroupBackoff *time.Duration                // length of time to wait between re-joining
	Logging          bool                          // if true subscriber logging enabled
	ManualCommit     *SubscriberManualCommitConfig // configuration for manual commit behavior
	DLQProducer      Producer                      // dead-letter queue producer for handling failed messages
}

type SubscriberConfigBuilder interface {
	// GroupId allows load balancing messages within the same group
	GroupId(groupId string) SubscriberConfigBuilder
	// BatchTimeout sets timeout of batch fetching from kafka (default: 10s)
	BatchTimeout(to time.Duration) SubscriberConfigBuilder
	// MaxWait sets maximum amount of time to wait for new data to come when fetching batches (default: 10s)
	MaxWait(to time.Duration) SubscriberConfigBuilder
	// CommitInterval sets interval to commit to kafka (default: sync)
	CommitInterval(to time.Duration) SubscriberConfigBuilder
	// Workers sets number of workers (default: 4)
	Workers(num int) SubscriberConfigBuilder
	// StartOffset determines from which offset a new group starts to consume. it must be set to one of FirstOffset = -2 or LastOffset = -1 (Default: FirstOffset)
	// Only used when GroupID is set
	StartOffset(v int64) SubscriberConfigBuilder
	// JoinGroupBackoff optionally sets the length of time to wait between re-joining
	JoinGroupBackoff(t time.Duration) SubscriberConfigBuilder
	// Logging if true subscriber logging enabled
	Logging(v bool) SubscriberConfigBuilder
	// ManualCommitMessageMaxRetryCount sets max retry count for manual commit retry if an error occurs.
	ManualCommitMessageMaxRetryCount(v int) SubscriberConfigBuilder
	// ManualCommitMessageRetryBackoffStepMs sets backoff delay step in ms between retries for manual commit retry if an error occurs.
	ManualCommitMessageRetryBackoffStepMs(v int) SubscriberConfigBuilder
	// ManualCommitHandleMessageMaxRetryCount sets max retry count for handling message retry if an error occurs.
	ManualCommitHandleMessageMaxRetryCount(v int) SubscriberConfigBuilder
	// ManualCommitHandleMessageRetryBackoffStepMs sets backoff delay step in ms between retries for handling message retry if an error occurs.
	ManualCommitHandleMessageRetryBackoffStepMs(v int) SubscriberConfigBuilder
	// DLQProducer sets a dead-letter queue producer for handling failed messages and returns the modified builder instance.
	DLQProducer(p Producer) SubscriberConfigBuilder
	// Validate checks the current configuration for any errors or missing mandatory fields and returns an error if invalid.
	Validate(ctx context.Context) error
	// Build builds config
	Build() *SubscriberConfig
}

type subscriberConfigBuilder struct {
	cfg *SubscriberConfig
}

func NewSubscriberCfgBuilder() SubscriberConfigBuilder {
	w := subWorkersPerTopic
	return &subscriberConfigBuilder{
		cfg: &SubscriberConfig{
			Workers: &w,
		},
	}
}

func (p *subscriberConfigBuilder) MaxWait(to time.Duration) SubscriberConfigBuilder {
	p.cfg.MaxWait = &to
	return p
}

func (p *subscriberConfigBuilder) GroupId(groupId string) SubscriberConfigBuilder {
	p.cfg.GroupId = groupId
	return p
}

func (p *subscriberConfigBuilder) CommitInterval(to time.Duration) SubscriberConfigBuilder {
	p.cfg.CommitInterval = &to
	return p
}

func (p *subscriberConfigBuilder) Workers(num int) SubscriberConfigBuilder {
	p.cfg.Workers = &num
	return p
}

func (p *subscriberConfigBuilder) BatchTimeout(to time.Duration) SubscriberConfigBuilder {
	p.cfg.BatchTimeout = &to
	return p
}

func (p *subscriberConfigBuilder) StartOffset(v int64) SubscriberConfigBuilder {
	p.cfg.StartOffset = &v
	return p
}

func (p *subscriberConfigBuilder) JoinGroupBackoff(t time.Duration) SubscriberConfigBuilder {
	p.cfg.JoinGroupBackoff = &t
	return p
}

func (p *subscriberConfigBuilder) Logging(v bool) SubscriberConfigBuilder {
	p.cfg.Logging = v
	return p
}

func (p *subscriberConfigBuilder) ManualCommitMessageMaxRetryCount(v int) SubscriberConfigBuilder {
	if p.cfg.ManualCommit == nil {
		p.cfg.ManualCommit = &SubscriberManualCommitConfig{}
	}
	p.cfg.ManualCommit.CommitMessageMaxRetryCount = v
	return p
}

func (p *subscriberConfigBuilder) ManualCommitMessageRetryBackoffStepMs(v int) SubscriberConfigBuilder {
	if p.cfg.ManualCommit == nil {
		p.cfg.ManualCommit = &SubscriberManualCommitConfig{}
	}
	p.cfg.ManualCommit.CommitMessageRetryBackoffStepMs = v
	return p
}

func (p *subscriberConfigBuilder) ManualCommitHandleMessageMaxRetryCount(v int) SubscriberConfigBuilder {
	if p.cfg.ManualCommit == nil {
		p.cfg.ManualCommit = &SubscriberManualCommitConfig{}
	}
	p.cfg.ManualCommit.HandleMessageMaxRetryCount = v
	return p
}

func (p *subscriberConfigBuilder) ManualCommitHandleMessageRetryBackoffStepMs(v int) SubscriberConfigBuilder {
	if p.cfg.ManualCommit == nil {
		p.cfg.ManualCommit = &SubscriberManualCommitConfig{}
	}
	p.cfg.ManualCommit.HandleMessageRetryBackoffStepMs = v
	return p
}

func (p *subscriberConfigBuilder) DLQProducer(v Producer) SubscriberConfigBuilder {
	p.cfg.DLQProducer = v
	return p
}

func (p *subscriberConfigBuilder) Validate(ctx context.Context) error {

	// auto commit
	if p.cfg.CommitInterval != nil && *p.cfg.CommitInterval > 0 {
		if p.cfg.ManualCommit != nil {
			return ErrKafkaSubscriberConfigInvalid(ctx, "manual commit and auto commit are mutually exclusive")
		}
		if p.cfg.DLQProducer != nil {
			return ErrKafkaSubscriberConfigInvalid(ctx, "dead-letter queue producer not suppoerted for auto commit")
		}
	} else {
		if p.cfg.ManualCommit == nil {
			p.cfg.ManualCommit = &SubscriberManualCommitConfig{}
		}
	}

	return nil
}

func (p *subscriberConfigBuilder) Build() *SubscriberConfig {
	if p.cfg.GroupId == "" {
		p.cfg.GroupId = kit.NewRandString()
	}
	if p.cfg.Workers == nil || *p.cfg.Workers < 1 {
		p.cfg.Workers = kit.IntPtr(1)
	}
	return p.cfg
}
