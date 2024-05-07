package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var AwsConfig aws.Config

func LoadConfig() *aws.Config {
	loadedConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config")
	}

	AwsConfig = loadedConfig
	return &loadedConfig
}
