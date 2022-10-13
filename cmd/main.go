package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"time"
	"wasabiCleanup/internal/client/wasabi"
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

	log.Println("Working...")
	for _, object := range buckets.Buckets {
		if retention[*object.Name] == 0 {
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
		comparisonDate := time.Now().AddDate(0, 0, -retention[*object.Name]-1)

		// Iterate through the S3 object pages, printing each object returned.
		var i int
		for p.HasMorePages() {
			i++

			// Next Page takes a new context for each page retrieval. This is where
			// you could add timeouts or deadlines.
			page, err := p.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("failed to get page %v, %v", i, err)
			}

			// Log the objects found
			for _, obj := range page.Contents {
				if obj.LastModified.Before(comparisonDate) {
					objectList.Items = append(objectList.Items, types.ObjectIdentifier{
						Key: obj.Key,
					})
					objectList.Size += obj.Size
					//fmt.Printf("Object Name: %s Object Modified Date: %s\n", *obj.Key, obj.LastModified)
				} else {
					safeList.Items = append(safeList.Items, types.ObjectIdentifier{
						Key: obj.Key,
					})
					safeList.Size += obj.Size
				}
			}
		}

		// Delete the objects that match
		fmt.Printf("Deleting %d objects totallying %s from the %s Bucket\n", len(objectList.Items), ByteCountSI(objectList.Size), *object.Name)
		fmt.Printf("Remaining %d objects totallying %s in the %s Bucket\n", len(safeList.Items), ByteCountSI(safeList.Size), *object.Name)
		//bar := progressbar.Default(int64(len(objectList.Items)))

		//for x := 0; x < len(objectList.Items); x++ {
		del := &types.Delete{
			Objects: objectList.Items,
			Quiet:   true,
		}
		//fmt.Println(del)
		input := s3.DeleteObjectsInput{
			Bucket: object.Name,
			Delete: del,
		}
		//fmt.Println(input)
		output, err := client.DeleteObjects(context.Background(), &input)
		if err != nil {
			panic(err)
		}
		//bar.Add(1)
		//}
		fmt.Println(output)
	}
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
