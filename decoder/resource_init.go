package decoder

import (
	"fmt"
	"multy-go/resources"
	rg "multy-go/resources/resource_group"
	"multy-go/resources/types"
)

func InitResource(resourceType string, commonParams *resources.CommonResourceParams) (resources.Resource, error) {
	switch resourceType {
	case "virtual_network":
		return &types.VirtualNetwork{CommonResourceParams: commonParams}, nil
	case "subnet":
		return &types.Subnet{CommonResourceParams: commonParams}, nil
	case "network_security_group":
		return &types.NetworkSecurityGroup{CommonResourceParams: commonParams}, nil
	case "virtual_machine":
		return &types.VirtualMachine{CommonResourceParams: commonParams}, nil
	case "public_ip":
		return &types.PublicIp{CommonResourceParams: commonParams}, nil
	case "route_table":
		return &types.RouteTable{CommonResourceParams: commonParams}, nil
	case "route_table_association":
		return &types.RouteTableAssociation{CommonResourceParams: commonParams}, nil
	case "network_interface":
		return &types.NetworkInterface{CommonResourceParams: commonParams}, nil
	case "database":
		return &types.Database{CommonResourceParams: commonParams}, nil
	case "resource_group":
		return &rg.Type{ResourceId: commonParams.ResourceId}, nil
	case "object_storage":
		return &types.ObjectStorage{CommonResourceParams: commonParams}, nil
	case "object_storage_object":
		return &types.ObjectStorageObject{CommonResourceParams: commonParams}, nil
	case "vault":
		return &types.Vault{CommonResourceParams: commonParams}, nil
	case "vault_secret":
		return &types.VaultSecret{CommonResourceParams: commonParams}, nil
	case "vault_access_policy":
		return &types.VaultAccessPolicy{CommonResourceParams: commonParams}, nil
	case "lambda":
		return &types.Lambda{CommonResourceParams: commonParams}, nil
	case "kubernetes_service":
		return &types.KubernetesService{CommonResourceParams: commonParams}, nil
	case "kubernetes_node_pool":
		return &types.KubernetesServiceNodePool{CommonResourceParams: commonParams}, nil
	default:
		return nil, fmt.Errorf("unknown resource type %s", resourceType)
	}
}
