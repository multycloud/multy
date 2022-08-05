package db

import (
	"context"
	"runtime/trace"
	"sync"
	"time"
)

type LocalLockDatabase struct {
	lockCache sync.Map
}

func (d *LocalLockDatabase) LockConfig(ctx context.Context, userId string, lockId string) (lock *ConfigLock, err error) {
	region := trace.StartRegion(ctx, "lock wait")
	defer region.End()
	l, _ := d.lockCache.LoadOrStore(lockId, &sync.Mutex{})
	l.(*sync.Mutex).Lock()
	return &ConfigLock{
		userId:              userId,
		lockId:              lockId,
		expirationTimestamp: time.Now().Add(24 * time.Hour),
		active:              true,
	}, nil
}

func (d *LocalLockDatabase) UnlockConfig(_ context.Context, lock *ConfigLock) error {
	if !lock.IsActive() {
		return nil
	}
	l, _ := d.lockCache.Load(lock.lockId)
	l.(*sync.Mutex).Unlock()
	lock.active = false
	return nil
}

func newLocalLockDatabase() *LocalLockDatabase {
	return &LocalLockDatabase{
		sync.Map{},
	}
}
