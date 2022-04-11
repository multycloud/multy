package aws_client

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	region     = "eu-west-2"
	bucketName = "multy-users-tfstate"
)

type Client struct {
	s3Client *s3.S3
}

func Configure() Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)},
	))
	return Client{s3.New(sess)}
}

func (c Client) SaveFile(userId string, fileName string, content string) error {
	keyName := fmt.Sprintf("%s/%s", userId, fileName)

	_, err := c.s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("private"),
		Body:   bytes.NewReader([]byte(content)),
		Bucket: aws.String(bucketName),
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
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Printf("s3://%s/%s does not exist. Creating empty file\n", bucketName, keyName)
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
