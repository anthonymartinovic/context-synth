# Context Synth -- Architecture

## Purpose

This document describes the technical shape of Context Synth at a system level. Unlike `docs/VISION.md`, which defines the enduring product thesis and invariants, and `docs/design/v0.1.md`, which defines the first release target, this file focuses on components, data flow, and implementation boundaries. Details here are expected to evolve as the implementation evolves.

## System Overview

Context Synth is an agent context compiler. It turns configured external knowledge sources into a repo-local context artifact for AI agents. The initial release supports local markdown files to validate the pipeline architecture. The intended primary ingestion path is MCP-fetched content -- Confluence, Notion, Jira, and other external sources accessed through MCP servers. The architecture must not assume local files are the only source type.

```
Configured source documents
        │
        ▼
     Snap sources
        │
        ▼
   Extract candidates
        │
        ▼
     Rank by section
        │
        ▼
    Assemble artifact
        │
        ▼
  Verify source alignment
        │
        ▼
       context.md
```

## Main Components

### CLI Layer

Owns command entry points such as:
- `cs init`
- `cs snap`
- `cs synth`

Responsibilities:
- load configuration
- validate inputs
- invoke the synthesis pipeline
- report progress and results

### Config Layer

Parses and validates `contextsynth.yml`.

Responsibilities:
- resolve source declarations (file paths and globs in v0.1; MCP resources and URIs later)
- apply defaults
- validate section and budget configuration
- expose a normalized runtime config

### Source and Snapshot Layer

Captures configured source material into a deterministic snapshot. This layer owns the boundary between "content has been fetched" and "content enters the pipeline." Downstream stages operate on normalized content with provenance references and do not know or care whether the content originated from a local file, an MCP resource, or another source type.

The v0.1 implementation reads local markdown files. The abstraction boundary should be clean enough that adding MCP-fetched sources later is a new adapter behind the same interface, not a rework of the pipeline.

Responsibilities:
- resolve and fetch source content (local files in v0.1, MCP resources later)
- normalize contents into a source-type-agnostic representation
- attach provenance references (file path, URI, or MCP resource identifier)
- compute content hashes
- estimate token counts
- record the source set used for a synthesis run

### Extraction Layer

Produces candidate context items from the snapshot.

Responsibilities:
- operate over source content
- preserve source references on derived items
- classify or structure extracted context
- optionally cache intermediate outputs

This layer may be LLM-backed, but the exact prompting and item shape are intentionally still flexible.

### Ranking Layer

Orders candidate items within sections.

Responsibilities:
- prioritize for agent usefulness
- work within section-level budgets
- preserve auditable ordering decisions
- optionally cache ranked plans

### Assembly Layer

Builds the final artifact from ranked items.

Responsibilities:
- allocate section budgets
- include or omit items deterministically
- render output markdown
- attach artifact metadata

### Verification Layer

Performs an independent source-alignment check on derived items.

Responsibilities:
- compare included items against referenced source material
- record pass/fail outcomes
- surface excluded items and reasons where applicable

This layer increases confidence in derived context, but it does not replace canonical sources.

## Data Model

The main runtime objects are:

- `Source`: canonical input (local file, MCP resource, or other external content)
- `Snapshot`: normalized set of sources for a run
- `Extraction`: candidate item with source reference
- `SynthesisPlan`: ranked extractions grouped by section
- `ContextArtifact`: final rendered output plus metadata

The exact source-reference model remains open. It may eventually use paragraphs, spans, or another stable unit.

## Data Flow

### 1. Capture

The pipeline begins by resolving configured sources and producing a snapshot tied to specific content hashes.

### 2. Derive

The system derives candidate context items from source material, preserving references back to the snapshot.

### 3. Prioritize

Candidate items are ordered within sections so the most useful context is considered first during assembly.

### 4. Render

The final artifact is assembled within a configured budget and emitted as markdown with metadata.

### 5. Validate

Derived items are checked for source alignment so humans can inspect what was accepted, rejected, or omitted.

## Operating Modes

### Full Synthesis Mode

Uses the full pipeline, including derivation, prioritization, and verification.

### Deterministic Fallback Mode

Skips LLM-backed stages and assembles source contents directly in config order within budget.

This mode is less intelligent, but it preserves the core governance properties:
- bounded artifact generation
- source traceability
- reviewability

## Architectural Boundaries

These boundaries should remain stable even if internals change:

- the source/snapshot layer should present a source-type-agnostic interface to the rest of the pipeline; adding a new source type (e.g. MCP resources) should be a new adapter, not a pipeline change
- deterministic source capture should stay separate from LLM-backed derivation
- assembly should remain deterministic and inspectable
- source tracking should survive all modes
- the artifact should remain a derived working set, not a new source of truth

## Open Technical Questions

- What is the best canonical source-reference model?
- What should the primary extracted item type be?
- How much caching should exist in the initial implementation?
- How provider-specific should the architecture be early on?
- How should verification handle items supported by multiple source fragments?
