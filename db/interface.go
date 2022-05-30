package db

import (
	"context"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/flags"
)

type LockDatabase interface {
	LockConfig(ctx context.Context, userId string) (lock *ConfigLock, err error)
	UnlockConfig(ctx context.Context, lock *ConfigLock) error
}

type Database interface {
	LockDatabase
	GetUserId(ctx context.Context, apiKey string) (string, error)
	CreateUser(ctx context.Context, emailAddress string) (apiKey string, err error)
	StoreUserConfig(config *configpb.Config, lock *ConfigLock) error
	LoadUserConfig(userId string, lock *ConfigLock) (*configpb.Config, error)
	Close() error
}

func NewDatabase(awsClient aws_client.AwsClient) (Database, error) {
	if flags.Environment == flags.Local {
		return newLocalDatabase(awsClient)
	}

	return newDatabase(awsClient)
}
