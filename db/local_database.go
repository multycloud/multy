package db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	aws_client "github.com/multycloud/multy/api/aws"
	"os"
	"path"
	"path/filepath"
)

type localDatabase struct {
	*userConfigStorage
	*LocalLockDatabase
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

func (d *localDatabase) LoadTerraformState(_ context.Context, userId string) (string, error) {
	file, err := os.ReadFile(path.Join(filepath.Join(os.TempDir(), "multy", userId, "local"), TfState))
	// empty state is fine and expected in dry runs
	if os.IsNotExist(err) {
		return "{}", nil
	}
	return string(file), err
}

func (d *localDatabase) LoadTerraformPlan(_ context.Context, userId string) (string, error) {
	file, err := os.ReadFile(path.Join(filepath.Join(os.TempDir(), "multy", userId, "local"), TfPlan))
	// empty plan is fine and expected in dry runs
	if os.IsNotExist(err) {
		return "", nil
	}
	return string(file), err
}

func (d *localDatabase) StoreTerraformPlan(_ context.Context, userId string, plan string) error {
	err := os.WriteFile(path.Join(filepath.Join(os.TempDir(), "multy", userId, "local"), TfPlan), []byte(plan), os.ModePerm&0664)
	return err
}

func newLocalDatabase(awsClient aws_client.AwsClient) (*localDatabase, error) {
	userStg, err := newUserConfigStorage(awsClient)
	if err != nil {
		return nil, err
	}
	return &localDatabase{
		userConfigStorage: userStg,
		LocalLockDatabase: newLocalLockDatabase(),
	}, nil
}
