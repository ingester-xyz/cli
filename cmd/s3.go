package cmd

import (
	"fmt"
	"log"

	"github.com/ingester-xyz/cli/pkg/s3"
	"github.com/ingester-xyz/cli/pkg/walrus"
	"github.com/spf13/cobra"
)

// S3Cmd defines the command to interact with AWS S3 and ingest into Walrus.
var S3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Ingest data from AWS S3 and store in Walrus",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the flags
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")

		// Validate required flags
		if bucket == "" {
			log.Fatal("Bucket name is required")
		}
		if region == "" {
			log.Fatal("Region is required")
		}

		// Ingest data from S3 using environment credentials
		data, err := s3.IngestFromS3(bucket, region)
		if err != nil {
			log.Fatalf("Error ingesting data from S3: %v", err)
		}

		// Ingest the downloaded files into Walrus blockchain storage
		errors := walrus.IngestFiles(data)
		if len(errors) > 0 {
			// Report any file-specific errors
			for key, err := range errors {
				log.Printf("Failed to ingest %s into Walrus: %v", key, err)
			}
			log.Fatal("One or more files failed to ingest into Walrus")
		}

		fmt.Println("All files successfully ingested into Walrus storage:")
	},
}

func init() {
	// Adding flags to the S3 subcommand
	S3Cmd.Flags().String("bucket", "", "S3 bucket name")
	S3Cmd.Flags().String("prefix", "", "S3 prefix (folder)")
	S3Cmd.Flags().String("manifest-file", "", "Path to the manifest file")
	S3Cmd.Flags().String("region", "", "AWS region")
	// removed profile flag (credentials now via env vars)
	S3Cmd.Flags().Bool("encryption", false, "Encrypt the data before uploading")
	S3Cmd.Flags().String("tags", "", "Comma-separated list of tags")
	S3Cmd.Flags().Int("concurrency", 5, "Number of concurrent uploads")
}
