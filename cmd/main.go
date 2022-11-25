package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
	"wasabiCleanup/internal/client/wasabi"
	wasabiConfig "wasabiCleanup/internal/config"
	"wasabiCleanup/internal/reporting"
	"wasabiCleanup/internal/utils"
)

type S3Object struct {
	Key          string
	LastModified time.Time
	Size         int64
}

type S3Objects struct {
	Items []types.ObjectIdentifier
	Size  int64
}

func main() {
	appConfig := wasabiConfig.InitConfig()

	if (appConfig.Connection == wasabiConfig.S3Connection{}) {
		log.Fatal("Sorry configuration is incomplete and can't connect to the S3 instance.")
		os.Exit(1)
	}

	client := wasabi.Client()
	report := reporting.Report{}

	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Working...")
	for _, object := range buckets.Buckets {
		if viper.GetBool("verbose") {
			fmt.Printf("Checking Bucket %s\n", *object.Name)
		}

		if appConfig.Buckets[*object.Name] == 0 {
			if viper.GetBool("verbose") {
				fmt.Printf("\t- Bucket not in config, skipping\n")
			}
			continue
		}

		// Return files that need deleting from this bucket based on the Retention Policy
		objectList := S3Objects{}
		safeList := S3Objects{}
		maxKeys := 0

		params := &s3.ListObjectsV2Input{Bucket: object.Name}

		// Create the Paginator for the ListObjectsV2 operation.
		p := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
			if v := int32(maxKeys); v != 0 {
				o.Limit = v
			}
		})

		// The date we need to delete items prior to
		comparisonDate := time.Now().AddDate(0, 0, -appConfig.Buckets[*object.Name]-1)
		if viper.GetBool("verbose") {
			fmt.Printf("\t- Checking files date is before %s\n", comparisonDate)
		}

		// Iterate through the S3 object pages, printing each object returned.
		var i int
		for p.HasMorePages() {
			i++

			// Next Page takes a new context for each page retrieval. This is where
			// you could add timeouts or deadlines.
			page, err := p.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("\t\tfailed to get page %v, %v", i, err)
			}

			if viper.GetBool("verbose") {
				fmt.Printf("\t\t- Next page (%d)\n", i)
			}

			// Log the objects found
			for _, obj := range page.Contents {
				if obj.LastModified.Before(comparisonDate) {
					objectList.Items = append(objectList.Items, types.ObjectIdentifier{
						Key: obj.Key,
					})
					objectList.Size += obj.Size
					if viper.GetBool("verbose") {
						fmt.Printf("\t\t\t- Deleting object %s\n", *obj.Key)
					}
					_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
						Bucket: object.Name,
						Key:    obj.Key,
					})

					if err != nil {
						panic("Couldn't delete items")
					}
				} else {
					safeList.Items = append(safeList.Items, types.ObjectIdentifier{
						Key: obj.Key,
					})
					safeList.Size += obj.Size
				}
			}
		}

		result := reporting.Result{
			Name:        *object.Name,
			Kept:        len(safeList.Items),
			KeptSize:    utils.ByteCountSI(safeList.Size),
			Deleted:     len(objectList.Items),
			DeletedSize: utils.ByteCountSI(objectList.Size),
		}

		report.Result = append(report.Result, result)
	}

	reporting.Output(report)
}
