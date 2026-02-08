# Vision

Context Synth defines context synthesis as a pipeline.

Sources → Normalized Units → Canonical Context Document

During synthesis, normalized units are weighted and assembled to produce a canonical context window appropriate for downstream AI consumption, with explicit consideration for token and budget constraints.

The trajectory is toward a stable, versioned open-source pipeline that:

* Establishes schemas for structured context representation
* Supports adapter-driven source ingestion
* Produces deterministic synthesized documents
* Enables integration into developer tooling

Practically, this is meant to sit between upstream context providers (including context surfaced via MCP servers) and the downstream prompt window: adapters ingest/normalize sources into units, the core pipeline synthesizes, and the resulting canonical document becomes the context you provide to your chat/prompt.

## v1.0 Intent

The pipeline is intended to be delivered through:

* A command-line interface
* Integration with Cursor and VS Code
