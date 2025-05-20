#!/bin/bash

# Set environment variables required for your CLI
export WALRUS_ENDPOINT="https://aggregator.walrus-testnet.walrus.space,https://publisher.walrus-testnet.walrus.space"
echo "-> Environment variables set up"

source ./scripts/build.sh

# Ensure the assets/sample.webp file exists
if [ ! -f "./assets/sample.webp" ]; then
    echo "Error: File './assets/sample.webp' not found."
    exit 1
fi

# Step 3: Execute ingestion from local file
echo "-> Executing CLI command to ingest local file..."

./ingester local file --path ./assets/sample.webp

# Check if ingestion was successful
if [ $? -ne 0 ]; then
    echo "Ingestion failed. Exiting."
    exit 1
fi

source ./scripts/clean.sh

echo "-> Ingestion completed."
