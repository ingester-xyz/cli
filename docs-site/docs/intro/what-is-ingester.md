# What is Ingester?

**Ingester** is a Go-based CLI tool that makes ingesting data into Walrus seamless and predictable.

> **Just like opening the tap 🚰**

In its first iteration, it focuses on migrating existing data from _AWS S3_ into Walrus. Along the way, it generates a metadata file that enables efficient search operations (e.g., `list` or `get`), similar to the native behavior of AWS S3. This simplifies the management of large volumes of blobs stored in Walrus.

![Ingester AWS Logo](../assets/aws-s3-ingester.png){ style="max-width: 20%; height: auto;" }

In the future, more integrations are ready to come (e.g. Arweave -> Walrus).

## Context

Today, Walrus stores data as scattered blobs, with no built-in structure, and ingesting data into it is a manual, fragmented process.
As a result, every team has to build their own solution.

Ingester solves these two critical pain points by:

- Adding an automatic, schema-driven structure that organizes and relates blobs predictably.

- Providing ready-to-use integrations (like AWS S3) so teams can easily ingest data without extra development effort.

With Ingester, Walrus transforms from a raw storage system into a plug-and-play, developer-friendly platform.
Therefore accelerating adoption, reducing barriers, and unlocking its full potential across real-world use cases.
