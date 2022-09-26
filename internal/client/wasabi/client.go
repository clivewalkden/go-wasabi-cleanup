package wasabi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

func Client() *s3.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "https://s3.eu-central-1.wasabisys.com",
			SigningRegion: "eu-central-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("wasabi"), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}
