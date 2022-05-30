package db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	aws_client "github.com/multycloud/multy/api/aws"
)

type localDatabase struct {
	*userConfigStorage
	lockDatabase *LocalLockDatabase
}

func (d *localDatabase) Close() error {
	return nil
}

func (d *localDatabase) GetUserId(ctx context.Context, apiKey string) (string, error) {
	return apiKey, nil
}

func (d *localDatabase) CreateUser(ctx context.Context, emailAddress string) (string, error) {
	return emailAddress, nil
}

func (d *localDatabase) LockConfig(ctx context.Context, userId string) (lock *ConfigLock, err error) {
	return d.lockDatabase.LockConfig(ctx, userId)
}

func (d *localDatabase) UnlockConfig(ctx context.Context, lock *ConfigLock) error {
	return d.lockDatabase.UnlockConfig(ctx, lock)
}

func newLocalDatabase(awsClient aws_client.AwsClient) (*localDatabase, error) {
	userStg, err := newUserConfigStorage(awsClient)
	if err != nil {
		return nil, err
	}
	return &localDatabase{
		userConfigStorage: userStg,
		lockDatabase:      newLocalLockDatabase(),
	}, nil
}
