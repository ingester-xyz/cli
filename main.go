// main.go
package main

import (
	"fmt"

	"github.com/ingester-xyz/cli/cmd"
)

func main() {
	// This is the only code in main.go!
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
