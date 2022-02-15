package terraform

import "fmt"

const DataResourceName = "azurerm_client_config"

type ClientConfig struct {
	ResourceName string `hcl:",key"`
	ResourceId   string `hcl:",key"`
}

func (r *ClientConfig) GetResult() string {
	return fmt.Sprintf("azurerm_client_config.%s.result", r.ResourceId)
}
