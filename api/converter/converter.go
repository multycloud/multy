package converter

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	common_resources "github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ResourceConverters[Arg proto.Message, OutT proto.Message] interface {
	Convert(resourceId string, request Arg, state *output.TfState) (OutT, error)
}

type MultyResourceConverter interface {
	ConvertToMultyResource(resourceId string, arg proto.Message, resources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error)
	GetResourceType() string
}

type ResourceInitFunc[T proto.Message, O resources.Resource] func(string, T, resources.Resources) (O, error)

type ResourceMetadata struct {
	InitFunc        ResourceInitFunc[proto.Message, resources.Resource]
	AbbreviatedName string
}

var Converters = map[proto.Message]ResourceMetadata{
	&resourcespb.VirtualNetworkArgs{}:        {cast(types.NewVirtualNetwork), "vn"},
	&resourcespb.SubnetArgs{}:                {cast(types.NewSubnet), "vn"},
	&resourcespb.DatabaseArgs{}:              {cast(types.NewDatabase), "db"},
	&resourcespb.PublicIpArgs{}:              {cast(types.NewPublicIp), "pip"},
	&resourcespb.RouteTableArgs{}:            {cast(types.NewRouteTable), "vn"},
	&resourcespb.RouteTableAssociationArgs{}: {InitFunc: cast(types.NewRouteTableAssociation)},
	&resourcespb.KubernetesNodePoolArgs{}:    {cast(types.NewKubernetesNodePool), "ks"},
	&resourcespb.KubernetesClusterArgs{}:     {cast(types.NewKubernetesCluster), "ks"},
	&resourcespb.LambdaArgs{}:                {cast(types.NewLambda), "fun"},
	&resourcespb.NetworkInterfaceArgs{}:      {cast(types.NewNetworkInterface), "nic"},
	&resourcespb.NetworkSecurityGroupArgs{}:  {cast(types.NewNetworkSecurityGroup), "nsg"},
	&resourcespb.ObjectStorageArgs{}:         {cast(types.NewObjectStorage), "st"},
	&resourcespb.ObjectStorageObjectArgs{}:   {cast(types.NewObjectStorageObject), "st"},
	&resourcespb.VaultArgs{}:                 {cast(types.NewVault), "kv"},
	&resourcespb.VaultAccessPolicyArgs{}:     {cast(types.NewVaultAccessPolicy), "kv"},
	&resourcespb.VaultSecretArgs{}:           {cast(types.NewVaultSecret), "kv"},
	&resourcespb.VirtualMachineArgs{}:        {cast(types.NewVirtualMachine), "vm"},
}

func GetConverter(name protoreflect.FullName) (*ResourceMetadata, error) {
	for messageType, conv := range Converters {
		if name == proto.MessageName(messageType) {
			return &conv, nil
		}
	}
	return nil, fmt.Errorf("unknown resource type %s", name)
}

func cast[T proto.Message, O resources.Resource](f func(string, T, resources.Resources) (O, error)) ResourceInitFunc[proto.Message, resources.Resource] {
	return func(resourceId string, arg proto.Message, r resources.Resources) (resources.Resource, error) {
		return f(resourceId, arg.(T), r)
	}
}
