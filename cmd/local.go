package cmd

import (
	"log"

	"github.com/ingester-xyz/cli/pkg/local"
	"github.com/spf13/cobra"
)

// LocalCmd defines the command to ingest data from a local file or folder and store it into Walrus,
// then persist the references blob for later lookups.
var LocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Ingest data from a local file or folder and store it into Walrus",
}

// LocalFileCmd defines the command to ingest a single local file.
var LocalFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Ingest data from a local file and store it into Walrus",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		path, _ := cmd.Flags().GetString("path")

		// Validate required flags
		if path == "" {
			log.Fatal("Path is required")
		}

		// Ingest the file
		err := local.IngestFile(path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// LocalFolderCmd defines the command to ingest all files from a local folder.
var LocalFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Ingest all files from a local folder and store them into Walrus",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		path, _ := cmd.Flags().GetString("path")

		// Validate required flags
		if path == "" {
			log.Fatal("Path is required")
		}

		// Ingest the folder
		err := local.IngestFolder(path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// Register the "file" subcommand under "local"
	LocalCmd.AddCommand(LocalFileCmd)

	// Register the "folder" subcommand under "local"
	LocalCmd.AddCommand(LocalFolderCmd)

	// Define flags for "file" command
	LocalFileCmd.Flags().String("path", "", "Path to the local file to ingest")

	// Define flags for "folder" command
	LocalFolderCmd.Flags().String("path", "", "Path to the local folder to ingest")
}
