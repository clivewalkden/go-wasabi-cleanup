/*
 * Copyright (c) 2023 Clive Walkden <clivewalkden@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
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

// S3Client is an interface that includes the methods we need from s3.Client.
type S3Client interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// S3Object represents an object in an S3 bucket.
type S3Object struct {
	Key          string
	LastModified time.Time
	Size         int64
}

// S3Objects represents a list of objects in an S3 bucket.
type S3Objects struct {
	Items []types.ObjectIdentifier
	Size  int64
}

// init initializes the clean command and its flags.
func init() {
	cleanCmd.Flags().BoolVarP(&dryRun, "dryrun", "n", false, "Show what will be deleted but don't delete it")
}

// GetBuckets is a function that retrieves a list of all buckets from the provided S3 client.
// It returns a slice of Bucket objects and an error. If the operation is successful, the error is nil.
// If there is an error during the operation, the function returns nil and the error.
//
// Parameters:
// client: An instance of an S3 client.
//
// Returns:
// []types.Bucket: A slice of Bucket objects representing all the buckets retrieved from the S3 client.
// error: An error that will be nil if the operation is successful, and an error object if the operation fails.
func GetBuckets(client S3Client) ([]types.Bucket, error) {
	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return buckets.Buckets, nil
}

// ProcessBucket processes a single bucket. It checks if the bucket is in the config and if it needs to be cleaned.
func ProcessBucket(bucket types.Bucket, client S3Client, dryRun bool, verbose bool) (reporting.Result, error) {
	if verbose {
		fmt.Printf("Checking Bucket %s\n", *bucket.Name)
	}

	if config.AppConfig().Buckets[*bucket.Name] == 0 {
		if viper.GetBool("verbose") {
			fmt.Printf("\t- Bucket not in config, skipping\n")
		}
		return reporting.Result{}, nil
	}

	// The date we need to delete items prior to
	comparisonDate := time.Now().AddDate(0, 0, -config.AppConfig().Buckets[*bucket.Name]-1)
	if verbose {
		fmt.Printf("\t- Checking files date is before %s\n", comparisonDate)
	}

	objectList, safeList, err := DeleteOldObjects(bucket, client, comparisonDate, dryRun, verbose)
	if err != nil {
		return reporting.Result{}, err
	}

	result := reporting.Result{
		Name:        *bucket.Name,
		Kept:        len(safeList.Items),
		KeptSize:    utils.ByteCountSI(safeList.Size),
		Deleted:     len(objectList.Items),
		DeletedSize: utils.ByteCountSI(objectList.Size),
	}

	return result, nil
}

// GetObjects retrieves objects from the bucket.
func GetObjects(bucket types.Bucket, client S3Client) ([]types.Object, error) {
	params := &s3.ListObjectsV2Input{Bucket: bucket.Name}
	p := s3.NewListObjectsV2Paginator(client, params)

	var objects []types.Object
	for p.HasMorePages() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		objects = append(objects, page.Contents...)
	}

	return objects, nil
}

// DeleteObject deletes an object from the bucket.
func DeleteObject(bucket types.Bucket, client S3Client, object types.Object) error {
	_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: bucket.Name,
		Key:    object.Key,
	})

	return err
}

// DeleteOldObjects deletes objects in a bucket that are older than the comparison date.
func DeleteOldObjects(bucket types.Bucket, client S3Client, comparisonDate time.Time, dryRun bool, verbose bool) (S3Objects, S3Objects, error) {
	objectList := S3Objects{}
	safeList := S3Objects{}

	objects, err := GetObjects(bucket, client)
	if err != nil {
		return S3Objects{}, S3Objects{}, err
	}

	for _, obj := range objects {
		if obj.LastModified.Before(comparisonDate) {
			objectList.Items = append(objectList.Items, types.ObjectIdentifier{
				Key: obj.Key,
			})
			objectList.Size += aws.ToInt64(obj.Size)

			if dryRun {
				if verbose {
					fmt.Printf("\t\t- Deleting object %s\n", *obj.Key)
				} else {
					fmt.Printf("\t- Deleting object %s\n", *obj.Key)
				}
			} else {
				if verbose {
					fmt.Printf("\t\t- Deleting object %s\n", *obj.Key)
				}
				err = DeleteObject(bucket, client, obj)
				if err != nil {
					return S3Objects{}, S3Objects{}, err
				}
			}
		} else {
			safeList.Items = append(safeList.Items, types.ObjectIdentifier{
				Key: obj.Key,
			})
			safeList.Size += aws.ToInt64(obj.Size)
		}
	}

	return objectList, safeList, nil
}

// CreateReport creates a report based on the results of the cleaning process.
func CreateReport(results []reporting.Result) reporting.Report {
	report := reporting.Report{Result: results}
	return report
}

// clean is the main function for the clean command. It retrieves the list of buckets, processes each bucket, and outputs a report.
func clean(cmd *cobra.Command) {
	dryRun, _ := cmd.Flags().GetBool("dryrun")
	verbose, _ := cmd.Flags().GetBool("verbose")

	client := wasabi.Client()

	buckets, err := GetBuckets(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Working...")
	var results []reporting.Result
	for _, bucket := range buckets {
		result, err := ProcessBucket(bucket, client, dryRun, verbose)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	report := CreateReport(results)
	reporting.Output(report)
}
