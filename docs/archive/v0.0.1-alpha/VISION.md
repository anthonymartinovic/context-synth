# Vision

Context Synth is an open-source pipeline for synthesizing multiple context sources into a single structured context document usable by AI-assisted development workflows.

Pipeline:

Sources → Adapters → Context Units → Synthesis → Canonical Context Document

The user controls what the pipeline produces through three inputs: **sources** (what goes in, and how influential each source is), **templates** (what the output looks like), and **budget** (how much fits in the result). Everything between is handled by the pipeline.

Context Synth is defined by this flow and by the separation of responsibilities between its components.

---

## Architectural Structure

Context Synth separates defining context from assembling context.

The architecture has six components:

- **Context Synthesis Protocol (CSP)** defines what valid context looks like
- **Sources** are the user's weighted inputs to the pipeline
- **Adapters** convert raw source data into CSP-compliant Context Units
- **Templates** define the structure of the output document
- **Context Synthesis Engine (CSE)** routes, assembles, and renders the output
- **Drivers** expose the engine to runtime environments

---

### Context Synthesis Protocol (CSP)

The Context Synthesis Protocol defines the contracts and invariants that govern context synthesis.

CSP defines:

- Context Unit structure
- Template and Slot rules
- Routing contracts
- Validation rules
- Provenance tracking

CSP standardizes representation.  
It does not ingest sources or assemble documents.

---

### Sources and Weighting

Sources are the user's primary control surface.

A source is a weighted input — a document or data feed that the user selects for inclusion in the pipeline. Each source carries a weight (0..1) expressing how influential it is relative to other sources.

The user configures sources and their weights before the pipeline runs. This is a pre-pipeline decision: the user declares what goes in and how much it matters. They do not interact with Context Units or any internal pipeline representation.

Weight propagates from sources through the pipeline:

- Adapters normalize source content into Context Units, each inheriting the source's weight
- CSE uses weight to determine slot filling precedence, content ordering, and orphan inclusion
- Higher weight means greater influence on the final canonical document

Sources, templates, and budget together define the user's intent: sources determine what content is available and how influential it is; templates determine what the output looks like; budget determines how much fits.

---

### Budget

Budget is the user's control over output size, expressed as a token count.

A higher budget allows more content into the final document — more Context Units per slot, more supplementary material, more orphans included. A lower budget forces the engine to be aggressive: only the highest-weighted units survive, supplementary content is cut, and orphans are filtered more strictly.

Budget is a pre-pipeline decision, like sources and templates. The user sets it in configuration. CSE enforces it during synthesis.

Budget affects synthesis, not routing. Routing still assigns every unit to a slot (or marks it as an orphan). Budget determines how much of that assigned content actually makes it into the rendered output.

---

### Adapters

Adapters convert raw source data into CSP-compliant Context Units.

They:

- Retrieve or receive source material
- Extract relevant sections
- Normalize content into Context Unit structure
- Emit validated units

Adapters vary by source type. Each source format requires its own adapter.

Adapters face inward: they bridge raw data to the protocol.  
They do not perform synthesis or define output structure.

---

### Context Units

A Context Unit is the standardized representation of a discrete section of source content after normalization.

A single source produces many Context Units — one per identifiable section. Context Units are the atomic elements that flow through the pipeline.

Each Context Unit carries:

- Content (the normalized section)
- Provenance (which source it came from)
- Weight (inherited from the source — not configured per unit)

Context Units are an internal pipeline concept. Users do not see, configure, or interact with them directly. The user's control is at the source level; weight and provenance propagate automatically.

Context Units conform to CSP. They are produced by Adapters and consumed by CSE.

---

### Templates and Slots

A Template defines the structure of the output document.

A Slot is a named, routable position within a template. Each slot has a stable identifier that routing targets.

Templates express output intent. They determine what the final document looks like — which sections exist, in what order, with what structure. The engine does not invent structure; it fills the structure the template provides.

Slots are the contract between routing and synthesis:

- Routing maps Context Units to Slots
- Synthesis fills Slots with routed units and renders the result

---

### Routing

Routing is the process of mapping Context Units to Slots.

Given a set of Context Units and a set of Slots, a router produces assignments: each assignment maps a unit to a slot.

Routing rules:

- A Context Unit is assigned to at most one Slot
- Any unit not assigned is an orphan — never silently dropped
- Routing targets Slot identifiers, never display text

Routing is pluggable. Different routers may use different strategies (deterministic heuristics, model-assisted classification, etc.). The engine receives a router and executes it; it does not choose one.

