package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"time"
	"wasabiCleanup/internal/client/wasabi"
)

type S3Object struct {
	Key          string
	LastModified time.Time
}

type S3Objects struct {
	Items []S3Object
}

var retention map[string]int

func main() {
	retention = make(map[string]int)
	retention["sozo-db-backups"] = 90
	retention["sozo-log-backups"] = 180

	client := wasabi.Client()

	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("First page results:")
	for _, object := range buckets.Buckets {
		if retention[*object.Name] == 0 {
			continue
		}

		// Return files that need deleting from this bucket based on the Retention Policy
		objectList := S3Objects{}

		result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: object.Name})
		if err != nil {
			log.Fatal(err)
		}

		// The date we need to delete items prior to
		comparisonDate := time.Now().AddDate(0, 0, -retention[*object.Name]+1)

		for _, object := range result.Contents {
			//fmt.Printf("Object Name: %s Object Modifiled Date: %s\n", *object.Key, object.LastModified)

			if object.LastModified.Before(comparisonDate) {
				objectList.Items = append(objectList.Items, S3Object{
					Key:          *object.Key,
					LastModified: *object.LastModified,
				})
				fmt.Printf("Object Name: %s Object Modifiled Date: %s\n", *object.Key, object.LastModified)
			}
		}

		fmt.Printf("Bucket Name: %s Retention Policy: %d\n", *object.Name, retention[*object.Name])
	}

	//output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
	//	Bucket: aws.String("sozo-db-backups"),
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println("First page results:")
	//for _, object := range output.Contents {
	//	log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	//}
}
