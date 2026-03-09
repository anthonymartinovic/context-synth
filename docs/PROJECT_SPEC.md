# Context Synth -- Project Spec

## Summary

**Context Synth** is a weighted epistemic runtime for governed AI reasoning.

In plain English: **AI frames probability as certainty. This is not a bug at the edge of the system - it is a consequence of how current AI systems are designed to answer. To account for this, Context Synth governs what shapes the answer, and by how much.**

It compiles knowledge from diverse sources into a governed, reasoning-ready context environment for humans, agents, and tools.

Its purpose is not to help models access more information in the abstract, but to ensure that downstream reasoning operates over a bounded, reviewable, source-backed knowledge environment.

It accepts knowledge from local and external sources, including repository documents, file systems, and MCP-fetched systems such as Notion, Confluence, Jira, or other machine-readable knowledge stores. It normalizes those inputs, applies explicit source influence, and produces governed outputs that can be inspected, versioned, and used for grounded reasoning.

The scope of "relevant knowledge" is not limited to what lives in a repository. It means any knowledge materially relevant to the reasoning context — whether that context is a codebase, a project, a team, or a decision.

## Problem

Modern software work depends on more than code. The context that shapes decisions lives in architecture docs, spreadsheets, message threads, tickets, intranet pages — spread across systems that were never designed to talk to each other.

Tools already reach most of that information. But access is not the failure. The failure is that access alone does not produce a governed reasoning environment. When context is ungoverned, the same background questions get asked repeatedly, decisions get made against knowledge that wasn't visible at the point of work, and review quality becomes a function of who was in the room when the decision was made.

That problem is real for humans. For LLM systems, it becomes something worse. Ungoverned context leaves gaps, and LLMs do not leave gaps empty — they fill them from probability and present the result as fact. This is usually called hallucination, but that framing is too soft. It is a predictable consequence of systems designed to produce plausible answers regardless of whether the basis for those answers is complete, structured, or verified.

The problem is therefore epistemic, not merely ergonomic. Uncertainty is being converted into assertion inside a reasoning environment that has not been bounded, weighted, or made reviewable. Bounded, weighted context is not a quality improvement on top of this — it is the prerequisite that makes the reasoning trustworthy at all.

## Project Thesis

Context Synth introduces a governance layer between raw knowledge and downstream reasoning.

Instead of:

Raw sources -> model, tool, or human -> conclusion

Context Synth establishes:

Sources -> governed context assembly -> reasoning-ready snapshot -> grounded synthesis

The project exists to answer a stricter question than "what information is available?":

**What knowledge is allowed to influence this reasoning process, in what form, under what constraints, and with what traceability?**

## User Inputs

Four inputs govern a synthesis run:

- **Sources**: knowledge inputs to the pipeline — local files, repository documents, directories, or MCP-fetched content. External sources are first-class.
- **Weight**: each source's declared authority, expressed as a value between 0.0 and 1.0. 1.0 is full authority; 0.0 is none. Weight propagates through the pipeline: it governs extraction prominence, ranking order, and budget allocation during assembly.
- **Structure**: the shape of the output — sections and categories. The system fills the structure the user provides. It does not invent structure. If source material does not map to a declared section, that section is left empty and the user is notified.
- **Budget**: a token bound on the output. Forces the system to choose. Weight governs how it chooses.

## Pipeline

```
Sources → Snap → Extract → Rank → Assemble → Verify → Context Artifact
```

**Snap** — fetch and normalize configured sources into a source-type-agnostic snapshot. Attach provenance references and content hashes. This is the boundary between fetching and processing.

**Extract** — decompose snapshot content into discrete, attributable units of knowledge and classify each unit into the user-provided structure. A single source may yield many items; extractions are not 1:1 with sources. Each item inherits its source's weight and provenance. This is the judgment layer of the pipeline — segmentation and classification are LLM-backed and non-deterministic. The governance stages on either side (Snap before, Rank/Assemble/Verify after) exist to bound and make that judgment reviewable.

