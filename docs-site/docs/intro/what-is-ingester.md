# What is Ingester?

**Ingester** is a Go-based CLI tool that migrates existing data from AWS S3 to Walrus, a decentralized object store.

It solves the complexity of:

- Efficiently ingesting complete S3 buckets
- Persisting metadata for lookup
- Fetching and listing S3-origin files by key

With Ingester, you gain CLI control over storage workflows previously locked into centralized systems.
