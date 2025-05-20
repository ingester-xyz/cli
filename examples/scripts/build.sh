#!/bin/bash

echo "-> Building Go binary..."

# Navigate to the Go source directory (ensure correct relative path)
go clean
go build -o ingester ..

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Make the binary executable (in case it's not already set)
chmod +x ./ingester

echo "-> Ingester CLI built successfully."
