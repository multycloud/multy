package database

import (
	"fmt"

	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_db_instance"

// aws_db_instance
type AwsDbInstance struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_db_instance"`
	AllocatedStorage    int      `hcl:"allocated_storage" json:"allocated_storage"`
	Name                string   `hcl:"db_name" hcle:"omitempty" json:"name"`
	Engine              string   `hcl:"engine" json:"engine"`
	EngineVersion       string   `hcl:"engine_version" json:"engine_version"`
	Username            string   `hcl:"username" json:"username"`
	Password            string   `hcl:"password" json:"password"`
	InstanceClass       string   `hcl:"instance_class" json:"instance_class"`
	Identifier          string   `hcl:"identifier" json:"identifier"`
	SkipFinalSnapshot   bool     `hcl:"skip_final_snapshot" json:"skip_final_snapshot"`
	DbSubnetGroupName   string   `hcl:"db_subnet_group_name,expr" json:"db_subnet_group_name"`
	PubliclyAccessible  bool     `hcl:"publicly_accessible" json:"publicly_accessible"`
	VpcSecurityGroupIds []string `hcl:"vpc_security_group_ids,expr" json:"vpc_security_group_ids"`
	Port                int      `hcl:"port" hcle:"omitempty" json:"port"`

	// outputs
	Address string `json:"address" hcle:"omitempty" json:"address"`
	Arn     string `json:"arn" hcle:"omitempty" json:"arn"`
}

// aws_db_subnet_group
type AwsDbSubnetGroup struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_db_subnet_group"`
	Name                string   `hcl:"name"`
	Description         string   `hcl:"description"`
	SubnetIds           []string `hcl:"subnet_ids,expr"`

	Arn string `json:"arn" hcle:"omitempty"`
}

func (dbSubGroup AwsDbSubnetGroup) GetResourceName() string {
	return fmt.Sprintf("aws_db_subnet_group.%s.name", dbSubGroup.ResourceId)
}
