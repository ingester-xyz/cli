package walrus

import (
	"encoding/json"
	"fmt"
)

// LoadRefs retrieves and parses a metadata blob containing the map of S3 keys to Walrus blob IDs.
func LoadRefs(metaBlobID string) (map[string]string, error) {
	client := newClientFromEnv()

	// Read the metadata blob content
	data, err := client.Read(metaBlobID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read refs blob %s: %v", metaBlobID, err)
	}

	// Unmarshal JSON into map
	var refs map[string]string
	if err := json.Unmarshal(data, &refs); err != nil {
		return nil, fmt.Errorf("failed to decode refs JSON: %v", err)
	}
	return refs, nil
}

// ListKeys returns all S3 keys stored in the metadata blob
func ListKeys(metaBlobID string) ([]string, error) {
	refs, err := LoadRefs(metaBlobID)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(refs))
	for k := range refs {
		keys = append(keys, k)
	}
	return keys, nil
}

// Lookup returns the Walrus blob ID for a given S3 key
func Lookup(metaBlobID, key string) (string, error) {
	refs, err := LoadRefs(metaBlobID)
	if err != nil {
		return "", err
	}
	blobID, ok := refs[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in refs", key)
	}
	return blobID, nil
}

// GetBlob retrieves the raw content for a given S3 key by fetching its referenced blob
func GetBlob(metaBlobID, key string) ([]byte, error) {
	blobID, err := Lookup(metaBlobID, key)
	if err != nil {
		return nil, err
	}
	client := newClientFromEnv()

	// Read the blob content
	data, err := client.Read(blobID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob %s: %v", blobID, err)
	}
	return data, nil
}
