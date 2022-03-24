package converter

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
)

type ResourceConverters[Arg proto.Message, OutT proto.Message] interface {
	Convert(resourceId string, request []Arg, state *output.TfState) (OutT, error)
}

type MultyResourceConverter interface {
	ConvertToMultyResource(resourceId string, arg proto.Message, resources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error)
	GetResourceType() string
}

type VnConverter struct {
}
type SubnetConverter struct {
}
type NetworkInterfaceConverter struct {
}
type RouteTableConverter struct {
}
type NetworkSecurityGroupConverter struct {
}
type RouteTableAssociationConverter struct {
}
type DatabaseConverter struct {
}
type ObjectStorageConverter struct {
}
type ObjectStorageObjectConverter struct {
}
type PublicIpConverter struct {
}
type KubernetesClusterConverter struct {
}
type KubernetesNodePoolConverter struct {
}
type LambdaConverter struct {
}
type VaultConverter struct {
}
type VaultAccessPolicyConverter struct {
}
type VaultSecretConverter struct {
}
type VirtualMachineConverter struct {
}

func (v VnConverter) GetResourceType() string {
	return "virtual_network"
}

func (v SubnetConverter) GetResourceType() string {
	return "subnet"
}

func (v NetworkInterfaceConverter) GetResourceType() string {
	return "network_interface"
}
func (v RouteTableConverter) GetResourceType() string {
	return "route_table"
}
func (v NetworkSecurityGroupConverter) GetResourceType() string {
	return "network_security_group"
}
func (v RouteTableAssociationConverter) GetResourceType() string {
	return "route_table_association"
}
func (v DatabaseConverter) GetResourceType() string {
	return "database"
}
func (v ObjectStorageConverter) GetResourceType() string {
	return "object_storage"
}
func (v ObjectStorageObjectConverter) GetResourceType() string {
	return "object_storage_object"
}
func (v PublicIpConverter) GetResourceType() string {
	return "public_ip"
}
func (v KubernetesClusterConverter) GetResourceType() string {
	return "kubernetes_service"
}
func (v KubernetesNodePoolConverter) GetResourceType() string {
	return "kubernetes_node_pool"
}
func (v LambdaConverter) GetResourceType() string {
	return "lambda"
}
func (v VaultConverter) GetResourceType() string {
	return "vault"
}
func (v VaultAccessPolicyConverter) GetResourceType() string {
	return "vault_access_policy"
}
func (v VaultSecretConverter) GetResourceType() string {
	return "vault_secret"
}
func (v VirtualMachineConverter) GetResourceType() string {
	return "virtual_machine"
}

func (v VnConverter) ConvertToMultyResource(resourceId string, m proto.Message, _ map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVirtualNetworkArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vn := types.VirtualNetwork{
		CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:                 arg.Name,
		CidrBlock:            arg.CidrBlock,
	}
	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vn,
		ImplicitlyCreated: false,
	}, nil
}

func getCommonParams(resourceId string, arg *common.CloudSpecificResourceCommonArgs, c cloud_providers.CloudProvider) *common_resources.CommonResourceParams {
	return &common_resources.CommonResourceParams{
		ResourceId:             resourceId,
		ResourceGroupId:        arg.ResourceGroupId,
		Location:               strings.ToLower(arg.Location.String()),
		Clouds:                 []string{string(c)},
		ResourceValidationInfo: validate.NewResourceValidationInfoWithId(resourceId),
	}
}

func (v SubnetConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificSubnetArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	subnet := types.Subnet{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:             arg.Name,
		CidrBlock:        arg.CidrBlock,
		AvailabilityZone: int(arg.AvailabilityZone),
	}

	if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
		// Connect to vn in the same cloud
		subnet.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &subnet,
		ImplicitlyCreated: false,
	}, nil
}

func (v NetworkInterfaceConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificNetworkInterfaceArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	ni := types.NetworkInterface{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
	}

	if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SubnetId, c)]; ok {
		// Connect to subnet in the same cloud
		ni.SubnetId = subnet.Resource.(*types.Subnet)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", arg.SubnetId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &ni,
		ImplicitlyCreated: false,
	}, nil
}

