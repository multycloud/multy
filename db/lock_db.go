package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/trace"
	"time"
)

type ConfigLock struct {
	userId              string
	lockId              string
	expirationTimestamp time.Time

	active bool
	// used to cache lock so we don't do unnecessary db calls when running in a single server
	localLock   *ConfigLock
	traceRegion *trace.Region
}

type RemoteLockDatabase struct {
	sqlConnection     *sql.DB
	localLockDatabase *LocalLockDatabase
}

func (l *ConfigLock) IsActive() bool {
	return l != nil && l.active && time.Now().UTC().Before(l.expirationTimestamp)
}

type lockErr struct {
	retryable bool
	error
}

func (d *RemoteLockDatabase) LockConfig(ctx context.Context, userId string, lockId string) (lock *ConfigLock, err error) {
	localLock, err := d.localLockDatabase.LockConfig(ctx, userId, lockId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			// something went wrong, which means unlocking is not the caller's responsability
			d.localLockDatabase.UnlockConfig(ctx, localLock)
		}
	}()
	retryPeriod := lockRetryPeriod
	for {
		configLock, err := d.lockConfig(ctx, userId, lockId)
		if err != nil {
			if !err.retryable {
				return nil, err.error
			} else {
				log.Println(err.Error())
				if configLock != nil {
					log.Printf("[INFO] configLock is locked (until %s), waiting for %s and then trying again\n", configLock.expirationTimestamp, retryPeriod)
				}
				select {
				case <-time.After(retryPeriod):
					newRetryPeriod := float64(retryPeriod) * lockRetryBackoffMult
					retryPeriod = time.Duration(newRetryPeriod)
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
		} else {
			configLock.localLock = localLock
			return configLock, nil
		}
	}
}

func isRetryableSqlErr(err error) bool {
	return err != sql.ErrTxDone && err != sql.ErrConnDone
}

func (d *RemoteLockDatabase) lockConfig(ctx context.Context, userId string, lockId string) (*ConfigLock, *lockErr) {
	log.Println("[DEBUG] locking")
	tx, err := d.sqlConnection.BeginTx(ctx, nil)
	if err != nil {
		return nil, &lockErr{false, err}
	}
	committed := false
	defer func() {
		if !committed {
			err := tx.Rollback()
			if err != nil {
				log.Printf("[ERROR] %s\n", err.Error())
			}
		}
	}()

	row := ConfigLock{
		active: true,
	}
	err = d.sqlConnection.
		QueryRowContext(ctx, "SELECT UserId, LockId, LockExpirationTimestamp FROM Locks WHERE UserId = ? AND LockId = ?;", userId, lockId).
		Scan(&row.userId, &row.lockId, &row.expirationTimestamp)
	now := time.Now().UTC()
	expirationTimestamp := now.Add(lockExpirationPeriod)
	if err == sql.ErrNoRows {
		_, err := d.sqlConnection.
			ExecContext(ctx, "INSERT INTO Locks (UserId, LockId, LockExpirationTimestamp) VALUES (?, ?, ?);", userId, lockId, expirationTimestamp)
		if err != nil {
			return nil, &lockErr{isRetryableSqlErr(err), err}
		}
		err = tx.Commit()
		if err != nil {
			return nil, &lockErr{isRetryableSqlErr(err), err}
		}
		committed = true
		row.userId = userId
		row.lockId = lockId
		row.expirationTimestamp = expirationTimestamp
		return &row, nil
	} else if err != nil {
		return nil, &lockErr{false, err}
	} else if err == nil && now.After(row.expirationTimestamp) {
		log.Println("[WARNING] ConfigLock has expired, overwriting it")
		_, err := d.sqlConnection.
			ExecContext(ctx, "UPDATE Locks SET LockExpirationTimestamp = ? WHERE UserId = ? AND LockId = ?;", expirationTimestamp, userId, lockId)
		if err != nil {
			return nil, &lockErr{isRetryableSqlErr(err), err}
		}
		err = tx.Commit()
		if err != nil {
			return nil, &lockErr{isRetryableSqlErr(err), err}
		}
		committed = true
		row.userId = userId
		row.lockId = MainConfigLock
		row.expirationTimestamp = expirationTimestamp
		return &row, nil
	} else {
		return &row, &lockErr{true, fmt.Errorf("config lock is already taken")}
	}
}

func (d *RemoteLockDatabase) UnlockConfig(ctx context.Context, lock *ConfigLock) error {
	log.Println("[DEBUG] unlocking")
	defer func() {
		d.localLockDatabase.UnlockConfig(ctx, lock.localLock)
	}()
	if !lock.IsActive() {
		return nil
	}
	_, err := d.sqlConnection.Exec("DELETE FROM Locks WHERE UserId = ? AND LockId = ?;", lock.userId, lock.lockId)
	if err != nil {
		log.Printf("[ERROR] error unlocking, %s\n", err.Error())
		return err
	}
	lock.active = false
	return nil
}

func NewLockDatabase(dbConnection *sql.DB) (*RemoteLockDatabase, error) {
	return &RemoteLockDatabase{
		sqlConnection:     dbConnection,
		localLockDatabase: &LocalLockDatabase{},
	}, nil
}
