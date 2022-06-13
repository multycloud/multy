package aws_client

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/multycloud/multy/flags"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"
)

const (
	region = "eu-west-2"
	logUrl = "https://c42wtfa4z6.execute-api.eu-west-1.amazonaws.com/logs"
)

type Client struct {
	s3Client         *s3.S3
	cloudWatchClient *cloudwatch.CloudWatch
	userStorageName  string
}

func newS3Client() (*Client, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)},
	))
	userStorageName, exists := os.LookupEnv("USER_STORAGE_NAME")
	if !exists {
		return nil, fmt.Errorf("USER_STORAGE_NAME not found")
	}
	return &Client{s3.New(sess), cloudwatch.New(sess), userStorageName}, nil
}

func (c Client) SaveFile(userId string, fileName string, content string) error {
	keyName := fmt.Sprintf("%s/%s", userId, fileName)

	_, err := c.s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("private"),
		Body:   bytes.NewReader([]byte(content)),
		Bucket: aws.String(c.userStorageName),
		Key:    aws.String(keyName),
	})

	if err != nil {
		log.Printf("[ERROR] Error saving file to S3: %s\n", err)
		return err
	}

	return nil
}

func (c Client) ReadFile(userId string, fileName string) (string, error) {
	keyName := fmt.Sprintf("%s/%s", userId, fileName)

	object, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.userStorageName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Printf("s3://%s/%s does not exist. Creating empty file\n", c.userStorageName, keyName)
				err := c.SaveFile(userId, fileName, "")
				if err != nil {
					return "", err
				}
				return "", nil
			default:
				return "", aerr
			}

		} else {
			return "", err
		}
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(object.Body)
	if err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}

func (c Client) UpdateQPSMetric(userId string, service string, method string) error {
	if flags.DryRun || flags.NoTelemetry {
		return nil
	}
	var wg errgroup.Group
	wg.Go(func() error {
		err := logAction(userId, service, method)
		if err != nil {
			log.Printf("[WARNING] Logging error ocurred: %s\n", err)
		}
		return err
	})

	wg.Go(func() error {
		metric := &cloudwatch.MetricDatum{
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String("service"),
					Value: aws.String(service),
				},
				{
					Name:  aws.String("method"),
					Value: aws.String(method),
				},
				{
					Name:  aws.String("env"),
					Value: aws.String(string(flags.Environment)),
				},
			},
			MetricName: aws.String("qps"),
			Timestamp:  aws.Time(time.Now()),
			Value:      aws.Float64(1),
		}
		_, err := c.cloudWatchClient.PutMetricData(&cloudwatch.PutMetricDataInput{
			MetricData: []*cloudwatch.MetricDatum{metric},
			Namespace:  aws.String("multy/server/"),
		})
		if err != nil {
			log.Printf("[WARNING] Cloudwatch error ocurred: %s\n", err.Error())
			return err
		}
		return nil
	})

	return wg.Wait()
}

func (c Client) UpdateErrorMetric(service string, method string, errorCode string) error {
	if flags.DryRun || flags.NoTelemetry {
		return nil
	}
	metric := &cloudwatch.MetricDatum{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("service"),
				Value: aws.String(service),
			},
			{
				Name:  aws.String("method"),
				Value: aws.String(method),
			},
			{
				Name:  aws.String("env"),
				Value: aws.String(string(flags.Environment)),
			},
			{
				Name:  aws.String("error_code"),
				Value: aws.String(errorCode),
			},
		},
		MetricName: aws.String("error"),
		Timestamp:  aws.Time(time.Now()),
		Value:      aws.Float64(1),
	}
	_, err := c.cloudWatchClient.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{metric},
		Namespace:  aws.String("multy/server/"),
	})
	if err != nil {
		log.Printf("[WARNING] Cloudwatch error ocurred: %s\n", err.Error())
	}
	return err
}
