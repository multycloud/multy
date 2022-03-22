package db

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/config"
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
	configFile = "config.pb.json"
	tfState    = "terraform.tfstate"
)

func (d *Database) StoreUserConfig(config *config.Config) error {
	fmt.Printf("Storing user config from api_key %s\n", config.UserId)
	str, err := d.marshaler.MarshalToString(config)
	if err != nil {
		return err
	}

	err = d.client.SaveFile(config.UserId, configFile, str)
	if err != nil {
		return err
	}

	tmpDir := filepath.Join(os.TempDir(), "multy", config.UserId)
	data, err := os.ReadFile(filepath.Join(tmpDir, tfState))
	if err != nil {
		return err
	}

	err = d.client.SaveFile(config.UserId, tfState, string(data))
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) LoadUserConfig(userId string) (*config.Config, error) {
	fmt.Printf("Loading config from api_key %s\n", userId)
	result := config.Config{
		UserId: userId,
	}

	//str, exists := d.keyValueStore[userId]
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return nil, fmt.Errorf("error creating output file: %s", err.Error())
	}

	tfFileStr, er := d.client.ReadFile(userId, configFile)
	if er != nil {
		return nil, *er
	}
	if tfFileStr != "" {
		er := jsonpb.UnmarshalString(tfFileStr, &result)
		if er != nil {
			return nil, er
		}
	}

	tfStateStr, er := d.client.ReadFile(userId, tfState)
	if er != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(tmpDir, tfState), []byte(tfStateStr), os.ModePerm&0664)
	if er != nil {
		return nil, err
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