**Rank** — order items within sections by weight and relevance. This is where declared source influence becomes operational. When sources share equal weight, declaration order is the tiebreaker — earlier-declared sources take precedence.

**Assemble** — build the artifact within budget. Higher-weight items included first. Omissions recorded, not dropped.

**Verify** — surface what was included, what was omitted, and the source basis for each. Produces a reviewable verification surface — not an automated pass/fail gate. The output of Verify is material for human or agent review.

## Data Model

- **Source** — a canonical knowledge input with a declared weight
- **Snapshot** — normalized, hashed set of sources for a run
- **Extraction** — a candidate context item with inherited weight and source reference
- **SynthPlan** — ranked extractions grouped by section
- **ContextArtifact** — the final bounded, source-backed output

## Invariants

These hold across all modes and versions:

- **Source traceability** — every included item traces to a specific source, reference, and content hash
- **Boundedness** — output stays within the declared budget
- **Human reviewability** — what was included, omitted, and why is always inspectable
- **Evidence-backed synthesis** — generated output is never self-justifying; source alignment remains visible
- **Governance without LLMs** — provenance, boundedness, weight ordering, and reviewability all survive without LLM-backed stages. Synthesis capability does not — semantic extraction and classification require the judgment layer. Without it, the pipeline falls back to assembling source contents directly in config order within budget.
- **Deterministic conflict resolution** — when sources of equal weight conflict, declaration order governs precedence. Conflict resolution is always inspectable, never implicit.
- **Derived, not canonical** — the artifact is a working set, not a replacement for source documents

## Architectural Boundaries

- The source/snapshot layer presents a source-type-agnostic interface. Adding a new source type is a new adapter, not a pipeline change.
- The pipeline has an explicit determinism boundary. Snap is deterministic. Extract is LLM-backed and non-deterministic — this is where the system exercises judgment. Rank, Assemble, and Verify are deterministic: same ranked inputs and budget produce the same output.
- Source tracking survives all pipeline stages and modes.

## What Context Synth Is Not

- Not a new source of truth
- Not a wiki or documentation system
- Not a generic RAG layer — RAG retrieves context at query time to augment a prompt. Context Synth compiles a governed artifact *before* any query. Sources are declared, weighted, and bounded into a reviewable snapshot. RAG answers "what is relevant to this question"; Context Synth answers "what knowledge is allowed to shape reasoning, and how much."
- Not an agent — Context Synth is a runtime. It does not set goals, select actions, or decide what to do next. It executes a declared pipeline over declared inputs. Agents and humans use its output as a foundation for reasoning — it does not reason on their behalf.
- Not a prompt manager
- Not a model provider

## Design Hypotheses

The thesis, inputs, pipeline shape, and invariants above are intended to be durable. The following are hypotheses this project exists to validate. They guide implementation but are expected to evolve as the system is built.

1. Multiple knowledge sources — including sources outside the repository — can be compiled into a bounded, governed working set.
2. User-declared source influence visibly and materially shapes the output.
3. Downstream reasoning over governed context is more trustworthy than reasoning over ad hoc context.
4. Teams can review the source basis, weight distribution, and omissions behind any output.
5. The same governed artifact serves humans, agents, and tools.
6. Governance properties hold without LLM-backed stages.

## Long-Term Direction

*The following are orientations, not planned features or commitments. They describe where the project points, not what it promises.*

Prove that governed, weighted, source-backed context is a better foundation for downstream work than ad hoc context collection.

Mixed-source ingestion including MCP-backed inputs, explicit weight with visible downstream effects, bounded artifact creation with recorded omissions, provenance on every item, and a synthesis mode that shows why a conclusion was reached.

Context compilation evolves into reasoning governance — richer influence models, conflict-aware synthesis, claim-level traceability, reusable reasoning snapshots, cross-model portability, policy-aware agent context environments.

**Context Synth governs what knowledge is allowed to shape reasoning, how strongly it shapes it, and how that influence remains visible. The aim is for it to become an infrastructure layer for constructing, versioning, and supplying context to AI systems.**
