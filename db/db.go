package db

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/config"
)

type Database struct {
	// TODO: store this in S3
	keyValueStore map[string]string
	marshaler     *jsonpb.Marshaler
}

func (d *Database) Store(config *config.Config) error {
	str, err := d.marshaler.MarshalToString(config)
	if err != nil {
		return err
	}
	d.keyValueStore[config.UserId] = str
	return nil
}

func (d *Database) Load(userId string) (*config.Config, error) {
	result := config.Config{
		UserId: userId,
	}
	str, exists := d.keyValueStore[userId]
	if exists {

		err := jsonpb.UnmarshalString(str, &result)
		if err != nil {
			return nil, err
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
	}, nil
}
