package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"os"
	"time"
)

type database struct {
	*userConfigStorage
	lockDatabase  *RemoteLockDatabase
	sqlConnection *sql.DB
	AwsClient     aws_client.AwsClient
}

const (
	configFile = "config.pb.json"
)

type LockType string

const (
	MainConfigLock       = "main"
	lockRetryPeriod      = 1 * time.Second
	lockRetryBackoffMult = 1.3
	lockExpirationPeriod = 2 * time.Hour
)

func (d *database) Close() error {
	err := d.sqlConnection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *database) GetUserId(ctx context.Context, apiKey string) (string, error) {
	var userId string
	err := d.sqlConnection.QueryRowContext(ctx, "SELECT UserId FROM ApiKeys WHERE ApiKey = ?;", apiKey).Scan(&userId)
	if err == sql.ErrNoRows {
		return "", errors.PermissionDenied(fmt.Sprintf("Api key '%s' is invalid", apiKey))
	} else if err != nil {
		return "", err
	}

	return userId, err
}

func (d *database) LockConfig(ctx context.Context, userId string) (lock *ConfigLock, err error) {
	return d.lockDatabase.LockConfig(ctx, userId)
}

func (d *database) UnlockConfig(ctx context.Context, lock *ConfigLock) error {
	return d.lockDatabase.UnlockConfig(ctx, lock)
}

func newDatabase(awsClient aws_client.AwsClient) (*database, error) {
	connectionString, exists := os.LookupEnv("MULTY_DB_CONN_STRING")
	if !exists {
		return nil, fmt.Errorf("db_connection_string env var is not set")
	}
	userStg, err := newUserConfigStorage(awsClient)
	if err != nil {
		return nil, err
	}

	dbConnection, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	db, err := NewLockDatabase(dbConnection)
	if err != nil {
		return nil, err
	}
	return &database{
		userConfigStorage: userStg,
		lockDatabase:      db,
		sqlConnection:     dbConnection,
		AwsClient:         awsClient,
	}, nil
}
