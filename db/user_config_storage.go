package db

import (
	"context"
	"fmt"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/configpb"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"runtime/trace"
)

type userConfigStorage struct {
	AwsClient aws_client.AwsClient
}

func (d *userConfigStorage) StoreUserConfig(ctx context.Context, config *configpb.Config, lock *ConfigLock) error {
	if !lock.IsActive() {
		return fmt.Errorf("unable to store user config because lock is invalid")
	}
	log.Printf("[INFO] Storing user config from api_key %s\n", config.UserId)
	region := trace.StartRegion(ctx, "config store")
	defer region.End()
	b, err := protojson.Marshal(config)
	if err != nil {
		return err
	}

	err = d.AwsClient.SaveFile(config.UserId, configFile, string(b))
	if err != nil {
		return errors.InternalServerErrorWithMessage("error storing configuration", err)
	}
	return nil
}

func (d *userConfigStorage) LoadUserConfig(ctx context.Context, userId string, lock *ConfigLock) (*configpb.Config, error) {
	if lock != nil && !lock.IsActive() {
		return nil, fmt.Errorf("unable to load user config because lock is invalid")
	}
	log.Printf("[INFO] Loading config from api_key %s\n", userId)
	region := trace.StartRegion(ctx, "config load")
	defer region.End()
	result := configpb.Config{
		UserId: userId,
	}

	configFileStr, err := d.AwsClient.ReadFile(userId, configFile)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error reading configuration", err)
	}

	if configFileStr != "" {
		err := protojson.Unmarshal([]byte(configFileStr), &result)
		if err != nil {
			return nil, errors.InternalServerErrorWithMessage("error parsing configuration", err)
		}
	}

	if result.UserId != userId {
		return nil, fmt.Errorf("config file is in unexpcted state. user should be %s but is %s", userId, result.UserId)
	}

	return &result, nil
}

func newUserConfigStorage(awsClient aws_client.AwsClient) (*userConfigStorage, error) {

	return &userConfigStorage{
		AwsClient: awsClient,
	}, nil
}
