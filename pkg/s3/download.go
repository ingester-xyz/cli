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

// DownloadFileFromS3 downloads a file from S3 using an existing S3 client.
func DownloadFileFromS3(s3Client *s3.S3, bucket, key string) ([]byte, error) {
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

// IngestFromS3 lists and downloads all objects from the given S3 bucket using environment credentials.
// It returns a map of object keys to their content.
func IngestFromS3(bucket, region string) (map[string][]byte, error) {
	// Default region to AWS_REGION if not provided
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

	result := make(map[string][]byte, len(listOutput.Contents))
	for _, obj := range listOutput.Contents {
		content, err := DownloadFileFromS3(s3Client, bucket, *obj.Key)
		if err != nil {
			log.Printf("error downloading %s: %v", *obj.Key, err)
			continue
		}
		result[*obj.Key] = content
	}

	return result, nil
}
