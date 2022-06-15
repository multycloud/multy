package metadata

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/types"
	aws_resources "github.com/multycloud/multy/resources/types/aws"
	azure_resources "github.com/multycloud/multy/resources/types/azure"
	gcp_resources "github.com/multycloud/multy/resources/types/gcp"
	"google.golang.org/protobuf/proto"
)

var Metadatas = map[proto.Message]resources.ResourceMetadataInterface{
	&resourcespb.VirtualNetworkArgs{}: &resources.ResourceMetadata[*resourcespb.VirtualNetworkArgs, *types.VirtualNetwork, *resourcespb.VirtualNetworkResource]{
		AbbreviatedName: "vn",
		ResourceType:    "virtual_network",
		Translators: map[commonpb.CloudProvider]func(*types.VirtualNetwork) resources.ResourceTranslator[*resourcespb.VirtualNetworkResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitVirtualNetwork,
			commonpb.CloudProvider_AZURE: azure_resources.InitVirtualNetwork,
			commonpb.CloudProvider_GCP:   gcp_resources.InitVirtualNetwork,
		},
	},
	&resourcespb.SubnetArgs{}: &resources.ResourceMetadata[*resourcespb.SubnetArgs, *types.Subnet, *resourcespb.SubnetResource]{
		AbbreviatedName: "vn",
		ResourceType:    "subnet",
		Translators: map[commonpb.CloudProvider]func(*types.Subnet) resources.ResourceTranslator[*resourcespb.SubnetResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitSubnet,
			commonpb.CloudProvider_AZURE: azure_resources.InitSubnet,
		},
	},
	&resourcespb.DatabaseArgs{}: &resources.ResourceMetadata[*resourcespb.DatabaseArgs, *types.Database, *resourcespb.DatabaseResource]{
		AbbreviatedName: "db",
		ResourceType:    "database",
		Translators: map[commonpb.CloudProvider]func(*types.Database) resources.ResourceTranslator[*resourcespb.DatabaseResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitDatabase,
			commonpb.CloudProvider_AZURE: azure_resources.InitDatabase,
		},
	},
	&resourcespb.PublicIpArgs{}: &resources.ResourceMetadata[*resourcespb.PublicIpArgs, *types.PublicIp, *resourcespb.PublicIpResource]{
		AbbreviatedName: "pip",
		ResourceType:    "public_ip",
		Translators: map[commonpb.CloudProvider]func(*types.PublicIp) resources.ResourceTranslator[*resourcespb.PublicIpResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitPublicIp,
			commonpb.CloudProvider_AZURE: azure_resources.InitPublicIp,
		},
	},
	&resourcespb.RouteTableArgs{}: &resources.ResourceMetadata[*resourcespb.RouteTableArgs, *types.RouteTable, *resourcespb.RouteTableResource]{
		AbbreviatedName: "rt",
		ResourceType:    "route_table",
		Translators: map[commonpb.CloudProvider]func(*types.RouteTable) resources.ResourceTranslator[*resourcespb.RouteTableResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitRouteTable,
			commonpb.CloudProvider_AZURE: azure_resources.InitRouteTable,
		},
	},
	&resourcespb.RouteTableAssociationArgs{}: &resources.ResourceMetadata[*resourcespb.RouteTableAssociationArgs, *types.RouteTableAssociation, *resourcespb.RouteTableAssociationResource]{
		AbbreviatedName: "rt",
		ResourceType:    "route_table_association",
		Translators: map[commonpb.CloudProvider]func(*types.RouteTableAssociation) resources.ResourceTranslator[*resourcespb.RouteTableAssociationResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitRouteTableAssociation,
			commonpb.CloudProvider_AZURE: azure_resources.InitRouteTableAssociation,
		},
	},
	&resourcespb.KubernetesNodePoolArgs{}: &resources.ResourceMetadata[*resourcespb.KubernetesNodePoolArgs, *types.KubernetesNodePool, *resourcespb.KubernetesNodePoolResource]{
		AbbreviatedName: "ks",
		ResourceType:    "kubernetes_node_pool",
		Translators: map[commonpb.CloudProvider]func(*types.KubernetesNodePool) resources.ResourceTranslator[*resourcespb.KubernetesNodePoolResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitKubernetesNodePool,
			commonpb.CloudProvider_AZURE: azure_resources.InitKubernetesNodePool,
		},
	},
	&resourcespb.KubernetesClusterArgs{}: &resources.ResourceMetadata[*resourcespb.KubernetesClusterArgs, *types.KubernetesCluster, *resourcespb.KubernetesClusterResource]{
		AbbreviatedName: "ks",
		ResourceType:    "kubernetes_cluster",
		Translators: map[commonpb.CloudProvider]func(*types.KubernetesCluster) resources.ResourceTranslator[*resourcespb.KubernetesClusterResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitKubernetesCluster,
			commonpb.CloudProvider_AZURE: azure_resources.InitKubernetesCluster,
		},
	},
	&resourcespb.NetworkInterfaceArgs{}: &resources.ResourceMetadata[*resourcespb.NetworkInterfaceArgs, *types.NetworkInterface, *resourcespb.NetworkInterfaceResource]{
		AbbreviatedName: "nic",
		ResourceType:    "network_interface",
		Translators: map[commonpb.CloudProvider]func(*types.NetworkInterface) resources.ResourceTranslator[*resourcespb.NetworkInterfaceResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitNetworkInterface,
			commonpb.CloudProvider_AZURE: azure_resources.InitNetworkInterface,
		},
	},
	&resourcespb.NetworkInterfaceSecurityGroupAssociationArgs{}: &resources.ResourceMetadata[*resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, *types.NetworkInterfaceSecurityGroupAssociation, *resourcespb.NetworkInterfaceSecurityGroupAssociationResource]{
		AbbreviatedName: "nic",
		ResourceType:    "network_interface_security_group_association",
		Translators: map[commonpb.CloudProvider]func(*types.NetworkInterfaceSecurityGroupAssociation) resources.ResourceTranslator[*resourcespb.NetworkInterfaceSecurityGroupAssociationResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitNetworkInterfaceSecurityGroupAssociation,
			commonpb.CloudProvider_AZURE: azure_resources.InitNetworkInterfaceSecurityGroupAssociation,
		},
	},
	&resourcespb.NetworkSecurityGroupArgs{}: &resources.ResourceMetadata[*resourcespb.NetworkSecurityGroupArgs, *types.NetworkSecurityGroup, *resourcespb.NetworkSecurityGroupResource]{
		AbbreviatedName: "nsg",
		ResourceType:    "network_security_group",
		Translators: map[commonpb.CloudProvider]func(*types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitNetworkSecurityGroup,
			commonpb.CloudProvider_AZURE: azure_resources.InitNetworkSecurityGroup,
		},
	},
	&resourcespb.ObjectStorageArgs{}: &resources.ResourceMetadata[*resourcespb.ObjectStorageArgs, *types.ObjectStorage, *resourcespb.ObjectStorageResource]{
		AbbreviatedName: "st",
		ResourceType:    "object_storage",
		Translators: map[commonpb.CloudProvider]func(*types.ObjectStorage) resources.ResourceTranslator[*resourcespb.ObjectStorageResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitObjectStorage,
			commonpb.CloudProvider_AZURE: azure_resources.InitObjectStorage,
		},
	},
	&resourcespb.ObjectStorageObjectArgs{}: &resources.ResourceMetadata[*resourcespb.ObjectStorageObjectArgs, *types.ObjectStorageObject, *resourcespb.ObjectStorageObjectResource]{
		AbbreviatedName: "st",
		ResourceType:    "object_storage_object",
		Translators: map[commonpb.CloudProvider]func(*types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitObjectStorageObject,
			commonpb.CloudProvider_AZURE: azure_resources.InitObjectStorageObject,
		},
	},
	&resourcespb.VaultArgs{}: &resources.ResourceMetadata[*resourcespb.VaultArgs, *types.Vault, *resourcespb.VaultResource]{
		AbbreviatedName: "kv",
		ResourceType:    "vault",
		Translators: map[commonpb.CloudProvider]func(*types.Vault) resources.ResourceTranslator[*resourcespb.VaultResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitVault,
			commonpb.CloudProvider_AZURE: azure_resources.InitVault,
		},
	},
	&resourcespb.VaultAccessPolicyArgs{}: &resources.ResourceMetadata[*resourcespb.VaultAccessPolicyArgs, *types.VaultAccessPolicy, *resourcespb.VaultAccessPolicyResource]{
		AbbreviatedName: "kv",
		ResourceType:    "vault_access_policy",
		Translators: map[commonpb.CloudProvider]func(*types.VaultAccessPolicy) resources.ResourceTranslator[*resourcespb.VaultAccessPolicyResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitVaultAccessPolicy,
			commonpb.CloudProvider_AZURE: azure_resources.InitVaultAccessPolicy,
		},
	},
	&resourcespb.VaultSecretArgs{}: &resources.ResourceMetadata[*resourcespb.VaultSecretArgs, *types.VaultSecret, *resourcespb.VaultSecretResource]{
		AbbreviatedName: "kv",
		ResourceType:    "vault_secret",
		Translators: map[commonpb.CloudProvider]func(secret *types.VaultSecret) resources.ResourceTranslator[*resourcespb.VaultSecretResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitVaultSecret,
			commonpb.CloudProvider_AZURE: azure_resources.InitVaultSecret,
		},
	},
	&resourcespb.VirtualMachineArgs{}: &resources.ResourceMetadata[*resourcespb.VirtualMachineArgs, *types.VirtualMachine, *resourcespb.VirtualMachineResource]{
		AbbreviatedName: "vm",
		ResourceType:    "virtual_machine",
		Translators: map[commonpb.CloudProvider]func(*types.VirtualMachine) resources.ResourceTranslator[*resourcespb.VirtualMachineResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitVirtualMachine,
			commonpb.CloudProvider_AZURE: azure_resources.InitVirtualMachine,
		},
	},
	&resourcespb.ResourceGroupArgs{}: &resources.ResourceMetadata[*resourcespb.ResourceGroupArgs, *types.ResourceGroup, *resourcespb.ResourceGroupResource]{
		AbbreviatedName: "rg",
		ResourceType:    "resource_group",
		Translators: map[commonpb.CloudProvider]func(*types.ResourceGroup) resources.ResourceTranslator[*resourcespb.ResourceGroupResource]{
			commonpb.CloudProvider_AWS:   aws_resources.InitResourceGroup,
			commonpb.CloudProvider_AZURE: azure_resources.InitResourceGroup,
			commonpb.CloudProvider_GCP:   gcp_resources.InitResourceGroup,
		},
	},
}
