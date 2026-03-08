# Context Synth -- Vision

## Problem

Repository-relevant context is scattered across requirements, architecture decisions, technical constraints, domain documents, tickets, and team memory. Repositories govern code well, but they do not govern this non-code context. As a result, the knowledge that makes code make sense is often invisible, ad hoc, and unauditable.

The problem is not access. Humans, agents, and tools can already read files, search repositories, inspect git history, and query external systems through MCP servers and other integrations. The gap is that raw access is not the same as bounded, reviewable context.

This matters for any consumer of repository context because repository work depends on more than code alone. The gap is more pronounced for agents, which consume context operationally and increasingly at scale, but it is not unique to them.

When this context is not governed, the failure mode is not merely inconvenience. Teams pay for it through repeated onboarding questions, inconsistent implementation decisions, review cycles spent re-explaining background, and avoidable reliance on tribal knowledge. The more important context lives across many documents, systems, and conversations, the more expensive this gap becomes.

The product thesis is therefore narrow and practical: repositories need a way to carry a governed working set of non-code context, not just pointers to wherever that context happens to live.

## Product Vision

Context Synth is a context compiler. It turns scattered external knowledge into a bounded, traceable, repo-local artifact that can travel with the repository and be consumed alongside the code.

In this document, a `consumer` is any reader or downstream system that uses repository context, including humans, agents, and tools.

```
External knowledge sources              Repository
├── domain knowledge                     ├── src/
├── requirements docs                    ├── internal/
├── ADRs                                 ├── go.mod
├── technical specs                      └── ...
├── constraints
└── ...
         \                              /
          \                            /
           ─── context-synth ─────────
                      |
                context artifact
                      |
                consumed alongside the code
```

Without Context Synth, consumers may have access to many sources but no curated working set of the knowledge that matters for a repository. With it, relevant context is compiled into a governed artifact that is bounded enough for intentional use and inspectable enough for meaningful review.

Category-wise, Context Synth is developer infrastructure for governed context. It is closer to a build step for repository-relevant knowledge than to a wiki, prompt manager, or agent runtime.

## Why Now

This category matters now because software projects already operate across many sources of repository-relevant knowledge, but most repositories still lack a disciplined way to carry that knowledge as a bounded and reviewable working set. Teams have become better at storing context, but not at governing which context should travel with the repository.

At the same time, modern teams already store key context in tools that are fetchable and machine-readable. The missing layer is not access plumbing. It is a disciplined way to compile that context into something repository-local, reviewable, and operationally useful.

## Who It Is For

Context Synth is for teams that already feel one or more of these pains:

- repository-critical context lives outside the repository
- contributors repeatedly need the same architectural or domain background to make safe changes
- important decisions depend on background that is not visible in the repository itself
- review quality depends on tribal knowledge that is not visible in the change itself
- teams need a context artifact that can be audited, versioned, and discussed like other repository assets

It is not aimed at teams that only need better note-taking, broader search, or a lighter documentation workflow.

## What Context Synth Is

- A context compiler
- Developer infrastructure for governed context
- A CLI that captures, prioritizes, assembles, and validates context

## What Context Synth Is Not

- Not a new source of truth; source documents remain canonical
- Not a wiki or documentation system
- Not a free-form content generator
- Not a generic RAG layer over arbitrary data
- Not a code analysis tool
- Not an agent framework, IDE plugin, or prompt manager

## Core Model

Context Synth operates on a small set of stable concepts:

- `Source`: canonical knowledge relevant to a repository
- `Section`: a user-defined category of context in the output artifact
- `Extraction`: a source-backed unit of derived context
- `SynthesisPlan`: the organized, prioritized intermediate result before final assembly
- `ContextArtifact`: the bounded, source-backed output consumed by any consumer

The initial implementation uses local markdown files as sources to validate the pipeline. The intended long-term ingestion path is MCP-fetched content such as Confluence, Notion, Jira, and other external systems. The product should not assume local files are the only source type.

## Invariants

These properties define the product more than any one file format or internal implementation:

### Source Traceability

Every included piece of context must trace back to a specific source and a stable source reference, along with content identity information such as a hash.

### Canonical Sources

The artifact is a derived working set, not a replacement for canonical documents.

### Boundedness

The output must stay within an explicit budget so included context remains intentional and reviewable.

### Human Reviewability

Humans must be able to inspect what was included, what was omitted, and where included context came from.

### Evidence-Backed Synthesis

LLMs may help extract, rank, or validate context, but generated output is never self-justifying. Source alignment must remain visible.

### Graceful Degradation

If no LLM is configured, the tool should still produce a useful governed artifact through deterministic assembly.

## Initial Value Proposition

The first version does not need to solve all context problems. It only needs to prove one concrete claim: a repository can ship a compact, source-backed context artifact that makes downstream work noticeably safer and more efficient than relying on scattered source material alone.

If Context Synth cannot produce an artifact that a reviewer would trust, a contributor would keep committed, and a team would treat as a meaningful repository asset, then the category claim is not yet proven.

## Product Shape

The long-term product shape is:

- configure sources and section budgets
- capture and normalize source material
- derive candidate context
- prioritize that context
- assemble a bounded artifact
- validate source alignment

This pipeline can evolve internally, and the final rendered artifact may become more configurable over time. What should remain stable is the contract that the artifact is bounded, source-backed, and reviewable.

## Scope Of This Document

This file captures the enduring product thesis, core model, and invariants. It is intentionally lighter on release-specific behavior and format decisions.

The concrete target for the first release lives in `docs/design/v0.1.md`. System-level implementation boundaries live in `docs/ARCHITECTURE.md`. Sequencing and likely future direction live in `docs/ROADMAP.md`.

## Success Conditions

Context Synth is successful if it proves that repositories can govern context with the same seriousness they apply to code:

1. Relevant non-code context can be compiled into a bounded repository artifact.
2. Included context remains traceable to canonical sources.
3. The artifact is useful to consumers while remaining inspectable by humans.
4. The system still provides governance value when running without LLM-backed stages.
5. The architecture can grow from local files toward MCP-fetched sources without changing the core model.
6. Teams can point to fewer context-related mistakes, less repeated explanation, or better shared understanding after adopting it.