func (v RouteTableConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificRouteTableArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	rt := types.RouteTable{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
		Routes: util.MapSliceValues(arg.Routes, func(route *resources.Route) types.RouteTableRoute {
			return types.RouteTableRoute{
				CidrBlock:   route.CidrBlock,
				Destination: strings.ToLower(route.Destination.String()),
			}
		}),
	}

	if arg.VirtualNetworkId != "" {
		if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
			// Connect to vn in the same cloud
			rt.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &rt,
		ImplicitlyCreated: false,
	}, nil
}

func (v RouteTableAssociationConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificRouteTableAssociationArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	rta := types.RouteTableAssociation{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c)}

	if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SubnetId, c)]; ok {
		// Connect to subnet in the same cloud
		rta.SubnetId = subnet.Resource.(*types.Subnet)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", arg.SubnetId, c)
	}

	if rt, ok := otherResources[common_resources.GetResourceIdForCloud(arg.RouteTableId, c)]; ok {
		// Connect to subnet in the same cloud
		rta.RouteTableId = rt.Resource.(*types.RouteTable)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("route table with id %s not found in %s", arg.RouteTableId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &rta,
		ImplicitlyCreated: false,
	}, nil
}

func (v NetworkSecurityGroupConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificNetworkSecurityGroupArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	nsg := types.NetworkSecurityGroup{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
		Rules: util.MapSliceValues(arg.Rules, func(rule *resources.NetworkSecurityRule) types.RuleType {
			return types.RuleType{
				Protocol:  rule.Protocol,
				Priority:  int(rule.Priority),
				FromPort:  convertPort(rule.PortRange.From),
				ToPort:    convertPort(rule.PortRange.To),
				CidrBlock: rule.CidrBlock,
				Direction: convertRuleDirection(rule.Direction),
			}
		}),
	}

	if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
		// Connect to vn in the same cloud
		nsg.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &nsg,
		ImplicitlyCreated: false,
	}, nil
}

func convertRuleDirection(direction resources.Direction) string {
	switch direction {
	case resources.Direction_BOTH_DIRECTIONS:
		return "both"
	case resources.Direction_INGRESS:
		return "ingress"
	case resources.Direction_EGRESS:
		return "egress"
	}

	return "unknown"
}

func convertPort(port int32) string {
	if port == 0 {
		return "*"
	}

	return strconv.Itoa(int(port))
}

