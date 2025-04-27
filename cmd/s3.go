package s3

import (
	"fmt"

	"github.com/ingester-xyz/cli/pkg/s3"
	"github.com/spf13/cobra"
)

// S3Cmd represents the S3 command
var S3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Interact with AWS S3",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new AWS S3 ingester
		ingester := s3.NewS3Ingester()
		// Call the function to list files from a bucket (example)
		err := ingester.ListFiles("example-bucket")
		if err != nil {
			fmt.Println("Error interacting with S3:", err)
		}
	},
}
