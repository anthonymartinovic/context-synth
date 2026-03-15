# Context Synth -- Project Spec

## Summary

**Context Synth** is a context-driven orchestration framework for governed AI reasoning.

In plain English: **AI frames probability as certainty. This is not a bug at the edge of the system — it is a consequence of how current AI systems are designed to answer. To account for this, Context Synth governs what shapes the answer, and by how much.**

It synthesizes a dependency graph of system capabilities into dynamic outputs, driven by governed context compiled from diverse knowledge sources.

Its purpose is not to help models access more information in the abstract, but to ensure that downstream reasoning and system behavior operate over a bounded, reviewable, source-backed knowledge environment — and that the capabilities activated by that environment are orchestrated through explicit, traceable dependency relationships.

It accepts knowledge from local and external sources, including repository documents, file systems, and MCP-fetched systems such as Notion, Confluence, Jira, or other machine-readable knowledge stores. It normalizes those inputs, applies explicit source influence, and uses the governed result to drive capability orchestration and output synthesis.

The scope of "relevant knowledge" is not limited to what lives in a repository. It means any knowledge materially relevant to the reasoning context — whether that context is a codebase, a project, a team, or a decision.

## Problem

Modern software work depends on more than code. The context that shapes decisions lives in architecture docs, spreadsheets, message threads, tickets, intranet pages — spread across systems that were never designed to talk to each other.

Tools already reach most of that information. But access is not the failure. The failure is that access alone does not produce a governed reasoning environment. When context is ungoverned, the same background questions get asked repeatedly, decisions get made against knowledge that wasn't visible at the point of work, and review quality becomes a function of who was in the room when the decision was made.

That problem is real for humans. For LLM systems, it becomes something worse. Ungoverned context leaves gaps, and LLMs do not leave gaps empty — they fill them from probability and present the result as fact. This is usually called hallucination, but that framing is too soft. It is a predictable consequence of systems designed to produce plausible answers regardless of whether the basis for those answers is complete, structured, or verified.

The problem is therefore epistemic, not merely ergonomic. Uncertainty is being converted into assertion inside a reasoning environment that has not been bounded, weighted, or made reviewable. Bounded, weighted context is not a quality improvement on top of this — it is the prerequisite that makes the reasoning trustworthy at all.

## Project Thesis

AI is non-deterministic. Context Synth doesn't fight that — it accepts it wholesale. Instead of trying to clamp down AI's reasoning, it creates the structural conditions that incentivize AI to make decisions within a bounded, governed context.

To do this, Context Synth introduces a governance layer between raw knowledge and downstream system behavior.

Instead of:

Raw sources → model, tool, or system → conclusion

Context Synth establishes:

Sources → governed context compilation → capability graph → orchestrated resolution → dynamic output

The project exists to answer a stricter question than "what information is available?":

**What knowledge is allowed to influence this reasoning process, in what form, under what constraints, and with what traceability?**

Context Synth extends this question beyond static documents. Governed context drives not just what knowledge is visible, but which system capabilities are activated, how they depend on each other, and how their synthesis produces output. Governance propagates from source through graph to output.

## Core Abstraction

The dependency graph is the central model.

External systems expose capabilities — things they can do. Context Synth does not define or host these capabilities. It discovers what is available, understands how those capabilities relate and depend on each other, and orchestrates their synthesis into coherent output driven by governed context.

- **Capabilities** are external. They represent what systems can do — generate content, run queries, transform data, produce artifacts. Context Synth treats them as declared nodes in a graph.
- **Dependencies** are relationships between capabilities. One capability may require the output of another. These relationships form a directed graph.
- **Context drives the graph.** The governed context artifact determines which capabilities are relevant, how they are weighted, and how the graph is constructed. Different context produces different graphs.
- **Resolution** traverses the graph and synthesizes a unified state. Dependencies are resolved in dependency order.
- **Projection** renders resolved state into observable output through adapters. Adapters are interchangeable — the same resolved state can project into different output forms without changing the graph or the governance.

## Governing Inputs

### Context Compilation

Four inputs govern how context is compiled:

- **Sources**: knowledge inputs — local files, repository documents, directories, or MCP-fetched content. External sources are first-class.
- **Weight**: each source's declared influence, expressed as a value between 0.0 and 1.0. 1.0 is full authority; 0.0 is none. Weight propagates through the pipeline: it governs extraction prominence, ranking order, budget allocation, and downstream capability weighting.
- **Structure**: the shape of the context artifact — sections and categories. The system fills the structure the user provides. It does not invent structure. If source material does not map to a declared section, that section is left empty and the user is notified.
- **Budget**: a token bound on the context artifact. Forces the system to choose. Weight governs how it chooses.

### Orchestration

Two additional inputs govern how the framework operates beyond compilation:

- **Capabilities**: what external systems make available. Context Synth discovers and maps these into the dependency graph. It does not implement them.
- **Adapters**: projection implementations that render resolved state into a specific output form. Adapters are interchangeable and do not influence graph resolution.

## Pipeline

### Phase 1: Context Compilation

```
Sources → Snap → Extract → Rank → Assemble → Verify → Context Artifact
```

**Snap** — fetch and normalize configured sources into a source-type-agnostic snapshot. Attach provenance references and content hashes. This is the boundary between fetching and processing.

**Extract** — decompose snapshot content into discrete, attributable units of knowledge and classify each unit into the user-provided structure. A single source may yield many items; extractions are not 1:1 with sources. Each item inherits its source's weight and provenance. This is the judgment layer of the pipeline — segmentation and classification are LLM-backed and non-deterministic. The governance stages on either side (Snap before, Rank/Assemble/Verify after) exist to bound and make that judgment reviewable.

