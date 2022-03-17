package public_ip

import (
	"fmt"
	"multy/resources/common"
	"multy/validate"
)

const AzureResourceName = "azurerm_public_ip"

type AzurePublicIp struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_public_ip"`
	AllocationMethod   string `hcl:"allocation_method"`
}

func (pIp AzurePublicIp) GetId(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("%s.%s.id", AzureResourceName, pIp.ResourceId)
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}
