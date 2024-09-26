package aws

import (
	"bialekredki/atik/lib/strings"
	"context"
	"fmt"
	stdstrings "strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var operationalDirectories = map[string]types.StorageClass{
	"hot/":     types.StorageClassStandard,
	"cold/":    types.StorageClassGlacier,
	"archive/": types.StorageClassDeepArchive,
}

func formatToS3DirectoryObjectName(value string) string {
	if stdstrings.HasSuffix(value, "/") {
		return value
	}
	return fmt.Sprintf("%s/", value)
}

type S3RepositoryParams struct {
	fx.In
	AwsConfig *aws.Config
	Logger    *zap.Logger
}

type S3Repository struct {
	client                *s3.Client
	presignClient         *s3.PresignClient
	logger                *zap.Logger
	operationalBucketName string
	awsRegion             string
}

func NewS3Repository(p S3RepositoryParams) (*S3Repository, error) {
	ctx := context.Background()
	client := s3.NewFromConfig(*p.AwsConfig)
	repository := &S3Repository{
		client:                client,
		presignClient:         s3.NewPresignClient(client),
		logger:                p.Logger,
		operationalBucketName: fmt.Sprintf("atik-operational-bucket-%s", p.AwsConfig.Region),
		awsRegion:             p.AwsConfig.Region,
	}
	_, err := repository.CreateOperationalBucket(ctx)
	if err != nil {
		return nil, err
	}
	err = repository.createOperationalDirectories(ctx)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func (r *S3Repository) CreateOperationalBucket(ctx context.Context) (bool, error) {
	_, err := r.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &r.operationalBucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(r.awsRegion),
		},
	})
	if err != nil {
		if apiErr := ToAPIError[*types.BucketAlreadyExists](err); apiErr != nil {
			LogAwsError(err)
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *S3Repository) CreateBucket(name string, ctx context.Context) (bool, error) {
	encodedName := strings.MD5Encode(name)
	_, err := r.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &encodedName, CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: "eu-central-1"}})
	if err != nil {
		if apiErr := ToAPIError[*types.BucketAlreadyExists](err); apiErr == nil {
			LogAwsError(err)
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (r *S3Repository) CreateDirectory(ctx context.Context, bucket string, name string, storageClass types.StorageClass) error {
	name = formatToS3DirectoryObjectName(name)
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &name, StorageClass: storageClass})
	return err
}

func (r *S3Repository) createOperationalDirectories(ctx context.Context) error {
	r.logger.Info("creating operational directories for ", zap.String("bucket", r.operationalBucketName))
	for name, storageClass := range operationalDirectories {
		r.logger.Info("creating operation directory", zap.String("directory", name))
		err := r.CreateDirectory(ctx, r.operationalBucketName, name, storageClass)
		if err != nil {
			r.logger.Sugar().Error(err)
			return err
		}
	}
	return nil
}

func (r *S3Repository) PresignUploadFile(bucket string, dir string, file string, ctx context.Context) {
}
