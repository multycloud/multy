package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"google.golang.org/protobuf/proto"
)

var Metadatas = map[proto.Message]resources.ResourceMetadataInterface{
	&resourcespb.VirtualNetworkArgs{}:        &virtualNetworkMetadata,
	&resourcespb.SubnetArgs{}:                &subnetMetadata,
	&resourcespb.DatabaseArgs{}:              &dbMetadata,
	&resourcespb.PublicIpArgs{}:              &publicIpMetadata,
	&resourcespb.RouteTableArgs{}:            &routeTableMetadata,
	&resourcespb.RouteTableAssociationArgs{}: &routeTableAssociationMetadata,
	&resourcespb.KubernetesNodePoolArgs{}:    &kubernetesNodePoolMetadata,
	&resourcespb.KubernetesClusterArgs{}:     &kubernetesClusterMetadata,
	&resourcespb.LambdaArgs{}:                &lambdaMetadata,
	&resourcespb.NetworkInterfaceArgs{}:      &networkInterfaceMetadata,
	&resourcespb.NetworkSecurityGroupArgs{}:  &networkSecurityGroupMetadata,
	&resourcespb.ObjectStorageArgs{}:         &objectStorageMetadata,
	&resourcespb.ObjectStorageObjectArgs{}:   &objectStorageObjectMetadata,
	&resourcespb.VaultArgs{}:                 &vaultMetadata,
	&resourcespb.VaultAccessPolicyArgs{}:     &vaultAccessPolicyMetadata,
	&resourcespb.VaultSecretArgs{}:           &vaultSecretMetadata,
	&resourcespb.VirtualMachineArgs{}:        &virtualMachineMetadata,
	&resourcespb.ResourceGroupArgs{}:         &resourceGroupMetadata,
}
