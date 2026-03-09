# Context Synth -- Architecture

## Purpose

This document describes the technical shape of Context Synth at a system level. It is derived from and subordinate to `docs/PROJECT_SPEC.md`. Details here are expected to evolve as the implementation evolves.

## System Overview

Context Synth is a governed context runtime. It compiles knowledge from configured sources into a bounded, source-backed context artifact for humans, agents, and tools.

The pipeline accepts four governing inputs — sources, weight, structure, and budget — and produces a reviewable artifact where every included item traces to its origin and every omission is recorded.

```
Sources with weight
        │
        ▼
   Snap (deterministic)
        │
        ▼
   Extract (LLM-backed, non-deterministic)
        │
        ▼
   Rank (deterministic)
        │
        ▼
   Assemble (deterministic)
        │
        ▼
   Verify (review surface)
        │
        ▼
   Contextfile
```

## Components

### CLI

Entry point. Loads configuration, validates inputs, invokes the pipeline, reports results.

### Config

Parses and validates the configuration file. Responsible for expressing the four governing inputs:

- **Sources** with declared weights
- **Structure** as user-defined sections with budget allocations
- **Budget** as a total token bound
- **LLM** configuration (optional)

### Source / Snapshot

Fetches and normalizes configured sources into a source-type-agnostic snapshot. This is the boundary between content acquisition and pipeline processing. Downstream stages do not know or care about source type.

Responsibilities:
- resolve and fetch source content
- normalize into a source-type-agnostic representation
- attach provenance references
- compute content hashes
- estimate token counts

Adding a new source type is a new adapter behind this interface, not a pipeline change.

### Extraction

Decomposes snapshot content into discrete, attributable knowledge items and classifies each into user-provided sections. A single source may yield many items. Each item inherits its source's weight and provenance.

This is the judgment layer. It is LLM-backed and non-deterministic. The governance stages on either side exist to bound and make that judgment reviewable.

### Ranking

Orders items within sections by declared source weight descending. When sources share equal weight, declaration order is the tiebreaker. Deterministic given the same extracted inputs.

### Assembly

Builds the artifact within budget. Higher-weight items are included first. Omissions are recorded, not silently dropped. Deterministic.

### Verification

Produces a reviewable surface showing what was included, what was omitted, and the source basis for each. Not an automated pass/fail gate. The output is material for human or agent review.

## Data Model

- **Source** — a knowledge input with a declared weight
- **Snapshot** — normalized, hashed set of sources for a run
- **Extraction** — a candidate item with inherited weight and source reference
- **SynthPlan** — ranked extractions grouped by section
- **ContextArtifact** — the final bounded, source-backed output

## Operating Modes

### Full Synthesis

Uses the complete pipeline including LLM-backed extraction, ranking, and verification.

### Deterministic Fallback

Skips LLM-backed stages. Assembles source contents directly in weight order within budget. Core governance properties — boundedness, provenance, reviewability, omission recording — are preserved.

## Architectural Boundaries

These boundaries should remain stable even if internals change:

- Source/snapshot layer presents a source-type-agnostic interface. New source types are new adapters.
- Deterministic source capture is separate from LLM-backed derivation.
- Assembly is deterministic and inspectable.
- Source tracking survives all pipeline stages and modes.
- The artifact is a derived working set, not a new source of truth.

## Future Capabilities

- **Per-section LLM configuration.** Sections represent different kinds of knowledge. The architecture should accommodate per-section LLM overrides so users can match model capability to the judgment required for each section's knowledge type. Sections without overrides inherit the global default. Sections with no LLM config fall back to deterministic mode individually.

## Open Technical Questions

- What is the best canonical source-reference model (paragraph, span, or other stable unit)?
- How should verification handle items supported by multiple source fragments?
- What caching strategy, if any, should exist across runs?
