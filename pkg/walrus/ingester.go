package walrus

import (
	"log"
	"os"
	"strings"

	walrusSDK "github.com/namihq/walrus-go"
)

// IngestFiles uploads multiple files to Walrus using the walrus-go SDK.
// It supports environment variables for overriding endpoints (WALRUS_ENDPOINT,
// WALRUS_AGGREGATOR_URLS, WALRUS_PUBLISHER_URLS). Returns two maps:
//   - refs: mapping S3 keys to their corresponding Walrus blob IDs
//   - errs: mapping S3 keys to any errors encountered during upload
func IngestFiles(files map[string][]byte) (refs map[string]string, errs map[string]error) {
	errs = make(map[string]error, len(files))
	refs = make(map[string]string, len(files))

	// Determine endpoint URLs from env vars
	var aggURLs, pubURLs []string
	if v := os.Getenv("WALRUS_ENDPOINT"); v != "" {
		urls := strings.Split(v, ",")
		aggURLs = urls
		pubURLs = urls
	} else {
		if v := os.Getenv("WALRUS_AGGREGATOR_URLS"); v != "" {
			aggURLs = strings.Split(v, ",")
		}
		if v := os.Getenv("WALRUS_PUBLISHER_URLS"); v != "" {
			pubURLs = strings.Split(v, ",")
		}
	}

	// Build client options based on provided URLs
	opts := []walrusSDK.ClientOption{}
	if len(aggURLs) > 0 {
		opts = append(opts, walrusSDK.WithAggregatorURLs(aggURLs))
	}
	if len(pubURLs) > 0 {
		opts = append(opts, walrusSDK.WithPublisherURLs(pubURLs))
	}

	// Initialize the SDK client
	client := walrusSDK.NewClient(opts...)

	// Upload each file and collect blob references
	for key, content := range files {
		resp, err := client.Store(content, &walrusSDK.StoreOptions{Epochs: 1})
		if err != nil {
			errs[key] = err
			continue
		}

		// Extract blob ID (newly created or existing)
		var blobID string
		if resp.NewlyCreated != nil {
			blobID = resp.NewlyCreated.BlobObject.BlobID
			log.Printf("Uploaded %s as new blob %s", key, blobID)
		} else if resp.AlreadyCertified != nil {
			blobID = resp.AlreadyCertified.BlobID
			log.Printf("File %s already stored as blob %s", key, blobID)
		}

		// Store reference for lookup by original S3 key
		refs[key] = blobID
	}

	return refs, errs
}
