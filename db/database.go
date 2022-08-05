package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"log"
	"os"
	"time"
)

type database struct {
	*userConfigStorage
	*RemoteLockDatabase
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
		log.Printf("[WARNING] api key %s not found\n", apiKey)
		return "", errors.PermissionDenied(fmt.Sprintf("Api key '%s' is invalid", apiKey))
	} else if err != nil {
		log.Printf("[ERROR] unable to retrieve user id from api key %s \n", apiKey)
		return "", err
	}

	return userId, err
}

func (d *database) CreateUser(ctx context.Context, emailAddress string) (apiKey string, err error) {
	tx, err := d.sqlConnection.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var userId string
	err = tx.QueryRowContext(ctx, "SELECT * FROM Users WHERE UserId = ?;", emailAddress).Scan(&userId)
	if err == nil {
		return "", errors.UserAlreadyExists(emailAddress)
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO Users (UserId) VALUES (?);", emailAddress)
	if err != nil {
		return "", err
	}
	apiKey = uuid.New().String()
	_, err = tx.ExecContext(ctx, "INSERT INTO ApiKeys (ApiKey, UserId) VALUES (?, ?);", apiKey, emailAddress)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return apiKey, err
}

func (d *database) LoadTerraformState(ctx context.Context, userId string) (string, error) {
	result, err := d.AwsClient.ReadFile(userId, TfState)
	if err != nil {
		return "", errors.InternalServerErrorWithMessage("error reading terraform state", err)
	}

	return result, nil
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
	lockDb, err := NewLockDatabase(dbConnection)
	if err != nil {
		return nil, err
	}
	return &database{
		userConfigStorage:  userStg,
		RemoteLockDatabase: lockDb,
		sqlConnection:      dbConnection,
		AwsClient:          awsClient,
	}, nil
}
