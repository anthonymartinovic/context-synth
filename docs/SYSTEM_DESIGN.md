# Context Synth -- System Design

> **Note:** This document reflects the v0.1.0 architecture — a two-phase pipeline (context compilation → capability resolution) mediated by a single monolithic artifact. The v1 direction inverts this model: artifacts themselves form the dependency graph, and the application emerges from the live graph rather than from a pipeline endpoint. This document needs to be revised to account for that structural change. See [STATUS.md](../STATUS.md) and [Reflections — v0.1.0](reflections/v0.1.0.md) for the full rationale.

## Purpose

This document describes the technical shape of Context Synth at a system level. It is derived from and subordinate to `docs/PROJECT_SPEC.md`. Details here are expected to evolve as the implementation evolves.

## System Overview

Context Synth is a context-driven orchestration framework. It compiles knowledge from configured sources into a governed context artifact, then uses that artifact to construct and resolve a dependency graph of external system capabilities into dynamic outputs.

The system operates in two phases:

```
Phase 1: Context Compilation

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
   Context Artifact

Phase 2: Capability Resolution

Context Artifact + External Capabilities
        │
        ▼
   Graph Construction
        │
        ▼
   Dependency Resolution
        │
        ▼
   Projection (via adapter)
        │
        ▼
   Dynamic Output
```

## Protocol

The protocol defines the contract surface of Context Synth — the formats, interfaces, and declarations that components communicate through.

### Context Artifact Format

The context artifact is the governed, bounded output of the engine and the primary input to the runtime. Its format is the stable interface between compilation and resolution.

The artifact contains:
- Front matter (generator, mode, snapshot hash, config hash, timestamp, budget, tokens used)
- Sections with source-attributed extractions
- Omission records

The artifact format must be machine-readable and self-describing. The runtime consumes it without knowledge of how it was produced.

### Capability Declaration

External systems declare their capabilities through a standard format. A capability declaration specifies:
- Identity and description
- Inputs required
- Outputs produced
- Dependencies on other capabilities

The framework discovers and maps these declarations into the dependency graph. It does not validate whether the external system can fulfill its declaration — it trusts the declaration and resolves the graph accordingly.

### Adapter Interface

Adapters receive resolved state from the runtime and produce observable output. The adapter interface defines:
- What resolved state looks like (the contract adapters implement against)
- What projection output looks like
- That adapters are stateless with respect to graph resolution — they render, they do not influence

Adapters are interchangeable. The same resolved state projects into different output forms through different adapters.

## Engine

The engine compiles raw context into the governed context artifact. It implements Phase 1 of the pipeline.

### CLI

Entry point. Loads configuration, validates inputs, invokes the pipeline, reports results.

### Config

Parses and validates the configuration file. Responsible for expressing the compilation inputs:

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

## Runtime

The runtime consumes the context artifact and orchestrates external capabilities. It implements Phase 2 of the pipeline.

### Graph Construction

Builds the capability dependency graph from the governed context artifact and available external capability declarations. Context determines which capabilities are relevant and how they are weighted.

The graph is explicit and inspectable. Its construction is deterministic given the same context artifact and capability set.

### Dependency Resolution

Traverses the graph in dependency order and resolves each capability. The result is a unified resolved state.

Resolution respects governance propagation — the weight and provenance from the context artifact influence how capabilities are resolved. A resolved capability traces back to the governed context that activated it.

### Projection

Renders resolved state through an adapter into observable output. The runtime selects the configured adapter and passes it the resolved state. The adapter produces the output form; the runtime does not know or care what form that takes.

## SDK

The SDK is the developer-facing surface for integrating external systems with Context Synth.

It provides:
- Capability declaration — how external systems describe what they can do
- Adapter implementation — how developers build new projection adapters
- Framework integration — how developers interact with the engine and runtime programmatically

The SDK defines contracts, not implementations. It exposes the protocol formats and provides the interfaces that external systems and adapters implement against.

## Data Model

### Compilation

- **Source** — a knowledge input with a declared weight
- **Snapshot** — normalized, hashed set of sources for a run
- **Extraction** — a candidate item with inherited weight and source reference
- **SynthPlan** — ranked extractions grouped by section
- **ContextArtifact** — the governed, bounded, source-backed compilation

### Resolution

- **Capability** — a declared external system capability with inputs, outputs, and dependencies
- **CapabilityGraph** — the dependency graph of capabilities constructed from governed context
- **ResolvedState** — the unified output of graph resolution
- **Projection** — the observable output rendered by an adapter

## Operating Modes

### Engine Modes

**Full Synthesis** — uses the complete pipeline including LLM-backed extraction, classification, and verification.

**Deterministic Fallback** — skips LLM-backed stages. Assembles source contents directly in weight order within budget. Core governance properties — boundedness, provenance, reviewability, omission recording — are preserved.

### Runtime Modes

To be defined as the runtime is implemented. Open questions include whether the runtime operates in single-shot mode (resolve once, project once), continuous mode (re-resolve on context changes), or both.

## Architectural Boundaries

These boundaries should remain stable even if internals change:

- Source/snapshot layer presents a source-type-agnostic interface. New source types are new adapters.
- Deterministic source capture is separate from LLM-backed derivation.
- Assembly is deterministic and inspectable.
- Source tracking survives all pipeline stages and modes.
- The context artifact is a stable interface between engine and runtime. The engine produces it; the runtime consumes it. Neither knows the other's internals.
- The artifact is a derived working set, not a new source of truth.
- Adapters are interchangeable. The same resolved state projects into different output forms without changes to the graph or the engine.
- Capability discovery is external to the framework core. New capabilities are available by connecting new systems, not by changing the framework.

## Open Technical Questions

- What is the best canonical source-reference model (paragraph, span, or other stable unit)?
- How should verification handle items supported by multiple source fragments?
- What caching strategy, if any, should exist across runs?
- Should per-section LLM configuration be supported, allowing users to match model capability to the judgment required for each section's knowledge type?
- What is the serialization format for the context artifact protocol (JSON, protobuf, YAML, or custom)?
- How are capability declarations discovered — static configuration, runtime registration, or both?
- What are the semantics of graph resolution when a capability's dependency cannot be fulfilled?
- Should the runtime support partial resolution (resolve what is possible, report what is not)?
- What does the adapter lifecycle look like — instantiated per-resolution, long-lived, or configurable?
