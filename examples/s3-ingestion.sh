#!/bin/bash

# Set environment variables required for your CLI
export WALRUS_ENDPOINT="https://aggregator.walrus-testnet.walrus.space,https://publisher.walrus-testnet.walrus.space"
echo "-> Environment variables set up"

source ./env

# Step 1: Check if AWS environment variables are set
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ] || [ -z "$AWS_REGION" ] || [ -z "$AWS_S3_BUCKET" ]; then
    echo "Error: One or more required AWS environment variables are missing. Exiting."
    exit 1
fi

echo "-> AWS environment variables set up"

source ./scripts/build.sh

# Step 3: Execute ingestion from local file
echo "-> Executing CLI command to s3 bucket files..."

./ingester s3 --bucket $AWS_S3_BUCKET --region $AWS_REGION

# Check if ingestion was successful
if [ $? -ne 0 ]; then
    echo "Ingestion failed. Exiting."
    exit 1
fi

source ./scripts/clean.sh

echo "-> Ingestion completed."
