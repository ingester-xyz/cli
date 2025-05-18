package cmd

import (
	"fmt"
	"log"

	"github.com/ingester-xyz/cli/pkg/s3"
	"github.com/ingester-xyz/cli/pkg/walrus"
	"github.com/spf13/cobra"
)

// S3Cmd defines the command to ingest data from AWS S3 and store into Walrus,
// then persist the references blob for later lookups.
var S3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Ingest data from AWS S3, store in Walrus, and persist refs metadata",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")

		// Validate required flags
		if bucket == "" {
			log.Fatal("Bucket name is required")
		}
		if region == "" {
			log.Fatal("Region is required")
		}

		// Download from S3
		data, err := s3.IngestFromS3(bucket, region)
		if err != nil {
			log.Fatalf("Error ingesting data from S3: %v", err)
		}

		// Upload to Walrus
		refs, errs := walrus.IngestFiles(data)
		if len(errs) > 0 {
			for key, e := range errs {
				log.Printf("Failed to ingest %s: %v", key, e)
			}
			log.Fatal("One or more files failed to ingest into Walrus")
		}

		// Persist the reference mapping as a metadata blob
		metaBlobID, err := walrus.PersistRefs(refs)
		if err != nil {
			log.Fatalf("Error persisting refs metadata: %v", err)
		}

		// Output the metadata blob ID for later lookup
		fmt.Printf("Refs metadata stored as blob: %s\n", metaBlobID)
	},
}

func init() {
	S3Cmd.Flags().String("bucket", "", "S3 bucket name")
	S3Cmd.Flags().String("prefix", "", "S3 prefix (folder)")
	S3Cmd.Flags().String("manifest-file", "", "Path to manifest file (unused)")
	S3Cmd.Flags().String("region", "", "AWS region")
	S3Cmd.Flags().Bool("encryption", false, "Encrypt before upload (unused)")
	S3Cmd.Flags().String("tags", "", "Comma-separated tags (unused)")
	S3Cmd.Flags().Int("concurrency", 5, "Concurrent downloads (unused)")
}
