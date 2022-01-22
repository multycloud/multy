package database

import (
	"fmt"
	"multy-go/resources/common"
)

// aws_db_instance
type AwsDbInstance struct {
	common.AwsResource `hcl:",squash"`
	AllocatedStorage   int    `hcl:"allocated_storage"`
	Name               string `hcl:"name"`
	Engine             string `hcl:"engine"`
	EngineVersion      string `hcl:"engine_version"`
	Username           string `hcl:"username"`
	Password           string `hcl:"password"`
	InstanceClass      string `hcl:"instance_class"`
	Identifier         string `hcl:"identifier"`
	SkipFinalSnapshot  bool   `hcl:"skip_final_snapshot"`
	DbSubnetGroupName  string `hcl:"db_subnet_group_name,expr"`
	PubliclyAccessible bool   `hcl:"publicly_accessible"`
}

// aws_db_subnet_group
type AwsDbSubnetGroup struct {
	common.AwsResource `hcl:",squash"`
	SubnetIds          []string `hcl:"subnet_ids"`
}

func (dbSubGroup AwsDbSubnetGroup) GetResourceName() string {
	return fmt.Sprintf("aws_db_subnet_group.%s.name", dbSubGroup.ResourceId)
}
