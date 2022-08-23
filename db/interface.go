package db

import (
	"context"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/flags"
)

const TfState = "terraform.tfstate"
const TfPlan = "terraform.tfplan"

type LockDatabase interface {
	LockConfig(ctx context.Context, userId string, lockId string) (lock *ConfigLock, err error)
	UnlockConfig(ctx context.Context, lock *ConfigLock) error
}

type TfStateReader interface {
	LoadTerraformState(ctx context.Context, userId string) (string, error)
	LoadTerraformPlan(ctx context.Context, configPrefix string) (string, error)
	StoreTerraformPlan(ctx context.Context, configPrefix string, plan string) error
}

type Database interface {
	TfStateReader
	LockDatabase
	GetUserId(ctx context.Context, apiKey string) (string, error)
	CreateUser(ctx context.Context, emailAddress string) (apiKey string, err error)
	StoreUserConfig(ctx context.Context, config *configpb.Config, configPrefix string, lock *ConfigLock) error
	LoadUserConfig(ctx context.Context, userId string, configPrefix string, lock *ConfigLock) (*configpb.Config, error)

	Close() error
}

func NewDatabase(awsClient aws_client.AwsClient) (Database, error) {
	if flags.Environment == flags.Local {
		return newLocalDatabase(awsClient)
	}

	return newDatabase(awsClient)
}
