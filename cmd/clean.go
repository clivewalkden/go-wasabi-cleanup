/*
Copyright © 2022 Clive Walkden <clivewalkden@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
	"log"
	"time"
	"wasabiCleanup/internal/client/wasabi"
	"wasabiCleanup/internal/config"
	"wasabiCleanup/internal/reporting"
	"wasabiCleanup/internal/utils"

	"github.com/spf13/cobra"
)

var (
	dryRun bool

	// cleanCmd represents the clean command
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean up the outdated files.",
		Run: func(cmd *cobra.Command, args []string) {
			clean(cmd)
		},
	}
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

func init() {
	cleanCmd.Flags().BoolVarP(&dryRun, "dryrun", "n", false, "Show what will be deleted but don't delete it")
}

func clean(cmd *cobra.Command) {
	dryRun, _ := cmd.Flags().GetBool("dryrun")
	verbose, _ := cmd.Flags().GetBool("verbose")

	client := wasabi.Client()

	report := reporting.Report{DryRun: dryRun}

	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Working...")
	for _, object := range buckets.Buckets {
		if verbose {
			fmt.Printf("Checking Bucket %s\n", *object.Name)
		}

		if config.AppConfig().Buckets[*object.Name] == 0 {
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
		comparisonDate := time.Now().AddDate(0, 0, -config.AppConfig().Buckets[*object.Name]-1)
		if verbose {
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

			if verbose {
				fmt.Printf("\t\t- Next page (%d)\n", i)
			}

			// Log the objects found
			for _, obj := range page.Contents {
				if obj.LastModified.Before(comparisonDate) {
					objectList.Items = append(objectList.Items, types.ObjectIdentifier{
						Key: obj.Key,
					})
					objectList.Size += obj.Size

					if dryRun {
						if verbose {
							fmt.Printf("\t\t\t- Deleting object %s\n", *obj.Key)
						} else {
							fmt.Printf("\t- Deleting object %s\n", *obj.Key)
						}
					} else {
						if verbose {
							fmt.Printf("\t\t\t- Deleting object %s\n", *obj.Key)
						}
						_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
							Bucket: object.Name,
							Key:    obj.Key,
						})

						if err != nil {
							panic("Couldn't delete items")
						}
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
