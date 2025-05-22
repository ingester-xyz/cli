## Prerequisites

1. **Go** (v1.18 or later) installed and on your `PATH`.

2. **AWS credentials** (with S3 read access) available via environment variables:

   ```bash
   export AWS_ACCESS_KEY_ID=<your-aws-access-key-id>
   export AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key>
   export AWS_REGION=<your-aws-default-region>
   ```

3. **Walrus endpoints** (public gateway URLs) configured in your shell:

   ```bash
   # Single var for both read/write
   export WALRUS_ENDPOINT="https://aggregator.walrus-testnet.walrus.space,https://publisher.walrus-testnet.walrus.space"

   # Or separate for more control
   export WALRUS_AGGREGATOR_URLS="https://aggregator.walrus-testnet.walrus.space"
   export WALRUS_PUBLISHER_URLS="https://publisher.walrus-testnet.walrus.space"
   ```

---

## Building

```bash
# Clone the repository
git clone git@github.com:ingester-xyz/cli.git
cd cli

# Download dependencies
go mod tidy

# Build the binary
go build -o ingester .
```

This generates an executable named `ingester` in your project root.

---

## Commands Overview

All commands live under the `ingester` root command:

| Command | Description                                                              |
| ------- | ------------------------------------------------------------------------ |
| `s3`    | Ingests all objects from a bucket into Walrus and persists refs metadata |
| `list`  | Lists all ingested S3 keys from a metadata blob                          |
| `get`   | Retrieves a single file by its S3 key and writes to stdout or file       |
| `url`   | Get public URL for AWS S3 ingested data key in Walrus                    |

### 1. `ingester s3`

Downloads every object from an S3 bucket, uploads each to Walrus, then writes a single metadata blob containing the S3-to-BlobID map.

```bash
ingester s3 \
  --bucket my-bucket       \
  --region us-west-2
# → Refs metadata stored as blob: QmSvz…Yz123
```

- **Flags**:

  - `--bucket` (string, required)
  - `--region` (string, required)
  - Other flags (`--prefix`, `--tags`, etc.) are reserved for future use and currently ignored.

### 2. `ingester list`

Prints all original S3 keys stored in the given metadata blob.

```bash
ingester list --meta-blob-id QmSvz…Yz123
```

- **Flags**:

  - `--meta-blob-id` (string, required)

### 3. `ingester get`

Fetches one ingested file by its original S3 key and writes the raw bytes to stdout or a file.

```bash
# To stdout:
ingester get \
  --meta-blob-id QmSvz…Yz123 \
  --key path/to/file.txt

# To a file:
ingester get \
  --meta-blob-id QmSvz…Yz123 \
  --key images/photo.png \
```

- **Flags**:

  - `--meta-blob-id` (string, required)
  - `--key` (string, required)

---

## Example Workflow

1. **Ingest** an S3 bucket into Walrus:

   ```bash
   export AWS_ACCESS_KEY=<you-aws-access-key>
   export AWS_SECRET_ACCESS_KEY=<you-aws-secret-key>
   export AWS_REGION=<you-aws-region>
   export WALRUS_ENDPOINT="https://aggregator.walrus-testnet.walrus.space,https://publisher.walrus-testnet.walrus.space"

   ./ingester s3 --bucket my-test-bucket --region eu-west-1
   # → Refs metadata stored as blob: QmSvz…Yz123
   ```

2. **List** all ingested keys:

   ```bash
   ./ingester list --meta-blob-id QmSvz…Yz123
   ```

3. **Retrieve** one file:

   ```bash
   ./ingester get --meta-blob-id QmSvz…Yz123 --key path/to/file.txt
   ```
