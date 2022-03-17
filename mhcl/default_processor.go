package mhcl

import (
	"github.com/multycloud/multy/validate"
	"reflect"
	"strings"
)

type DefaultTagProcessor struct {
}

type NameSetter interface {
	SetName(string)
}

func (p *DefaultTagProcessor) Process(r any) {
	t := reflect.TypeOf(r)
	tValue := reflect.ValueOf(r)

	if t.Kind() == reflect.Ptr {
		t = reflect.TypeOf(r).Elem()
		tValue = reflect.ValueOf(r).Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		if tagValue, ok := t.Field(i).Tag.Lookup("default"); ok {
			values := strings.Split(tagValue, ",")
			for _, v := range values {
				keyVal := strings.SplitN(v, "=", 2)
				key := keyVal[0]
				defaultVal := keyVal[1]
				switch key {
				case "name":
					tValue.Field(i).Interface().(NameSetter).SetName(defaultVal)
				default:
					validate.LogInternalError("unknown key '%s' in tag %s", key, tagValue)
				}
			}
		}
	}
}
