package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/wal1251/pkg/tools/size"
)

type (
	Client struct {
		client      s3.Client
		config      Config
		urlProvider func(string) *url.URL
	}

	StorageObject struct {
		Key          string
		ETag         string
		Size         size.Size
		LastModified *time.Time
		Body         io.ReadCloser
		Expires      *time.Time
		Metadata     map[string]string
		PartNumber   int32
	}

	FileObject struct {
		Name        string
		Body        io.Reader
		Metadata    map[string]string
		ContentType *string
		PartNumber  int32
		Expires     *time.Time
	}
)

func NewClient(ctx context.Context, cfg *Config) (ClientInterface, error) {
	storageURL, err := url.Parse(cfg.ServiceURL)
	if err != nil {
		return nil, fmt.Errorf("incorrect storage url: %s: %w", cfg.ServiceURL, err)
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			HostnameImmutable: false,
			PartitionID:       PartitionID,
			URL:               storageURL.String(),
			SigningRegion:     SigningRegion,
		}, nil
	})

	s3Config, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, ""),
		),
		config.WithRegion(SigningRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMode(aws.RetryModeStandard),
		config.WithRetryMaxAttempts(cfg.RetryCount),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(s3Config)

	_, err = client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(cfg.BucketName),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to access bucket %s: %w", cfg.BucketName, err)
	}

	return &Client{
		client:      *client,
		config:      *cfg,
		urlProvider: makeURLProvider(storageURL, cfg.BucketName),
	}, nil
}
