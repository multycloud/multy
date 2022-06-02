package aws_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/multycloud/multy/flags"
	"log"
	"net/http"
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
		fmt.Println(err)
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

func (c Client) UpdateQPSMetric(apiKey string, service string, method string) error {
	if flags.DryRun || flags.NoTelemetry {
		return nil
	}
	postBody, _ := json.Marshal(map[string]string{
		"action":  method,
		"service": service,
		"api_key": apiKey,
	})
	resp, err := http.Post(logUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("Logging error occured %v", err)
		return err
	}
	defer resp.Body.Close()

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
		},
		MetricName: aws.String("qps"),
		Timestamp:  aws.Time(time.Now()),
		Value:      aws.Float64(1),
	}
	_, err = c.cloudWatchClient.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{metric},
		Namespace:  aws.String("multy/server/"),
	})
	if err != nil {
		log.Printf("[WARNING] %s\n", err.Error())
	}
	return err
}

func (c Client) UpdateErrorMetric(service string, method string, errorCode string) error {
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
		log.Printf("[WARNING] %s\n", err.Error())
	}
	return err
}
