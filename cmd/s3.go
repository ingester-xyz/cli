package cmd

import (
	"fmt"
	"log"

	"github.com/ingester-xyz/cli/pkg/s3"
	"github.com/spf13/cobra"
)

// S3Cmd defines the command to interact with AWS S3.
var S3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Ingest data from AWS S3",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the flags
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")
		profile, _ := cmd.Flags().GetString("profile")

		// Validate required flags
		if bucket == "" {
			log.Fatal("Bucket name is required")
		}
		if region == "" {
			log.Fatal("Region is required")
		}
		if profile == "" {
			log.Fatal("profile is required")
		}

		// Ingest data from S3
		data, err := s3.IngestFromS3(bucket, region, profile)
		if err != nil {
			log.Fatalf("Error ingesting data from S3: %v", err)
		}

		// Here, you could implement logic to save the manifestFile with the ingested data.
		// For now, let's print the result for demonstration purposes.
		fmt.Printf("Ingestion completed. Data from the following files:\n")
		for key, content := range data {
			// You could save content to the manifestFile, for example:
			// err = saveToManifest(manifestFile, key, content)
			fmt.Printf("File: %s, Size: %d bytes\n", key, len(content))
		}
	},
}

func init() {
	// Adding flags to the S3 subcommand
	S3Cmd.Flags().String("bucket", "", "S3 bucket name")
	S3Cmd.Flags().String("prefix", "", "S3 prefix (folder)")
	S3Cmd.Flags().String("manifest-file", "", "Path to the manifest file")
	S3Cmd.Flags().String("region", "", "AWS region")
	S3Cmd.Flags().String("profile", "", "AWS profile to use")
	S3Cmd.Flags().Bool("encryption", false, "Encrypt the data before uploading")
	S3Cmd.Flags().String("tags", "", "Comma-separated list of tags")
	S3Cmd.Flags().Int("concurrency", 5, "Number of concurrent uploads")
}
