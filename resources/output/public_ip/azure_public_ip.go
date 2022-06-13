package public_ip

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

const AzureResourceName = "azurerm_public_ip"

type AzurePublicIp struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_public_ip"`
	AllocationMethod   string `hcl:"allocation_method"`

	IpAddress string `json:"ip_address" hcle:"omitempty"`
}

func (pIp AzurePublicIp) GetId() string {

	return fmt.Sprintf("%s.%s.id", AzureResourceName, pIp.ResourceId)
}
