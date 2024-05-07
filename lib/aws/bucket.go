package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateBucket(name string) error {
	client := s3.NewFromConfig(AwsConfig)
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: &name})
	if err != nil {
		return err
	}
	return nil
}
