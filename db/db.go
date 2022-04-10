package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/jsonpb"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/configpb"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Database struct {
	// TODO: store this in S3
	keyValueStore map[string]string
	marshaler     *jsonpb.Marshaler
	client        aws_client.Client
	sqlConnection *sql.DB
}

const (
	configFile = "config.pb.json"
	tfState    = "terraform.tfstate"
)

type LockType string

const (
	MainConfigLock       = "main"
	lockRetryPeriod      = 1 * time.Second
	lockRetryBackoffMult = 1.3
	lockExpirationPeriod = 2 * time.Hour
)

type ConfigLock struct {
	userId              string
	lockId              string
	expirationTimestamp time.Time

	active bool
}

func (l *ConfigLock) IsActive() bool {
	return l != nil && l.active && time.Now().UTC().Before(l.expirationTimestamp)
}

type lockErr struct {
	retryable bool
	error
}

func (d *Database) LockConfig(ctx context.Context, userId string) (*ConfigLock, error) {
	retryPeriod := lockRetryPeriod
	for {
		configLock, err := d.lockConfig(ctx, userId)
		if err != nil {
			if !err.retryable {
				return nil, err.error
			} else {
				log.Println(err.Error())
				if configLock != nil {
					log.Printf("configLock is locked (until %s), waiting for %s and then trying again\n", configLock.expirationTimestamp, retryPeriod)
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
			return configLock, nil
		}
	}
}

func isRetryableSqlErr(err error) bool {
	return err != sql.ErrTxDone && err != sql.ErrConnDone
}

func (d *Database) lockConfig(ctx context.Context, userId string) (*ConfigLock, *lockErr) {
	log.Println("locking")
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
		QueryRowContext(ctx, "SELECT UserId, LockId, LockExpirationTimestamp FROM Locks WHERE UserId = ? AND LockId = ?;", userId, MainConfigLock).
		Scan(&row.userId, &row.lockId, &row.expirationTimestamp)
	now := time.Now().UTC()
	expirationTimestamp := now.Add(lockExpirationPeriod)
	if err == sql.ErrNoRows {
		_, err := d.sqlConnection.
			ExecContext(ctx, "INSERT INTO Locks (UserId, LockId, LockExpirationTimestamp) VALUES (?, ?, ?);", userId, MainConfigLock, expirationTimestamp)
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
	} else if err != nil {
		return nil, &lockErr{false, err}
	} else if err == nil && now.After(row.expirationTimestamp) {
		log.Println("ConfigLock has expired, overwriting it")
		_, err := d.sqlConnection.
			ExecContext(ctx, "UPDATE Locks SET LockExpirationTimestamp = ? WHERE UserId = ? AND LockId = ?;", expirationTimestamp, userId, MainConfigLock)
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

func (d *Database) UnlockConfig(_ context.Context, lock *ConfigLock) error {
	log.Println("unlocking")
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

func (d *Database) StoreUserConfig(config *configpb.Config, lock *ConfigLock) error {
	if !lock.IsActive() {
		return fmt.Errorf("unable to store user config because lock is invalid")
	}
	log.Printf("Storing user config from api_key %s\n", config.UserId)
	str, err := d.marshaler.MarshalToString(config)
	if err != nil {
		return errors.InternalServerErrorWithMessage("unable to marshal configuration", err)
	}

	err = d.client.SaveFile(config.UserId, configFile, str)
	if err != nil {
		return errors.InternalServerErrorWithMessage("error storing configuration", err)
	}
	tmpDir := filepath.Join(os.TempDir(), "multy", config.UserId)
	data, err := os.ReadFile(filepath.Join(tmpDir, tfState))
	if err != nil {
		return errors.InternalServerErrorWithMessage("error reading current infra state cache", err)
	}

	err = d.client.SaveFile(config.UserId, tfState, string(data))
	if err != nil {
		return errors.InternalServerErrorWithMessage("error storing current infra state", err)
	}

	return nil
}

func (d *Database) LoadUserConfig(userId string, lock *ConfigLock) (*configpb.Config, error) {
	if lock != nil && !lock.IsActive() {
		return nil, fmt.Errorf("unable to load user config because lock is invalid")
	}
	log.Printf("Loading config from api_key %s\n", userId)
	result := configpb.Config{
		UserId: userId,
	}

	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	if lock == nil {
		tmpDir = filepath.Join(os.TempDir(), "multy/readonly", userId)
	}
	err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error creating output file", err)
	}

	tfFileStr := ""
	tfFileStr, err = d.client.ReadFile(userId, configFile)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error read configuration", err)
	}
	tfStateStr, err := d.client.ReadFile(userId, tfState)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error reading current infra state", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, tfState), []byte(tfStateStr), os.ModePerm&0664)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error caching current infra state", err)
	}

	if tfFileStr != "" {
		err := jsonpb.UnmarshalString(tfFileStr, &result)
		if err != nil {
			return nil, errors.InternalServerErrorWithMessage("error parsing configuration", err)
		}
	}
	return &result, nil
}

func (d *Database) Close() error {
	err := d.sqlConnection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetUserId(ctx context.Context, apiKey string) (string, error) {
	var userId string
	err := d.sqlConnection.QueryRowContext(ctx, "SELECT UserId FROM ApiKeys WHERE ApiKey = ?;", apiKey).Scan(&userId)
	if err == sql.ErrNoRows {
		return "", errors.PermissionDenied("wrong api key")
	} else if err != nil {
		return "", err
	}

	return userId, err
}

func NewDatabase(connectionString string) (*Database, error) {
	marshaler, err := proto.GetMarshaler()
	if err != nil {
		return nil, err
	}

	// "goncalo:@tcp(localhost:3306)/multydb?parseTime=true"

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return &Database{
		keyValueStore: map[string]string{},
		marshaler:     marshaler,
		client:        aws_client.Configure(),
		sqlConnection: db,
	}, nil
}
