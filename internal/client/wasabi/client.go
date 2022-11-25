package wasabi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	wasabiConfig "wasabiCleanup/internal/config"
)

func Client(appConfig wasabiConfig.S3Connection) *s3.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           appConfig.Url,
			SigningRegion: appConfig.Region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(appConfig.Profile), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}
