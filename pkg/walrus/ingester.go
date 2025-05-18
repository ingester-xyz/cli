package walrus

import (
	"encoding/json"
	"fmt"
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

	client := newClientFromEnv()

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

// PersistRefs serializes a map of S3 keys to Walrus blob IDs and stores
// it as a single metadata blob in Walrus. Returns the blob ID of the metadata
// blob or an error.
func PersistRefs(refs map[string]string) (string, error) {
	// Serialize the refs map to JSON
	data, err := json.Marshal(refs)
	if err != nil {
		return "", fmt.Errorf("failed to serialize refs: %v", err)
	}

	client := newClientFromEnv()

	// Store the JSON blob with a single epoch
	resp, err := client.Store(data, &walrusSDK.StoreOptions{Epochs: 1})
	if err != nil {
		return "", fmt.Errorf("failed to persist refs blob: %v", err)
	}

	// Extract the new metadata blob ID
	var metaBlobID string
	if resp.NewlyCreated != nil {
		metaBlobID = resp.NewlyCreated.BlobObject.BlobID
	} else if resp.AlreadyCertified != nil {
		metaBlobID = resp.AlreadyCertified.BlobID
	}

	log.Printf("Refs metadata stored as blob %s", metaBlobID)
	return metaBlobID, nil
}

// newClientFromEnv constructs a walrus-go SDK client using environment variables
// for endpoint overrides: WALRUS_ENDPOINT, WALRUS_AGGREGATOR_URLS, WALRUS_PUBLISHER_URLS.
func newClientFromEnv() *walrusSDK.Client {
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

	opts := []walrusSDK.ClientOption{}
	if len(aggURLs) > 0 {
		opts = append(opts, walrusSDK.WithAggregatorURLs(aggURLs))
	}
	if len(pubURLs) > 0 {
		opts = append(opts, walrusSDK.WithPublisherURLs(pubURLs))
	}

	return walrusSDK.NewClient(opts...)
}
