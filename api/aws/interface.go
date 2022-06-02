package aws_client

import "github.com/multycloud/multy/flags"

type AwsClient interface {
	SaveFile(userId string, fileName string, content string) error
	ReadFile(userId string, fileName string) (string, error)
	UpdateQPSMetric(apiKey string, service string, method string) error
	UpdateErrorMetric(service string, method string, errorCode string) error
}

func NewClient() (AwsClient, error) {
	if flags.Environment == flags.Local {
		return newLocalClient()
	}

	return newS3Client()
}
