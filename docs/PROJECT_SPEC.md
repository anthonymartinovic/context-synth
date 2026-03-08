# Context Synth -- Project Spec

## Summary

**Context Synth** is a weighted epistemic runtime for AI reasoning over repository-relevant knowledge.

In plain English: **AI frames probability as certainty. To address that, Context Synth governs what shapes the answer, and by how much.**

It compiles repository-relevant knowledge into a governed, reasoning-ready context environment for humans, agents, and tools.

Its purpose is not to help models access more information in the abstract, but to ensure that downstream reasoning operates over a bounded, reviewable, source-backed knowledge environment.

It accepts knowledge from local and external sources, including repository documents, file systems, and MCP-fetched systems such as Notion, Confluence, Jira, or other machine-readable knowledge stores. It normalizes those inputs, applies explicit source influence, and produces governed outputs that can be inspected, versioned, and used for grounded reasoning.

In this framing, repository-relevant does not mean repository-hosted. It means materially relevant to work performed in or around a repository.

## Problem

Modern software work depends on more than code. Critical context often lives across:

- architecture documents
- requirements
- ADRs
- tickets
- domain references
- external knowledge systems
- team memory
- operational constraints

Tools already make this information accessible. The failure is not access. The failure is that access does not produce a governed reasoning environment.

As a result:

- contributors repeatedly ask the same background questions
- agents and tools operate over incomplete or loosely assembled context
- decisions depend on knowledge that is not visible at the point of work
- outputs are difficult to review against source material
- plausible but unsupported conclusions are treated as factual

The deeper issue is epistemic, not merely ergonomic. When the reasoning environment is unbounded, unstructured, or unaudited, consumers fill gaps with inference. For AI systems, this often appears as hallucination. For teams, it appears as inconsistent judgment, re-explanation, and low trust.

## Project Thesis

Context Synth introduces a governance layer between raw knowledge and downstream reasoning.

Instead of:

Raw sources -> model, tool, or human -> conclusion

Context Synth establishes:

Sources -> governed context assembly -> reasoning-ready snapshot -> grounded synthesis

The project exists to answer a stricter question than "what information is available?":

**What knowledge is allowed to influence this reasoning process, in what form, under what constraints, and with what traceability?**

## Project Definition

Context Synth is a project for constructing governed context from multiple sources and making that context available for reliable downstream use.

It produces two closely related outputs:

### 1. Governed Context Artifact

A bounded, inspectable, source-backed artifact representing the working set of context selected for a given use.

This artifact is intended to be:

- reviewable by humans
- versionable
- auditable
- portable across tools and models
- derived rather than canonical

### 2. Reasoning-Ready Snapshot

A structured context state suitable for grounded synthesis, assessment, and conclusion generation.

This snapshot is intended to support:

- weighted synthesis across multiple sources
- traceable claims
- conflict-aware interpretation
- bounded reasoning
- reproducibility across runs and models

The artifact and the snapshot are related but not identical. The artifact is the human-visible governed output. The snapshot is the machine-usable epistemic state derived from governed inputs.

## What Makes Context Synth Distinct

Context Synth is not primarily a search layer, a note-taking system, or a generic RAG implementation.

Its distinguishing idea is that source influence is explicit and governed.

Users do not merely select sources. They declare how those sources should matter. Context Synth then preserves and operationalizes that influence through assembly and synthesis.

This enables a different kind of system:

- not "retrieve what seems relevant"
- not "summarize everything available"
- not "let the model decide what matters"

Instead:

- select relevant knowledge sources
- assign epistemic influence
- normalize and bound the working set
- produce outputs whose claims can be traced to inputs

## Core Inputs

Context Synth is controlled through a small number of user-declared inputs.

### Sources

A source is any knowledge input relevant to the task, repository, or decision space.

Sources may include:

- local files
- repository documents
- directories or globbed file sets
- remote content snapshots
- MCP-fetched resources
- external knowledge systems
- structured or semi-structured documents

The project must treat external MCP-backed sources as first-class inputs, not as future edge cases or integrations bolted onto a repository-only model.

### Weight

Each source carries a weight representing its epistemic influence.

Weight is not just ordering. It is a declared signal about how much a source should shape the resulting context and synthesis.

By default, weight represents a blend of:

- authority: how trustworthy or canonical the source is
- relevance: how pertinent the source is to the current task or domain
- priority: how strongly the user wants the source to influence the result

Weight should not imply absolute truth. It is a governance signal, not a guarantee.

### Budget

Budget defines how much context is allowed into the governed output and reasoning state.

Budget exists to preserve intentionality. Without boundedness, context accumulation becomes another form of noise.

### Structure

The user may also define the desired output structure, such as sections, templates, categories, or expected synthesis shape.

Structure matters because governed context is not merely collected. It is shaped for use.

## Core Outputs

Context Synth should support outputs such as:

