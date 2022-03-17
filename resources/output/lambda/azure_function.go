package lambda

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output/local_exec"
)

const AzureResourceName = "azurerm_function_app"
const AzureAppServicePlanResourceName = "azurerm_app_service_plan"

type AzureFunctionApp struct {
	*common.AzResource      `hcl:",squash" default:"name=azurerm_function_app"`
	StorageAccountName      string               `hcl:"storage_account_name,expr"`
	StorageAccountAccessKey string               `hcl:"storage_account_access_key,expr"`
	AppServicePlanId        string               `hcl:"app_service_plan_id,expr"`
	OperatingSystem         string               `hcl:"os_type" hcle:"omitempty"`
	LocalExec               local_exec.LocalExec `hcl:"provisioner" hcle:"omitempty"`
	AppSettings             map[string]string    `hcl:"app_settings" hcle:"omitempty"`
}

type AzureAppServicePlan struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_app_service_plan"`
	Kind               string   `hcl:"kind"`
	Reserved           bool     `hcl:"reserved"`
	Sku                AzureSku `hcl:"sku"`
}

type AzureSku struct {
	Tier string `hcl:"tier"`
	Size string `hcl:"size"`
}
