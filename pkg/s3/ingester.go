package s3

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DownloadFileFromS3 downloads a file from S3 and returns the file content in memory.
func DownloadFileFromS3(bucket string, key string, region string, profile string) ([]byte, error) {
	// Load the AWS session using the provided profile and region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),                            // Use the region passed as a parameter
		Credentials: credentials.NewSharedCredentials("", profile), // Use the profile from AWS credentials
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create AWS session: %v", err)
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Set up the parameters for the S3 GetObject call
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Make the S3 GetObject call to retrieve the object
	resp, err := s3Client.GetObject(params)
	if err != nil {
		return nil, fmt.Errorf("unable to get object %s from bucket %s: %v", key, bucket, err)
	}
	defer resp.Body.Close()

	// Read the file content into memory (as a byte slice)
	fileContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file content from S3: %v", err)
	}

	// Return the file content
	return fileContent, nil
}

// IngestFromS3 will list and download each object from S3 into memory.
func IngestFromS3(bucket string, region string, profile string) (map[string][]byte, error) {
	// Initialize AWS session with the provided profile and region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),                            // Use the region passed as a parameter
		Credentials: credentials.NewSharedCredentials("", profile), // Use the profile from AWS credentials
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create AWS session: %v", err)
	}

	// Create the S3 client
	s3Client := s3.New(sess)

	// Define input parameters for listing objects
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	// List objects in the specified bucket and prefix
	result, err := s3Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("unable to list objects in S3: %v", err)
	}

	// Map to store the ingested data (file names and their content)
	ingestedData := make(map[string][]byte)

	// Download each object from S3 and store in memory
	for _, item := range result.Contents {
		// Download the file from S3
		fileContent, err := DownloadFileFromS3(bucket, *item.Key, region, profile)
		if err != nil {
			log.Printf("Error downloading file %s: %v\n", *item.Key, err)
			continue // Proceed to next file if there's an error
		}

		// Store the downloaded file content in the map
		ingestedData[*item.Key] = fileContent
	}

	// Return the map of ingested data
	return ingestedData, nil
}
