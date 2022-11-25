package wasabi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"log"
)

func Client() *s3.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           viper.GetString("connection.url"),
			SigningRegion: viper.GetString("connection.region"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(viper.GetString("connection.profile")), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}
