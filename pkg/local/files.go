package local

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ingester-xyz/cli/pkg/walrus"
)

// IngestFile takes a single file and ingests it into Walrus.
func IngestFile(filePath string) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filePath, err)
	}

	// Prepare file data in the format walrus expects (map[string][]byte)
	filesData := map[string][]byte{
		filePath: data,
	}

	// Upload to Walrus
	_, errs := walrus.IngestFiles(filesData)
	if len(errs) > 0 {
		for key, e := range errs {
			log.Printf("Failed to ingest %s: %v", key, e)
		}
		return fmt.Errorf("One or more files failed to ingest into Walrus")
	}

	// Just print the success message, no metadata persistence
	fmt.Printf("File %s successfully ingested into Walrus.\n", filePath)
	return nil
}

// IngestFolder takes a folder path and ingests all files in the folder into Walrus.
func IngestFolder(folderPath string) error {
	// Read all files in the folder
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("Error reading folder %s: %v", folderPath, err)
	}

	// Prepare map to hold file names and their data
	filesData := make(map[string][]byte)

	// Loop through all files in the folder and read their data
	for _, file := range files {
		if !file.IsDir() { // Only process files, skip directories
			filePath := filepath.Join(folderPath, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				continue
			}
			filesData[file.Name()] = data
		}
	}

	// Upload to Walrus
	_, errs := walrus.IngestFiles(filesData)
	if len(errs) > 0 {
		for key, e := range errs {
			log.Printf("Failed to ingest %s: %v", key, e)
		}
		return fmt.Errorf("One or more files failed to ingest into Walrus")
	}

	// Just print the success message, no metadata persistence
	fmt.Printf("Folder %s successfully ingested into Walrus.\n", folderPath)
	return nil
}