**Rank** — order items within sections by weight and relevance. This is where declared source influence becomes operational. When sources share equal weight, declaration order is the tiebreaker — earlier-declared sources take precedence.

**Assemble** — build the artifact within budget. Higher-weight items included first. Omissions recorded, not dropped.

**Verify** — surface what was included, what was omitted, and the source basis for each. Produces a reviewable verification surface — not an automated pass/fail gate. The output of Verify is material for human or agent review.

### Phase 2: Capability Resolution

```
Context Artifact → Graph Construction → Dependency Resolution → Projection
```

**Graph Construction** — build the capability dependency graph from the governed context artifact and available external capabilities. Context determines which capabilities are relevant and how they are weighted. The graph is explicit and inspectable.

**Dependency Resolution** — traverse the graph in dependency order and resolve each capability. The result is a unified resolved state that reflects the governed context and the synthesized capabilities.

**Projection** — render resolved state through an adapter into observable output. The adapter determines the output form — not the output substance.

## Data Model

### Compilation

- **Source** — a canonical knowledge input with a declared weight
- **Snapshot** — normalized, hashed set of sources for a run
- **Extraction** — a candidate context item with inherited weight and source reference
- **SynthPlan** — ranked extractions grouped by section
- **ContextArtifact** — the governed, bounded, source-backed compilation

### Resolution

- **Capability** — a declared external system capability with inputs, outputs, and dependencies
- **CapabilityGraph** — the dependency graph of capabilities constructed from governed context
- **ResolvedState** — the unified output of graph resolution
- **Projection** — the observable output rendered by an adapter

## Invariants

These hold across all modes and versions:

- **Source traceability** — every included item traces to a specific source, reference, and content hash
- **Boundedness** — the context artifact stays within the declared budget
- **Human reviewability** — what was included, omitted, and why is always inspectable
- **Evidence-backed synthesis** — generated output is never self-justifying; source alignment remains visible
- **Governance without LLMs** — provenance, boundedness, weight ordering, and reviewability all survive without LLM-backed stages. Synthesis capability does not — semantic extraction and classification require the judgment layer. Without it, the pipeline falls back to assembling source contents directly in config order within budget.
- **Deterministic conflict resolution** — when sources of equal weight conflict, declaration order governs precedence. Conflict resolution is always inspectable, never implicit.
- **Derived, not canonical** — the context artifact is a working set, not a replacement for source documents
- **Projection independence** — adapters render resolved state but do not influence graph construction or resolution. Changing the adapter changes the output form, not the output substance.
- **Governance propagation** — source traceability and weight influence survive from compilation through graph resolution. A resolved capability traces back to the governed context that activated it.

## Architectural Boundaries

- The source/snapshot layer presents a source-type-agnostic interface. Adding a new source type is a new adapter, not a pipeline change.
- The pipeline has an explicit determinism boundary. Snap is deterministic. Extract is LLM-backed and non-deterministic — this is where the system exercises judgment. Rank, Assemble, and Verify are deterministic: same ranked inputs and budget produce the same output.
- Source tracking survives all pipeline stages and modes.
- The context artifact is a stable interface between compilation and resolution. The compiler produces it; the resolution phase consumes it. Neither phase knows the other's internals.
- Adapters are interchangeable. The same resolved state projects into different output forms without changes to the graph or the compilation.
- Capability discovery is external to the framework core. New capabilities are available by connecting new systems, not by changing the framework.

## What Context Synth Is Not

- Not a new source of truth
- Not a wiki or documentation system
- Not a generic RAG layer — RAG retrieves context at query time to augment a prompt. Context Synth compiles a governed artifact *before* any query. Sources are declared, weighted, and bounded into a reviewable snapshot. RAG answers "what is relevant to this question"; Context Synth answers "what knowledge is allowed to shape reasoning, and how much."
- Not an agent — Context Synth orchestrates a declared dependency graph over governed context. It does not set goals, select actions, or decide what to do next. Agents may consume its output or expose capabilities that it orchestrates — but the framework itself does not reason on their behalf.
- Not a capability runtime — Context Synth does not implement or host capabilities. It orchestrates capabilities that external systems provide.
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
7. System capabilities orchestrated over governed context produce more reliable outputs than capabilities operating over ad hoc context.
8. A dependency graph of external capabilities can be resolved deterministically given the same context artifact and capability set.
9. The same resolved state can be meaningfully projected into different output forms.

## Long-Term Direction

*The following are orientations, not planned features or commitments. They describe where the project points, not what it promises.*

Prove that governed, weighted, source-backed context is a better foundation for downstream system behavior than ad hoc context collection.

Mixed-source ingestion including MCP-backed inputs, explicit weight with visible downstream effects, bounded artifact creation with recorded omissions, provenance on every item, and a synthesis mode that shows why a conclusion was reached.

Context compilation evolves into capability orchestration — the governed artifact drives a dependency graph of external system capabilities, resolved into dynamic outputs. Adapters project the same resolved state into different forms: music, interfaces, automation, APIs.

The orchestration model evolves toward richer dependency semantics, reactive graph resolution, graph visualization and debugging tools, and broader capability ecosystems.

**Context Synth governs what knowledge is allowed to shape reasoning and system behavior, how strongly it shapes it, and how that influence remains visible. The aim is for it to become an infrastructure layer for orchestrating system capabilities over governed context.**