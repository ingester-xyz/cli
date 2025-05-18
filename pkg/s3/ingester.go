package s3

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DownloadFileFromS3 downloads a file from S3 using AWS environment variables for credentials.
// It takes the bucket name, object key, and optional region. If region is empty, it falls back to AWS_REGION.
func DownloadFileFromS3(bucket, key, region string) ([]byte, error) {
	// Use AWS_REGION environment variable if region not provided
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewEnvCredentials(), // Reads AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, optionally AWS_SESSION_TOKEN
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object %s from bucket %s: %v", key, bucket, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read object body: %v", err)
	}

	return data, nil
}

// IngestFromS3 lists and downloads all objects from the given S3 bucket using AWS environment variables for credentials.
// It returns a map of object keys to their content.
func IngestFromS3(bucket, region string) (map[string][]byte, error) {
	// Use AWS_REGION environment variable if region not provided
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	// List objects in the bucket
	listOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list objects in bucket %s: %v", bucket, err)
	}

	// Download each object
	result := make(map[string][]byte)
	for _, obj := range listOutput.Contents {
		content, err := DownloadFileFromS3(bucket, *obj.Key, region)
		if err != nil {
			log.Printf("error downloading %s: %v", *obj.Key, err)
			continue
		}
		result[*obj.Key] = content
	}

	return result, nil
}
