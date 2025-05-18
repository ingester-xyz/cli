// main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	walrus "github.com/namihq/walrus-go"
)

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "", "Path to the local file to ingest into Walrus")
	epochs := flag.Int("epochs", 1, "Number of storage epochs to request")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("--file is required (e.g. ./main --file path/to/file.txt)")
	}

	// Read the file into memory
	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Unable to read file %s: %v", *filePath, err)
	}

	// Initialize Walrus client using environment variables:
	//   WALRUS_ENDPOINT or WALRUS_AGGREGATOR_URLS / WALRUS_PUBLISHER_URLS
	client := walrus.NewClient()

	// Store the data with specified epochs
	resp, err := client.Store(data, &walrus.StoreOptions{Epochs: *epochs})
	if err != nil {
		log.Fatalf("Walrus store error: %v", err)
	}

	// Determine the BlobID from the response
	var blobID string
	if resp.NewlyCreated != nil {
		blobID = resp.NewlyCreated.BlobObject.BlobID
		fmt.Printf("🎉 New blob stored: %s\n", blobID)
	} else if resp.AlreadyCertified != nil {
		blobID = resp.AlreadyCertified.BlobID
		fmt.Printf("ℹ️  Blob already exists: %s\n", blobID)
	}

	// default aggregators URL list -> https://github.com/namihq/walrus-go/blob/main/endpoints.go#L41
	aggregatorEndpoint := "https://aggregator.walrus-testnet.walrus.space"

	// Construct and print the public access URL
	accessURL := fmt.Sprintf("%s/v1/blobs/%s\n", aggregatorEndpoint, blobID)
	fmt.Printf("🔗 Access your file at: %s\n", accessURL)

	// open this file with target file type (e.g. ".webp")

	os.Exit(0)
}
