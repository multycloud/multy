package output

import (
	"fmt"
)

type TfResource struct {
	Address string                 `json:"address"`
	Values  map[string]interface{} `json:"values"`
}

type TfState struct {
	Values struct {
		RootModule struct {
			Resources []TfResource `json:"resources"`
		} `json:"root_module"`
	} `json:"values"`
}

func (t *TfState) GetValues(resourceType any, resourceId string) (map[string]interface{}, error) {
	address := fmt.Sprintf("%s.%s", GetResourceName(resourceType), resourceId)
	for _, r := range t.Values.RootModule.Resources {
		if r.Address == address {
			return r.Values, nil
		}
	}

	return nil, fmt.Errorf("resource %s doesn't exist", address)
}