---

### Context Synthesis Engine (CSE)

The Context Synthesis Engine consumes Context Units and a Template, routes units to slots, and produces the canonical context document.

It is responsible for:

- Executing routing (via an injected router)
- Filling slots with assigned units
- Ordering content by weight and precedence
- Enforcing the token budget
- Handling orphaned units (weight-aware and budget-aware filtering)
- Rendering the final structured output

Weight and budget interact at every stage. Higher-weighted units take precedence in slot filling and appear earlier in ordering. When the budget is tight, lower-weighted content is the first to be excluded — supplementary units within slots are cut before primary units, and orphans are filtered more aggressively. When the budget is generous, more content survives: supplementary material is included, and orphans pass the threshold more easily.

CSE assembles context.  
It does not define unit structure or template rules.

The engine exposes a single entry point: one function that accepts a config, a router, and options. This is the only integration point between the engine and drivers.

---

### Drivers

Drivers bridge runtime environments to the engine.

A driver:

- Loads configuration
- Selects a router appropriate to the environment
- Calls the engine
- Delivers output

Drivers face outward: they expose the engine to a specific runtime (editor extension, CLI, serverless function, etc.).

Drivers do not define the protocol, produce Context Units, or perform synthesis. They invoke the engine and wire it to their environment.

---

## Separation Principle

Context Synth maintains clear boundaries:

- CSP defines structure and contracts
- Sources carry weight (the user's intent)
- Budget constrains output size
- Adapters produce Context Units
- Templates define output shape
- CSE assembles output within budget
- Drivers expose the engine

This separation allows:

- Adding new source types without protocol change
- Adding new output structures by changing templates alone
- Multiple engines operating on shared units
- Multiple drivers exposing the same engine
- Evolution of synthesis strategies independently
- Ecosystem extensibility

---

## Key Contracts

### Determinism

Given the same template, Context Units, routing assignments, and budget, the rendered output must be byte-for-byte identical.

Non-determinism is allowed only in routers (e.g., model-assisted routing) and only affects assignments. Once assignments exist, synthesis and rendering are fully deterministic.

### Routing Contract

Routing produces assignments mapping Context Units to Slots.

- A unit maps to at most one slot.
- Unassigned units are orphans. Orphans are never silently dropped.
- No implicit fallbacks. Behavior changes require explicit configuration.

### Engine Contract

The engine exposes one function. It receives:

- A validated configuration
- A router (injected by the driver)
- Optional callbacks (progress reporting, cancellation)

This is the only integration point. The engine does not import driver-specific dependencies. A CLI, a test harness, or a serverless function can call the engine without the editor installed.

---

## Differentiation from Adjacent Systems

### MCP (Model Context Protocol)

MCP is a protocol for accessing tools and services. CSP is a protocol for structuring context.

MCP governs how systems communicate with external capabilities.  
CSP governs how context is represented and assembled.

Context Synth may consume data retrieved via MCP, but it is not MCP. They address different concerns.

### Retrieval-Augmented Generation (RAG)

RAG is a retrieval technique for identifying relevant material. CSE is a synthesis engine for assembling structured output.

RAG determines what to retrieve.  
Context Synth determines how to structure and assemble it.

Context Synth may operate on material surfaced by RAG, but it is not RAG. They operate at different stages.

---

## Terminology

| Term | Meaning |
|---|---|
| **Source** | A weighted input document or data feed — the user's primary control surface |
| **Weight** | A 0..1 float on a source indicating influence; propagates to Context Units and affects synthesis precedence |
| **Budget** | A token count that constrains output size — the user's control over how much content fits in the result |
| **Context Unit** | The standardized representation of a discrete section of source content — an internal pipeline concept, not user-configured |
| **Adapter** | A component that converts raw source data into CSP-compliant Context Units |
| **Template** | A document that defines the structure of the output |
| **Slot** | A named, routable position within a template |
| **Routing** | The process of mapping Context Units to Slots |
| **Router** | A concrete implementation of routing |
| **Assignment** | A single Context Unit → Slot mapping |
| **Orphan** | A Context Unit that routing could not assign to any slot |
| **Synthesis** | Filling slots with assigned units and rendering output |
| **CSP** | Context Synthesis Protocol — defines structure and contracts |
| **CSE** | Context Synthesis Engine — routes, assembles, and renders output |
| **Driver** | A component that bridges a runtime environment to the engine |

---

## Diagram

See `VISION_DIAGRAM.md`.
