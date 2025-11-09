package kit

import (
	"context"
	"time"
)

const (
	ErrCodeDistributedLockFailed = "DLK-001"
)

var (
	ErrLock = func(ctx context.Context, err error, ref string) error {
		return NewAppErrBuilder(ErrCodeDistributedLockFailed, "failed to acquire lock").Wrap(err).C(ctx).F(KV{"ref": ref}).Err()
	}
)

type DistributedLockStorage interface {
	// Lock attempts to acquire a lock on an employee record
	Lock(ctx context.Context, ref, releaseId string) (bool, error)
	// UnLock releases a previously acquired lock
	UnLock(ctx context.Context, ref, releaseId string) error
}

type DistributedLock interface {
	Lock(ctx context.Context, ref string) (string, error)
	UnLock(ctx context.Context, ref, releaseId string)
}

type DistributedLockCfg struct {
	AwaitPeriod time.Duration
}

type distributedLockSvcImpl struct {
	storage DistributedLockStorage
	cfg     *DistributedLockCfg
	logger  CLoggerFunc
}

func NewDistributedLock(storage DistributedLockStorage, cfg *DistributedLockCfg, logger CLoggerFunc) DistributedLock {
	return &distributedLockSvcImpl{
		storage: storage,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *distributedLockSvcImpl) l() CLogger {
	return s.logger()
}

func (s *distributedLockSvcImpl) Lock(ctx context.Context, ref string) (string, error) {
	l := s.l().C(ctx).Mth("lock").F(KV{"ref": ref}).Dbg()

	releaseId := NewRandString()
	err := <-Await(func() (bool, error) {
		locked, err := s.storage.Lock(ctx, ref, releaseId)
		if err != nil {
			return false, err
		}
		return locked, nil
	}, time.Millisecond*300, s.cfg.AwaitPeriod)

	if err != nil {
		return "", ErrLock(ctx, err, ref)
	}

	l.Dbg("locked")
	return releaseId, nil
}

func (s *distributedLockSvcImpl) UnLock(ctx context.Context, ref, releaseId string) {
	l := s.l().C(ctx).Mth("unlock").F(KV{"ref": ref, "releaseId": releaseId}).Dbg()
	if err := s.storage.UnLock(ctx, ref, releaseId); err != nil {
		l.E(err).St().Err()
	}
}
