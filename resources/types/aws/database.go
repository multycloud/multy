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
	"github.com/multycloud/multy/util"
	"strings"
)

type AwsDatabase struct {
	*types.Database
}

func InitDatabase(r *types.Database) resources.ResourceTranslator[*resourcespb.DatabaseResource] {
	return AwsDatabase{r}
}

func (r AwsDatabase) FromState(state *output.TfState) (*resourcespb.DatabaseResource, error) {
	host := "dyrun"
	if !flags.DryRun {
		values, err := state.GetValues(database.AwsDbInstance{}, r.ResourceId)
		if err != nil {
			return nil, err
		}
		host = values["address"].(string)
	}

	return &resourcespb.DatabaseResource{
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
		Host:               host,
		ConnectionUsername: r.Args.Username,
	}, nil
}

func (r AwsDatabase) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds, err := util.MapSliceValuesErr(r.Subnets, func(v *types.Subnet) (string, error) {
		return resources.GetMainOutputId(AwsSubnet{v})
	})
	if err != nil {
		return nil, err
	}
	vpcId, err := resources.GetMainOutputId(AwsVirtualNetwork{r.Subnets[0].VirtualNetwork})
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
