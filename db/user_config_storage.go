package db

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/configpb"
	"log"
)

type userConfigStorage struct {
	marshaler *jsonpb.Marshaler
	AwsClient aws_client.AwsClient
}

func (d *userConfigStorage) StoreUserConfig(config *configpb.Config, lock *ConfigLock) error {
	if !lock.IsActive() {
		return fmt.Errorf("unable to store user config because lock is invalid")
	}
	log.Printf("[INFO] Storing user config from api_key %s\n", config.UserId)
	str, err := d.marshaler.MarshalToString(config)
	if err != nil {
		return errors.InternalServerErrorWithMessage("unable to marshal configuration", err)
	}

	err = d.AwsClient.SaveFile(config.UserId, configFile, str)
	if err != nil {
		return errors.InternalServerErrorWithMessage("error storing configuration", err)
	}
	return nil
}

func (d *userConfigStorage) LoadUserConfig(userId string, lock *ConfigLock) (*configpb.Config, error) {
	if lock != nil && !lock.IsActive() {
		return nil, fmt.Errorf("unable to load user config because lock is invalid")
	}
	log.Printf("[INFO] Loading config from api_key %s\n", userId)
	result := configpb.Config{
		UserId: userId,
	}

	tfFileStr, err := d.AwsClient.ReadFile(userId, configFile)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error reading configuration", err)
	}

	if tfFileStr != "" {
		err := jsonpb.UnmarshalString(tfFileStr, &result)
		if err != nil {
			return nil, errors.InternalServerErrorWithMessage("error parsing configuration", err)
		}
	}
	return &result, nil
}

func newUserConfigStorage(awsClient aws_client.AwsClient) (*userConfigStorage, error) {
	marshaler, err := proto.GetMarshaler()
	if err != nil {
		return nil, err
	}

	return &userConfigStorage{
		marshaler: marshaler,
		AwsClient: awsClient,
	}, nil
}
