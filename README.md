# CLI Tool for S3 ↔ Walrus

A Golang-based command‑line application to:

- **Ingest** objects from an AWS S3 bucket into Walrus storage
- **Persist** a metadata blob mapping original S3 keys to Walrus blob IDs
- **List** ingested S3 keys via the metadata blob
- **Fetch** (get) ingested files by their original S3 key

---

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

4. (Optional) **S3 emulator** for local tests (e.g. MinIO, LocalStack) if you don’t want to hit real AWS.

---

## Building the CLI

```bash
# Clone the repository
git clone <repo-url> cli-tool
cd cli-tool

# Download dependencies
go mod tidy

# Build the binary
go build -o cli ./cmd
```

This generates an executable named `cli` in your project root.

---

## Commands Overview

All commands live under the `cli` root command:

| Command  | Description                                                                            |
| -------- | -------------------------------------------------------------------------------------- |
| `s3`     | Ingests all objects from a bucket into Walrus and persists refs metadata               |
| `list`   | Lists all ingested S3 keys from a metadata blob                                        |
| `get`    | Retrieves a single file by its S3 key and writes to stdout or file                     |
| `lookup` | (alias) Same as `list` + `get` combined: lists when no key, fetches when `--key` given |

### 1. `cli s3`

Downloads every object from an S3 bucket, uploads each to Walrus, then writes a single metadata blob containing the S3-to-BlobID map.

```bash
cli s3 \
  --bucket my-bucket       \
  --region us-west-2
# → Refs metadata stored as blob: QmSvz…Yz123
```

- **Flags**:

  - `--bucket` (string, required)
  - `--region` (string, required)
  - Other flags (`--prefix`, `--tags`, etc.) are reserved for future use and currently ignored.

### 2. `cli list`

Prints all original S3 keys stored in the given metadata blob.

```bash
cli list --meta-blob-id QmSvz…Yz123
```

- **Flags**:

  - `--meta-blob-id` (string, required)

### 3. `cli get`

Fetches one ingested file by its original S3 key and writes the raw bytes to stdout or a file.

```bash
# To stdout:
cli get \
  --meta-blob-id QmSvz…Yz123 \
  --key path/to/file.txt

# To a file:
cli get \
  --meta-blob-id QmSvz…Yz123 \
  --key images/photo.png \
```

- **Flags**:

  - `--meta-blob-id` (string, required)
  - `--key` (string, required)

### 4. `cli lookup`

Combined behavior of `list` and `get`:

- No `--key` → lists all keys
- With `--key` → fetches that file to stdout

```bash
# list mode:
cli lookup --meta-blob-id QmSvz…Yz123

# get mode:
cli lookup --meta-blob-id QmSvz…Yz123 --key path/to/data.json
```

---

## Example Workflow

1. **Ingest** an S3 bucket into Walrus:

   ```bash
   export AWS_ACCESS_KEY=<you-aws-access-key>
   export AWS_SECRET_ACCESS_KEY=<you-aws-secret-key>
   export AWS_REGION=eu-west-1
   export WALRUS_ENDPOINT="https://aggregator.walrus-testnet.walrus.space,https://publisher.walrus-testnet.walrus.space"

   ./cli s3 --bucket my-test-bucket --region eu-west-1
   # → Refs metadata stored as blob: QmSvz…Yz123
   ```

2. **List** all ingested keys:

   ```bash
   ./cli list --meta-blob-id QmSvz…Yz123
   ```

3. **Retrieve** one file:

   ```bash
   ./cli get --meta-blob-id QmSvz…Yz123 --key path/to/file.txt
   ```

4. **Lookup** in combined mode:

   ```bash
   # list
   ./cli lookup --meta-blob-id QmSvz…Yz123

   # get
   ./cli lookup --meta-blob-id QmSvz…Yz123 --key path/to/file.txt
   ```

---

## Testing against AWS S3 in eu-west-1

Use a real AWS S3 bucket in the `eu-west-1` region for end‑to‑end testing:

1. **Create a test bucket**:

   ```bash
   aws s3 mb s3://my-test-bucket --region eu-west-1
   ```

2. **Upload a sample file**:

   ```bash
   aws s3 cp ./example.txt s3://my-test-bucket/example.txt --region eu-west-1
   ```

3. **Run the ingest command**:

   ```bash
   ./cli s3 --bucket my-test-bucket --region eu-west-1
   # → Refs metadata stored as blob: QmSvz…Yz123
   ```

4. **List and retrieve** using the metadata blob ID as usual:

   ```bash
   ./cli list --meta-blob-id QmSvz…Yz123
   ./cli get --meta-blob-id QmSvz…Yz123 --key example.txt
   ```

---

Happy ingesting on AWS S3 in eu-west-1! Feel free to open issues or contribute enhancements on GitHub.
