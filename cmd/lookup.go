package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ingester-xyz/cli/pkg/walrus"
	"github.com/spf13/cobra"
)

// ListCmd lists all S3 keys stored in a persisted refs metadata blob.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ingested AWS S3 keys from Walrus refs metadata",
	Run: func(cmd *cobra.Command, args []string) {
		metaBlobID, _ := cmd.Flags().GetString("meta-blob-id")
		if metaBlobID == "" {
			log.Fatal("--meta-blob-id is required to list keys")
		}

		keys, err := walrus.ListKeys(metaBlobID)
		if err != nil {
			log.Fatalf("Error listing keys: %v", err)
		}
		for _, k := range keys {
			fmt.Println(k)
		}
	},
}

// GetCmd retrieves a specific file by its original S3 key and writes content to stdout.
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch ingested file from Walrus by AWS S3 key",
	Run: func(cmd *cobra.Command, args []string) {
		metaBlobID, _ := cmd.Flags().GetString("meta-blob-id")
		key, _ := cmd.Flags().GetString("key")
		outFile, _ := cmd.Flags().GetString("out")

		if metaBlobID == "" || key == "" {
			log.Fatal("Both --meta-blob-id and --key are required to get a file")
		}

		data, err := walrus.GetBlob(metaBlobID, key)
		if err != nil {
			log.Fatalf("Error retrieving key %s: %v", key, err)
		}

		if outFile == "" {
			// write to stdout
			if _, err := os.Stdout.Write(data); err != nil {
				log.Fatalf("Error writing to stdout: %v", err)
			}
		} else {
			// write to specified file
			if err := ioutil.WriteFile(outFile, data, 0644); err != nil {
				log.Fatalf("Error writing to file %s: %v", outFile, err)
			}
			fmt.Printf("Saved %s to %s\n", key, outFile)
		}
	},
}

// UrlCmd returns the public URL for a given S3 key
var UrlCmd = &cobra.Command{
	Use:   "url",
	Short: "Get public URL for AWS S3 ingested data key in Walrus",
	Run: func(cmd *cobra.Command, args []string) {
		metaBlobID, _ := cmd.Flags().GetString("meta-blob-id")
		key, _ := cmd.Flags().GetString("key")

		if metaBlobID == "" || key == "" {
			log.Fatal("Both --meta-blob-id and --key are required to get a blob URL")
		}

		url, err := walrus.GetBlobURL(metaBlobID, key)
		if err != nil {
			log.Fatalf("Error retrieving URL for key %s: %v", key, err)
		}

		fmt.Println(url)
	},
}

func init() {
	// List command flag
	ListCmd.Flags().String("meta-blob-id", "", "Walrus metadata blob ID containing S3 refs")

	// Get command flags
	GetCmd.Flags().String("meta-blob-id", "", "Walrus metadata blob ID containing S3 refs")
	GetCmd.Flags().String("key", "", "Original S3 key to fetch")
	GetCmd.Flags().String("out", "", "Output file path; stdout if empty")

	// Url command flags
	UrlCmd.Flags().String("meta-blob-id", "", "Walrus metadata blob ID containing S3 refs")
	UrlCmd.Flags().String("key", "", "Original S3 key to get public URL for")
}
