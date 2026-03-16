# Context Synth -- Technology

## Purpose

This document records the foundational technology choices for Context Synth. These decisions are deliberate and should not change without explicit review. Implementation details that scope to a specific release belong in versioned design docs (`docs/design/`).

## Selection Principles

- Prefer open-source foundations. AI is trending toward utility infrastructure; open-source models are the durable foundation for that future. Proprietary models serve as frontier and premium options, not as dependencies.
- Prefer languages with strong concurrency primitives and explicit error handling for core infrastructure.
- Prefer technologies that support single-binary or minimal-dependency deployment.
- The developer-facing surface prioritizes ecosystem reach and type safety.

## Core Technologies

### Go — Engine and Runtime

Go is the implementation language for the engine (context compilation) and the runtime (capability resolution).

- Strong concurrency primitives for parallel pipeline stages and graph resolution
- Single binary deployment with no runtime dependencies
- Explicit error handling suits a governance-focused system where failures must be visible
- Performance characteristics appropriate for infrastructure that sits in the critical path

### Llama — Context Understanding

Llama is the primary model for context parsing, extraction, and classification in the engine's judgment layer.

The choice of an open-source model is deliberate. Context Synth is infrastructure — it should not depend on proprietary API access to function. Open-source models provide the control, reproducibility, and deployment flexibility that infrastructure demands.

The framework abstracts the LLM provider behind a client interface. Proprietary models (Gemini, OpenAI, Anthropic) can be used as alternative providers — they are not excluded, but they are not the default. The architecture does not assume or require them.

### TypeScript — SDK

TypeScript is the language for the developer-facing SDK — capability declaration, adapter implementation, and framework integration.

- Broad ecosystem reach for external system integrators
- Type system provides contract enforcement at the SDK boundary
- Natural fit for adapter implementations across diverse output domains

## Not In Scope

Technology choices for specific adapters or external capability implementations are not governed by this document. Those are decisions made by adapter and capability authors, constrained only by the protocol contracts defined in `docs/SYSTEM_DESIGN.md`.
