package proto

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"io/ioutil"
	"path/filepath"
)

func getAnyResolver() (jsonpb.AnyResolver, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("error while reading current directory: %s", err.Error())
	}
	var descs []*desc.FileDescriptor
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".proto" {
			bytes, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return nil, err
			}
			descriptor, err := desc.LoadFileDescriptor(string(bytes))
			if err != nil {
				return nil, err
			}
			descs = append(descs, descriptor)
		}
	}

	return dynamic.AnyResolver(nil, descs...), nil
}

func GetMarshaler() (*jsonpb.Marshaler, error) {
	resolver, err := getAnyResolver()
	if err != nil {
		return nil, err
	}
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "\t",
		OrigName:     true,
		AnyResolver:  resolver,
	}
	return &marshaler, nil
}
