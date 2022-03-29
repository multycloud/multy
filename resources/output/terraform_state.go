package output

import (
	"encoding/json"
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
	return t.Get(address)
}
func (t *TfState) Get(resourceRef string) (map[string]interface{}, error) {
	for _, r := range t.Values.RootModule.Resources {
		if r.Address == resourceRef {
			return r.Values, nil
		}
	}

	return nil, fmt.Errorf("resource %s doesn't exist", resourceRef)
}

func GetParsed[T any](state *TfState, resourceRef string) (*T, error) {
	rawResourceState, err := state.Get(resourceRef)
	if err != nil {
		return nil, err
	}

	jsonResourceState, err := json.Marshal(rawResourceState)
	if err != nil {
		return nil, err
	}

	stateResource := new(T)
	err = json.Unmarshal(jsonResourceState, stateResource)
	if err != nil {
		return nil, err
	}

	return stateResource, nil
}
