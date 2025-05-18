// cmd/root.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd defines the root command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A CLI tool for interacting with AWS S3 and Arweave",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the CLI tool! Use a subcommand.")
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Adding subcommands for S3
	rootCmd.AddCommand(S3Cmd)
	rootCmd.AddCommand(LookupCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(GetCmd)
}
