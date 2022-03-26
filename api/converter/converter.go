package converter

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
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

func getCommonChildResourceParams(resourceId string, arg *common.CloudSpecificChildResourceCommonArgs) *common_resources.CommonResourceParams {
	return &common_resources.CommonResourceParams{
		ResourceId:             resourceId,
		ResourceValidationInfo: validate.NewResourceValidationInfoWithId(resourceId),
	}
}

func setReference[T any](resources map[string]common_resources.CloudSpecificResource, id string, field *T, resourceId string, fieldName string) error {
	if vn, ok := resources[id]; ok {
		if _, okType := vn.Resource.(T); !okType {
			return errors.ValidationErrors([]validate.ValidationError{
				{
					ErrorMessage: fmt.Sprintf("resource with id %s is of the wrong type", id),
					ResourceId:   resourceId,
					FieldName:    fieldName,
				},
			})
		}
		// Connect to resource in the same cloud
		*field = vn.Resource.(T)
	} else {
		return errors.ValidationErrors([]validate.ValidationError{
			{
				ErrorMessage: fmt.Sprintf("resource with id %s not found", id),
				ResourceId:   resourceId,
				FieldName:    fieldName,
			},
		})
	}

	return nil
}

func (v SubnetConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificSubnetArgs)
	subnet := types.Subnet{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters),
		Name:             arg.Name,
		CidrBlock:        arg.CidrBlock,
		AvailabilityZone: int(arg.AvailabilityZone),
	}

	err := setReference(otherResources, arg.VirtualNetworkId, &subnet.VirtualNetwork, resourceId, "virtual_network_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(subnet.VirtualNetwork.Clouds[0]),
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

	err := setReference(otherResources, arg.SubnetId, &ni.SubnetId, resourceId, "subnet_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
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
		err := setReference(otherResources, arg.VirtualNetworkId, &rt.VirtualNetwork, resourceId, "virtual_network_id")
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
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
	rta := types.RouteTableAssociation{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters)}

	err := setReference(otherResources, arg.SubnetId, &rta.SubnetId, resourceId, "subnet_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	err = setReference(otherResources, arg.RouteTableId, &rta.RouteTableId, resourceId, "route_table_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(rta.RouteTableId.Clouds[0]),
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

	err := setReference(otherResources, arg.VirtualNetworkId, &nsg.VirtualNetwork, resourceId, "virtual_network_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
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

	db.SubnetIds = make([]*types.Subnet, len(arg.SubnetIds))
	for i, subnetId := range arg.SubnetIds {
		err := setReference(otherResources, subnetId, &db.SubnetIds[i], resourceId, fmt.Sprintf("subnet_ids[%d]", i))
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
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
		Name:       arg.Name,
		Versioning: arg.Versioning,
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &db,
		ImplicitlyCreated: false,
	}, nil
}

func (v ObjectStorageObjectConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificObjectStorageObjectArgs)
	obj := types.ObjectStorageObject{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters),
		Name:        arg.Name,
		Content:     arg.Content,
		ContentType: arg.ContentType,
		Acl:         strings.ToLower(arg.Acl.String()),
		Source:      arg.Source,
	}

	err := setReference(otherResources, arg.ObjectStorageId, &obj.ObjectStorage, resourceId, "object_storage_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(obj.ObjectStorage.Clouds[0]),
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

	err := setReference(otherResources, arg.NetworkInterfaceId, &obj.NetworkInterfaceId, resourceId, "network_interface_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
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

	kc.SubnetIds = make([]*types.Subnet, len(arg.SubnetIds))
	for i, subnetId := range arg.SubnetIds {
		err := setReference(otherResources, subnetId, &kc.SubnetIds[i], resourceId, fmt.Sprintf("subnet_ids[%d]", i))
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
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
	knp := types.KubernetesServiceNodePool{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters),
		Name:              arg.Name,
		IsDefaultPool:     arg.IsDefaultPool,
		StartingNodeCount: zeroToNil(arg.StartingNodeCount),
		MaxNodeCount:      int(arg.MaxNodeCount),
		MinNodeCount:      int(arg.MinNodeCount),
		Labels:            arg.Labels,
		VmSize:            strings.ToLower(arg.VmSize.String()),
		DiskSizeGiB:       int(arg.DiskSizeGb),
	}

	err := setReference(otherResources, arg.ClusterId, &knp.ClusterId, resourceId, "cluster_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	knp.SubnetIds = make([]*types.Subnet, len(arg.SubnetIds))
	for i, subnetId := range arg.SubnetIds {
		err := setReference(otherResources, subnetId, &knp.SubnetIds[i], resourceId, fmt.Sprintf("subnet_ids[%d]", i))
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(knp.ClusterId.Clouds[0]),
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

	err := setReference(otherResources, arg.SourceCodeObjectId, &l.SourceCodeObject, resourceId, "source_code_object_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
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
	vap := types.VaultAccessPolicy{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters),
		Identity: arg.Identity,
		Access:   strings.ToLower(arg.Access.String()),
	}

	err := setReference(otherResources, arg.VaultId, &vap.Vault, resourceId, "vault_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(vap.Vault.Clouds[0]),
		Resource:          &vap,
		ImplicitlyCreated: false,
	}, nil
}

func (v VaultSecretConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVaultSecretArgs)
	vs := types.VaultSecret{CommonResourceParams: getCommonChildResourceParams(resourceId, arg.CommonParameters),
		Name:  arg.Name,
		Value: arg.Value,
	}

	err := setReference(otherResources, arg.VaultId, &vs.Vault, resourceId, "vault_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	return common_resources.CloudSpecificResource{
		Cloud:             cloud_providers.CloudProvider(vs.Vault.Clouds[0]),
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
		err := setReference(otherResources, arg.PublicIpId, &vm.PublicIpId, resourceId, "public_ip_id")
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
		}
	}

	err := setReference(otherResources, arg.SubnetId, &vm.SubnetId, resourceId, "subnet_id")
	if err != nil {
		return common_resources.CloudSpecificResource{}, err
	}

	vm.NetworkInterfaceIds = make([]*types.NetworkInterface, len(arg.NetworkInterfaceIds))
	for i, niId := range arg.NetworkInterfaceIds {
		err := setReference(otherResources, niId, &vm.NetworkInterfaceIds[i], resourceId, fmt.Sprintf("network_interface_ids[%d]", i))
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
		}
	}

	vm.NetworkSecurityGroupIds = make([]*types.NetworkSecurityGroup, len(arg.NetworkSecurityGroupIds))
	for i, nsgId := range arg.NetworkSecurityGroupIds {
		err := setReference(otherResources, nsgId, &vm.NetworkSecurityGroupIds[i], resourceId, fmt.Sprintf("network_security_group_ids[%d]", i))
		if err != nil {
			return common_resources.CloudSpecificResource{}, err
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vm,
		ImplicitlyCreated: false,
	}, nil
}