- context artifacts
- reasoning snapshots
- weighted assessments
- grounded summaries
- synthesis memos
- decision-support outputs
- conflict-aware conclusions

The long-term goal is not merely a better context document. It is a governed substrate for trustworthy synthesis.

## Source Weight Semantics

Weight should be treated as epistemic influence that propagates through the system.

At a minimum, that means weight affects:

- prioritization during context selection
- inclusion under budget constraints
- conflict resolution behavior
- ordering and prominence in outputs
- synthesis emphasis
- explanation of why a conclusion leaned in one direction

A higher-weight source should generally have more influence on the final governed result than a lower-weight source, all else equal.

However, weight must not override source traceability, provenance, or reviewability. A highly weighted source may dominate a conclusion, but that dominance must remain visible and inspectable.

## Project Principles

### Governed, Not Generative-First

Generated output is never self-justifying. Source alignment must remain visible.

### Bounded By Design

Useful context is constrained context. Budgeting is a core project feature, not an implementation detail.

### Source-Backed By Default

Outputs must preserve provenance and maintain clear links to the sources that shaped them.

### Reviewable By Humans

A reviewer must be able to inspect what was included, what was omitted, and why the result looks the way it does.

### External-Source Native

The project must assume that important context may live outside the repository and be fetched through MCP or similar interfaces.

### Reproducible Enough To Trust

The same governed inputs should yield materially stable outputs, especially in deterministic or semi-deterministic modes.

### Derived, Not Canonical

Context Synth does not replace source systems. It produces governed working sets and reasoning states from canonical inputs.

## What Context Synth Is

- A governance layer for context
- A compiler for repository-relevant knowledge
- A weighted reasoning substrate
- Developer infrastructure for bounded, auditable synthesis
- A bridge between scattered knowledge and grounded downstream use

## What Context Synth Is Not

- Not a new system of record
- Not a wiki
- Not a general note-taking system
- Not a prompt manager
- Not a generic RAG layer over arbitrary data
- Not an agent framework
- Not a model provider
- Not a substitute for canonical documentation systems

## Primary Use Cases

### Repository Context Compilation

A team compiles architecture, requirements, constraints, and external references into a bounded working artifact that travels with the repository.

### Grounded Agent Workflows

An agent receives a governed reasoning-ready snapshot instead of unconstrained raw context and can generate outputs with better traceability and lower hallucination risk.

### Weighted Multi-Source Assessment

A user supplies multiple sources with explicit influence signals and asks for a grounded conclusion that reflects those declared priorities.

### Review and Decision Support

A reviewer or lead uses Context Synth to surface the most influential evidence behind a recommendation, change, or implementation path.

### Knowledge Portability Across Models and Tools

The same governed context state can be consumed by different models, tools, or workflows without redefining the knowledge environment each time.

## Project Experience

The expected project experience is:

1. Define sources.
2. Assign source influence.
3. Declare structure and budget.
4. Produce governed context.
5. Produce or enable grounded synthesis from that governed context.
6. Inspect provenance, omissions, and influence.
7. Reuse the result across tools, runs, and consumers.

The key shift is that the user governs the reasoning environment before downstream consumption begins.

## Key Project Claims

Context Synth should prove the following claims:

1. Multiple knowledge sources can be turned into a bounded, governed working set.
2. Repository-relevant knowledge can be compiled even when the most important sources live outside the repository.
3. User-declared source influence can materially shape the output in transparent ways.
4. Downstream reasoning can be more trustworthy when performed over a governed context state.
5. Teams can review not only outputs, but the context basis that produced them.
6. The same governed context can serve humans, agents, and tools.

## Success Criteria

Context Synth is successful if:

- teams treat the output as a meaningful working asset
- users can explain why a result looks the way it does
- conclusions are inspectable against source material
- source weighting produces visible, intelligible effects
- external MCP-based knowledge participates as naturally as local files
- consumers trust the governed context more than ad hoc prompting or loose document gathering
- the system reduces context-related ambiguity, repeated explanation, or unsupported synthesis

## Near-Term Project Goal

The first meaningful milestone is not full automation of reasoning. It is proving that governed, weighted, source-backed context is a better substrate for downstream work than ad hoc context collection.

A strong first version should demonstrate:

- mixed-source ingestion, including external MCP-backed inputs
- explicit weight assignment
- bounded artifact creation
- provenance visibility
- reviewable omissions and inclusions
- a grounded synthesis or assessment mode that shows why a conclusion was reached

## Long-Term Direction

Over time, Context Synth can evolve from context compilation into a full reasoning governance layer.

That longer-term direction may include:

- richer source influence models
- conflict-aware synthesis
- claim-level traceability
- reusable reasoning snapshots
- cross-model portability
- policy-aware agent context environments
- institutional knowledge governance for engineering systems

The enduring project idea should remain stable:

**Context Synth governs what knowledge is allowed to shape reasoning, how strongly it shapes it, and how that influence remains visible.**