func (v DatabaseConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificDatabaseArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	db := types.Database{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:          arg.Name,
		Engine:        strings.ToLower(arg.Engine.String()),
		EngineVersion: arg.EngineVersion,
		Storage:       int(arg.StorageGb),
		Size:          strings.ToLower(arg.Size.String()),
		DbUsername:    arg.Username,
		DbPassword:    arg.Password,
	}

	for _, subnetId := range arg.SubnetIds {

		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			db.SubnetIds = append(db.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &db,
		ImplicitlyCreated: false,
	}, nil
}

func (v ObjectStorageConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificObjectStorageArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	db := types.ObjectStorage{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &db,
		ImplicitlyCreated: false,
	}, nil
}

func (v ObjectStorageObjectConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificObjectStorageObjectArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	obj := types.ObjectStorageObject{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:        arg.Name,
		Content:     arg.Content,
		ContentType: arg.ContentType,
		Acl:         strings.ToLower(arg.Acl.String()),
		Source:      arg.Source,
	}

	if o, ok := otherResources[common_resources.GetResourceIdForCloud(arg.ObjectStorageId, c)]; ok {
		// Connect to vn in the same cloud
		obj.ObjectStorage = o.Resource.(*types.ObjectStorage)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("object storage with id %s not found in %s", arg.ObjectStorageId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &obj,
		ImplicitlyCreated: false,
	}, nil
}

func (v PublicIpConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificPublicIpArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	obj := types.PublicIp{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
	}

	if ni, ok := otherResources[common_resources.GetResourceIdForCloud(arg.NetworkInterfaceId, c)]; ok {
		// Connect to vn in the same cloud
		obj.NetworkInterfaceId = ni.Resource.(*types.NetworkInterface)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("network interface with id %s not found in %s", arg.NetworkInterfaceId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &obj,
		ImplicitlyCreated: false,
	}, nil
}

func (v KubernetesClusterConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificKubernetesClusterArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	kc := types.KubernetesService{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
	}

	for _, subnetId := range arg.SubnetIds {

		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			kc.SubnetIds = append(kc.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &kc,
		ImplicitlyCreated: false,
	}, nil
}

func zeroToNil(a int32) *int {
	var result *int
	if a != 0 {
		n := int(a)
		result = &n
	}
	return result
}

func (v KubernetesNodePoolConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificKubernetesNodePoolArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	knp := types.KubernetesServiceNodePool{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:              arg.Name,
		IsDefaultPool:     arg.IsDefaultPool,
		StartingNodeCount: zeroToNil(arg.StartingNodeCount),
		MaxNodeCount:      int(arg.MaxNodeCount),
		MinNodeCount:      int(arg.MinNodeCount),
		Labels:            arg.Labels,
		VmSize:            strings.ToLower(arg.VmSize.String()),
		DiskSizeGiB:       int(arg.DiskSizeGb),
	}

	if kc, ok := otherResources[common_resources.GetResourceIdForCloud(arg.ClusterId, c)]; ok {
		// Connect to vn in the same cloud
		knp.ClusterId = kc.Resource.(*types.KubernetesService)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("cluster with id %s not found in %s", arg.ClusterId, c)
	}

	for _, subnetId := range arg.SubnetIds {
		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			knp.SubnetIds = append(knp.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &knp,
		ImplicitlyCreated: false,
	}, nil
}

func (v LambdaConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificLambdaArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	l := types.Lambda{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		FunctionName: arg.Name,
		Runtime:      arg.Runtime,
	}

	if o, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SourceCodeObjectId, c)]; ok {
		// Connect to vn in the same cloud
		l.SourceCodeObject = o.Resource.(*types.ObjectStorageObject)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("object with id %s not found in %s", arg.SourceCodeObjectId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &l,
		ImplicitlyCreated: false,
	}, nil
}

func (v VaultConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVaultArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vault := types.Vault{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name: arg.Name,
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vault,
		ImplicitlyCreated: false,
	}, nil
}

func (v VaultAccessPolicyConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVaultAccessPolicyArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vap := types.VaultAccessPolicy{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Identity: arg.Identity,
		Access:   strings.ToLower(arg.Access.String()),
	}

	if v, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VaultId, c)]; ok {
		// Connect to vn in the same cloud
		vap.Vault = v.Resource.(*types.Vault)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("vault with id %s not found in %s", arg.VaultId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vap,
		ImplicitlyCreated: false,
	}, nil
}

func (v VaultSecretConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVaultSecretArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vs := types.VaultSecret{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:  arg.Name,
		Value: arg.Value,
	}

	if v, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VaultId, c)]; ok {
		// Connect to vn in the same cloud
		vs.Vault = v.Resource.(*types.Vault)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("vault with id %s not found in %s", arg.VaultId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vs,
		ImplicitlyCreated: false,
	}, nil
}

func (v VirtualMachineConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVirtualMachineArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vm := types.VirtualMachine{CommonResourceParams: getCommonParams(resourceId, arg.CommonParameters, c),
		Name:            arg.Name,
		OperatingSystem: strings.ToLower(arg.OperatingSystem.String()),
		Size:            strings.ToLower(arg.VmSize.String()),
		UserData:        arg.UserData,
		SshKey:          arg.PublicSshKey,
		PublicIp:        arg.GeneratePublicIp,
	}

	if arg.PublicIpId != "" {
		if pid, ok := otherResources[common_resources.GetResourceIdForCloud(arg.PublicIpId, c)]; ok {
			vm.PublicIpId = pid.Resource.(*types.PublicIp)
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("public ip with id %s not found in %s", arg.PublicIpId, c)
		}
	}

	if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SubnetId, c)]; ok {
		vm.SubnetId = subnet.Resource.(*types.Subnet)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", arg.SubnetId, c)
	}

	for _, niId := range arg.NetworkInterfaceIds {
		if ni, ok := otherResources[common_resources.GetResourceIdForCloud(niId, c)]; ok {
			vm.NetworkInterfaceIds = append(vm.NetworkInterfaceIds, ni.Resource.(*types.NetworkInterface))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("network interface with id %s not found in %s", niId, c)
		}
	}

	for _, nsgId := range arg.NetworkSecurityGroupIds {
		if nsg, ok := otherResources[common_resources.GetResourceIdForCloud(nsgId, c)]; ok {
			vm.NetworkSecurityGroupIds = append(vm.NetworkSecurityGroupIds, nsg.Resource.(*types.NetworkSecurityGroup))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("network security group with id %s not found in %s", nsgId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vm,
		ImplicitlyCreated: false,
	}, nil
}
