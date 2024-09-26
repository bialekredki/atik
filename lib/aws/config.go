package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.Logger) (*aws.Config, error) {
	loadedConfig, err := config.LoadDefaultConfig(context.Background())
	logger.Sugar().Debugf("aws config %+v", loadedConfig)
	if err != nil {
		logger.Fatal("failed to load AWS config", zap.Error(err))
		return nil, err
	}

	return &loadedConfig, nil
}
