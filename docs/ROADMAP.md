# Context Synth -- Roadmap

## Purpose

This roadmap tracks likely product and implementation direction without turning tentative ideas into design commitments. It should stay shorter and more changeable than `docs/VISION.md` and `docs/design/v0.1.md`.

## Near Term

### Validate the core workflow

Use local markdown files as the initial source type to validate the compilation pipeline without the complexity of external ingestion. Prove that Context Synth can:
- capture configured sources and compile them into a bounded context artifact
- produce extraction and ranking that is meaningfully better than raw concatenation
- preserve source tracking and reviewability throughout the pipeline
- remain useful even when the source set is small and uneven
- maintain a clean source-layer boundary so that adding MCP-fetched sources later is an adapter addition, not a pipeline rework

### Clarify the artifact contract

Resolve the most important design questions:
- what a source reference looks like
- what kind of extracted item is primary
- what verification is actually checking
- what "agent usefulness" means operationally

### Keep the implementation simple

Prefer a narrow, understandable first implementation over a feature-rich one. Avoid locking in provider abstractions, caching strategies, or schema details before the core workflow is validated.

## Later

These are plausible next steps, not promises:

- MCP-fetched source ingestion -- Confluence, Notion, Jira, and other external sources accessed through MCP servers; this is the intended primary ingestion path, not an afterthought
- richer verification semantics
- better source-reference granularity (URIs, content-addressable references for non-file sources)
- stronger artifact metadata and validation tooling
- improved caching and incremental regeneration
- multi-provider or per-stage model support
- artifact diffing and CI-oriented validation
- staleness detection and drift checking against live sources

## Non-Goals Right Now

- proving semantic correctness mathematically
- designing the final stable internal package structure
- committing to a permanent multi-model strategy
- building or maintaining source-system connectors -- MCP servers handle access; Context Synth owns the compilation

## Updating This Document

When priorities change, update this file rather than expanding the release design docs with speculative future work. The goal is to keep roadmap intent visible without scattering tentative plans across multiple files.
