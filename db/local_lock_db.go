package db

import (
	"context"
	"sync"
	"time"
)

type LocalLockDatabase struct {
	lockCache sync.Map
}

func (d *LocalLockDatabase) LockConfig(_ context.Context, userId string) (lock *ConfigLock, err error) {
	l, _ := d.lockCache.LoadOrStore(userId, &sync.Mutex{})
	l.(*sync.Mutex).Lock()
	return &ConfigLock{
		userId:              userId,
		expirationTimestamp: time.Now().Add(24 * time.Hour),
		active:              true,
	}, nil
}

func (d *LocalLockDatabase) UnlockConfig(_ context.Context, lock *ConfigLock) error {
	if !lock.IsActive() {
		return nil
	}
	l, _ := d.lockCache.Load(lock.userId)
	l.(*sync.Mutex).Unlock()
	lock.active = false
	return nil
}

func newLocalLockDatabase() *LocalLockDatabase {
	return &LocalLockDatabase{
		sync.Map{},
	}
}
