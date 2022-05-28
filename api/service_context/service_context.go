package service_context

import (
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/db"
)

type ServiceContext struct {
	db.Database
	aws_client.AwsClient
	deploy.DeploymentExecutor
}
