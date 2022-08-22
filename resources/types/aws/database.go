package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/types"
	"strings"
)

type AwsDatabase struct {
	*types.Database
}

func InitDatabase(r *types.Database) resources.ResourceTranslator[*resourcespb.DatabaseResource] {
	return AwsDatabase{r}
}

func (r AwsDatabase) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.DatabaseResource, error) {
	out := &resourcespb.DatabaseResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:               r.Args.Name,
		Engine:             r.Args.Engine,
		EngineVersion:      r.Args.EngineVersion,
		StorageGb:          r.Args.StorageGb,
		Size:               r.Args.Size,
		Username:           r.Args.Username,
		Password:           r.Args.Password,
		SubnetIds:          r.Args.SubnetIds,
		Port:               r.Args.Port,
		SubnetId:           r.Args.SubnetId,
		Host:               "",
		ConnectionUsername: r.Args.Username,
		GcpOverride:        r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}
	out.AwsOutputs = &resourcespb.DatabaseAwsOutputs{}

	if db, exists, err := output.MaybeGetParsedById[database.AwsDbInstance](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Host = db.Address
		out.Name = db.AwsResource.Tags["Name"]
		out.Username = db.Username
		out.Password = db.Password
		if e, ok := resourcespb.DatabaseEngine_value[strings.ToUpper(db.Engine)]; ok {
			out.Engine = resourcespb.DatabaseEngine(e)
		} else {
			out.Engine = resourcespb.DatabaseEngine_UNKNOWN_ENGINE
		}
		out.EngineVersion = db.EngineVersion
		out.Port = int32(db.Port)
		out.StorageGb = int64(db.AllocatedStorage)
		out.AwsOutputs.DbInstanceId = db.Arn
		output.AddToStatuses(statuses, "aws_db_instance", output.MaybeGetPlannedChageById[database.AwsDbInstance](plan, r.ResourceId))
	} else {
		statuses["aws_db_instance"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.AwsSecurityGroup](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.DefaultNetworkSecurityGroupId = stateResource.ResourceId
		output.AddToStatuses(statuses, "aws_default_network_security_group", output.MaybeGetPlannedChageById[network_security_group.AwsSecurityGroup](plan, r.ResourceId))
	} else {
		statuses["aws_default_network_security_group"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[database.AwsDbSubnetGroup](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.DbSubnetGroupId = stateResource.Arn
		output.AddToStatuses(statuses, "aws_db_subnet_group", output.MaybeGetPlannedChageById[database.AwsDbSubnetGroup](plan, r.ResourceId))
	} else {
		statuses["aws_db_subnet_group"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil

}

func (r AwsDatabase) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds := AwsSubnet{r.Subnet}.GetSubnetIds()
	vpcId, err := resources.GetMainOutputId(AwsVirtualNetwork{r.Subnet.VirtualNetwork})
	if err != nil {
		return nil, err
	}

	name := common.RemoveSpecialChars(r.Args.Name)
	dbSubnetGroup := database.AwsDbSubnetGroup{
		AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
		Name:        r.Args.Name,
		Description: "Managed by Multy",
		SubnetIds:   subnetIds,
	}
	nsg := network_security_group.AwsSecurityGroup{
		AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
		VpcId:       vpcId,
		Name:        r.Args.Name,
		Description: fmt.Sprintf("Default security group of %s", r.Args.Name),
		Ingress: []network_security_group.AwsSecurityGroupRule{{
			Protocol:   "-1",
			FromPort:   0,
			ToPort:     0,
			CidrBlocks: []string{"0.0.0.0/0"},
		}},
		Egress: []network_security_group.AwsSecurityGroupRule{{
			Protocol:   "-1",
			FromPort:   0,
			ToPort:     0,
			CidrBlocks: []string{"0.0.0.0/0"},
		}},
	}
	return []output.TfBlock{
		dbSubnetGroup,
		nsg,
		database.AwsDbInstance{
			AwsResource:         common.NewAwsResource(r.ResourceId, name),
			AllocatedStorage:    int(r.Args.StorageGb),
			Engine:              strings.ToLower(r.Args.Engine.String()),
			EngineVersion:       r.Args.EngineVersion,
			Username:            r.Args.Username,
			Password:            r.Args.Password,
			InstanceClass:       common.DBSIZE[r.Args.Size][r.GetCloud()],
			Identifier:          r.Args.Name,
			SkipFinalSnapshot:   true,
			DbSubnetGroupName:   dbSubnetGroup.GetResourceName(),
			PubliclyAccessible:  true,
			VpcSecurityGroupIds: []string{fmt.Sprintf("%s.%s.id", output.GetResourceName(nsg), nsg.ResourceId)},
			Port:                int(r.Args.Port),
		},
	}, nil
}

func (r AwsDatabase) GetMainResourceName() (string, error) {
	return database.AwsResourceName, nil
}
