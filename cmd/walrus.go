package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ingester-xyz/cli/pkg/walrus"
	"github.com/spf13/cobra"
)

// LookupCmd provides commands to search and retrieve objects by their original S3 key
// using the persisted Walrus metadata blob.
var LookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup and retrieve objects stored in Walrus by S3 key",
	Run: func(cmd *cobra.Command, args []string) {
		metaBlobID, _ := cmd.Flags().GetString("meta-blob-id")
		key, _ := cmd.Flags().GetString("key")

		// Validate required flag
		if metaBlobID == "" {
			log.Fatal("--meta-blob-id is required")
		}

		// List all keys if none specified
		if key == "" {
			keys, err := walrus.ListKeys(metaBlobID)
			if err != nil {
				log.Fatalf("error listing keys: %v", err)
			}
			for _, k := range keys {
				fmt.Println(k)
			}
			return
		}

		// Retrieve specific object
		data, err := walrus.GetBlob(metaBlobID, key)
		if err != nil {
			log.Fatalf("error retrieving key %s: %v", key, err)
		}
		// Write the blob content to stdout
		if _, err := os.Stdout.Write(data); err != nil {
			log.Fatalf("error writing output: %v", err)
		}
	},
}

func init() {
	// Flags for lookup command
	LookupCmd.Flags().String("meta-blob-id", "", "Walrus metadata blob ID containing S3 references")
	LookupCmd.Flags().String("key", "", "S3 key to look up; omit to list all keys")

	// Register with root command
	RootCmd.AddCommand(LookupCmd)
}
