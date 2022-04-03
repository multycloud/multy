package db

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/configpb"
	"os"
	"path/filepath"
)

type Database struct {
	// TODO: store this in S3
	keyValueStore map[string]string
	marshaler     *jsonpb.Marshaler
	client        aws_client.Client
}

const (
	configFile     = "config.pb.json"
	tfState        = "terraform.tfstate"
	multyLocalUser = "multy_local"
)

func (d *Database) StoreUserConfig(config *configpb.Config) error {
	fmt.Printf("Storing user config from api_key %s\n", config.UserId)
	str, err := d.marshaler.MarshalToString(config)
	if err != nil {
		return errors.InternalServerErrorWithMessage("unable to marshal configuration", err)
	}

	if config.UserId != multyLocalUser {
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
	} else {
		d.keyValueStore[config.UserId] = str
	}

	return nil
}

func (d *Database) LoadUserConfig(userId string) (*configpb.Config, error) {
	fmt.Printf("Loading config from api_key %s\n", userId)
	result := configpb.Config{
		UserId: userId,
	}

	//str, exists := d.keyValueStore[userId]
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error creating output file", err)
	}

	tfFileStr := ""
	if userId != multyLocalUser {
		var errPtr *error
		tfFileStr, errPtr = d.client.ReadFile(userId, configFile)
		if errPtr != nil {
			return nil, errors.InternalServerErrorWithMessage("error read configuration", *errPtr)
		}
		tfStateStr, errPtr := d.client.ReadFile(userId, tfState)
		if errPtr != nil {
			return nil, errors.InternalServerErrorWithMessage("error reading current infra state", *errPtr)
		}
		err = os.WriteFile(filepath.Join(tmpDir, tfState), []byte(tfStateStr), os.ModePerm&0664)
		if err != nil {
			return nil, errors.InternalServerErrorWithMessage("error caching current infra state", err)
		}
	} else {
		tfFileStr = d.keyValueStore[userId]
	}

	if tfFileStr != "" {
		err := jsonpb.UnmarshalString(tfFileStr, &result)
		if err != nil {
			return nil, errors.InternalServerErrorWithMessage("error parsing configuration", err)
		}
	}
	return &result, nil
}

func NewDatabase() (*Database, error) {
	marshaler, err := proto.GetMarshaler()
	if err != nil {
		return nil, err
	}
	return &Database{
		keyValueStore: map[string]string{},
		marshaler:     marshaler,
		client:        aws_client.Configure(),
	}, nil
}
