package database

import (
	"fmt"

	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_db_instance"

// aws_db_instance
type AwsDbInstance struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_db_instance"`
	AllocatedStorage    int      `hcl:"allocated_storage"`
	Name                string   `hcl:"db_name" hcle:"omitempty"`
	Engine              string   `hcl:"engine"`
	EngineVersion       string   `hcl:"engine_version"`
	Username            string   `hcl:"username"`
	Password            string   `hcl:"password"`
	InstanceClass       string   `hcl:"instance_class"`
	Identifier          string   `hcl:"identifier"`
	SkipFinalSnapshot   bool     `hcl:"skip_final_snapshot"`
	DbSubnetGroupName   string   `hcl:"db_subnet_group_name,expr"`
	PubliclyAccessible  bool     `hcl:"publicly_accessible"`
	VpcSecurityGroupIds []string `hcl:"vpc_security_group_ids,expr"`
	Port                int      `hcl:"port" hcle:"omitempty"`

	// outputs
	Address string `json:"address"`
}

// aws_db_subnet_group
type AwsDbSubnetGroup struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_db_subnet_group"`
	Name                string   `hcl:"name"`
	Description         string   `hcl:"description"`
	SubnetIds           []string `hcl:"subnet_ids,expr"`
}

func (dbSubGroup AwsDbSubnetGroup) GetResourceName() string {
	return fmt.Sprintf("aws_db_subnet_group.%s.name", dbSubGroup.ResourceId)
}
