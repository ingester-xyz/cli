# Ingester CLI examples 🚰

These examples provides a set of scripts for ingesting files into Walrus, either from a **local file** or from an **S3 bucket**.

- **`file-ingestion.sh`**: Ingests a local file into Walrus.
- **`s3-aws-ingestion.sh`**: Ingests files from an S3 bucket into Walrus.

Make sure to configure your environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `AWS_S3_BUCKET`) correctly before running the scripts.

The scripts use the `ingester` Go CLI to perform the ingestion process 🚰.

## Folder Structure

The `examples/` folder contains the following scripts:

- **`file-ingestion.sh`**: Ingests a local file into Walrus.
- **`s3-aws-ingestion.sh`**: Ingests files from an S3 bucket into Walrus.
- **`env`**: Contains the AWS environment variables required to interact with the S3 bucket.

## Prerequisites

### 1. **Go Environment**

Ensure Go is installed and configured on your machine. To check if Go is installed, run:

```bash
go version
```

If Go is not installed, follow the [Go installation guide](https://golang.org/doc/install) to install Go.

### 3. **Ingest CLI Binary**

The `ingester` binary must be built using the Go source code before running the ingestion scripts. The `file-ingestion.sh` and `s3-aws-ingestion.sh` scripts will automatically build the binary.

## Setup and Configuration

### 1. **Set Environment Variables**

In the `env` file located in the `examples/` folder, you must configure your AWS credentials and the S3 bucket name.

Example of the `env` file:

```bash
#!/bin/bash

# AWS env variables to provide S3 bucket read+list access
export AWS_ACCESS_KEY_ID="<your-aws-access-key-id>"
export AWS_SECRET_ACCESS_KEY="<your-aws-secret-access-key>"
export AWS_REGION="<your-aws-default-region>"

export AWS_S3_BUCKET="<your-aws-s3-bucket-name-to-ingest-from>"
```

Replace the placeholders (`<your-aws-access-key-id>`, `<your-aws-secret-access-key>`, etc.) with your actual AWS credentials and the S3 bucket name.

## Running the Scripts

### 1. **Ingesting a Local File**

To ingest a local file (`./assets/sample.webp`) into Walrus, run the following script:

```bash
./examples/file-ingestion.sh
```

This script does the following:

1. Sets the necessary environment variables.
2. Checks if the local file exists (`./assets/sample.webp`).
3. Builds the Go `ingester` binary.
4. Executes the ingestion of the local file using the `ingester local file` command.
5. Cleans up after the ingestion.

### 2. **Ingesting Files from an S3 Bucket**

To ingest files from an S3 bucket into Walrus, run the following script:

```bash
./examples/s3-aws-ingestion.sh
```

This script performs the following:

1. Sets the necessary environment variables for AWS access.
2. Checks if the required AWS environment variables are set.
3. Verifies that the AWS CLI is installed and configured.
4. Builds the Go `ingester` binary.
5. Executes the ingestion process using the `ingester s3` command with the configured S3 bucket and region.

### Troubleshooting

If you encounter any errors while running the scripts, check the following:

1. **AWS Environment Variables**: Ensure the AWS environment variables are correctly set in the `env` file.
2. **File Path**: For `file-ingestion.sh`, ensure the file path (`./assets/sample.webp`) exists and is accessible.
3. **Permission Issues**: Make sure your AWS user has the necessary permissions to access the S3 bucket (`s3:ListBucket`, `s3:GetObject`).
