// main.go
package main

import (
	"log"
	"os"

	"github.com/ingester-xyz/cli/cmd"
)

func main() {
	// Execute the CLI tool
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
