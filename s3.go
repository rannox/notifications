package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	// Set up a context and a configuration object for the AWS SDK.
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to load SDK config: %v", err))
	}

	// Create a new S3 client using the configuration object.
	client := s3.NewFromConfig(cfg)

	// Set the S3 bucket name and the last modified date you want to filter on.
	bucketName := "your-bucket-name"
	lastModifiedDate := time.Date(2022, 05, 03, 0, 0, 0, 0, time.UTC)

	// Set up a listObjectsV2Input object to list all the objects in the bucket with a last modified date
	// greater than or equal to the one you specified.
	listInput := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}

	// Use the S3 client's ListObjectsV2 method to get a list of all the objects in the bucket that match
	// the last modified date filter.
	var matchingObjects []types.Object
	paginator := s3.NewListObjectsV2Paginator(client, listInput)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to list objects in bucket: %v", err))
		}

		for _, object := range output.Contents {
			if object.LastModified.After(lastModifiedDate) || object.LastModified.Equal(lastModifiedDate) {
				matchingObjects = append(matchingObjects, object)
			}
		}
	}

	// Use the S3 client's GetObject method to get each object's JSON data and extract the desired attribute.
	for _, object := range matchingObjects {
		getInput := &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    object.Key,
		}

		// Use the S3 client's Download method to download the object's data.
		var data []byte
		downloader := manager.NewDownloader(client)
		_, err = downloader.Download(ctx, &data, getInput)
		if err != nil {
			panic(fmt.Sprintf("failed to download object: %v", err))
		}

        // Parse the JSON data and extract the desired attribute.
        // In this example, we're assuming that the JSON data has a top-level "attribute" field.
        // You should replace this with the name of the actual field you want to extract.
        var parsedData map[string]interface{}
        err = json.Unmarshal(data, &parsedData)
        if err != nil {
            panic(fmt.Sprintf("failed to parse JSON data: %v", err))
        }
        attributeValue, ok := parsedData["attribute"].(string)
        if !ok {
            panic("attribute field is not a string")
        }

        // Do something with the extracted attribute value.
        fmt.Printf("Found attribute value '%s' in object '%s'\n", attributeValue, *object.Key)
    }
}