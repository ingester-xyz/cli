# Architecture Overview

```mermaid
graph TD
  A[S3 Bucket] --> B[Ingester CLI]
  B --> C[Walrus Publisher]
  B --> D[Walrus Aggregator]
  B --> E[Metadata Blob (S3 key → BlobID map)]

  subgraph CLI Features
    B1[Command Line Interface]
    B2[Blob Ref Writer]
    B3[Retriever]
  end

  B --> B1
  B1 --> B2
  B1 --> B3
```
