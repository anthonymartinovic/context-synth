<div align="center">
  <img width="196" height="150" alt="ChatGPT Image Feb 9, 2026, 07_48_44 PM" src="https://github.com/user-attachments/assets/512488ba-9234-4c3c-a495-9b68ff45542a" />
</div>

---

A weighted epistemic runtime for governed AI reasoning.

AI frames probability as certainty. Context Synth governs what shapes the answer, and by how much. It compiles knowledge from diverse sources into a bounded, source-backed context artifact where every inclusion, omission, and influence is traceable and reviewable.

## What It Does

Four inputs govern a synthesis run:

- **Sources** — knowledge inputs to the pipeline (local files, repository documents, MCP-fetched content)
- **Weight** — each source's declared authority (0.0–1.0)
- **Structure** — the shape of the output (user-defined sections)
- **Budget** — a token bound on the output

The pipeline compiles these into a reviewable artifact:

```
Sources → Snap → Extract → Rank → Assemble → Verify → Contextfile
```

Weight propagates through every stage — it governs extraction prominence, ranking order, and budget allocation. The output is bounded, source-backed, and inspectable. It is a derived working set, not a new source of truth.

See the [Project Spec](docs/PROJECT_SPEC.md) and [System Design](docs/SYSTEM_DESIGN.md) for details on the pipeline, data model, invariants, and system-level technical shape.

## What It Is Not

- Not a wiki or documentation system
- Not a generic RAG layer
- Not an agent
- Not a prompt manager or model provider

## Status

Pre-implementation. The project specification, system design, and v0.1 design are complete. The v0.1 implementation target is a single Go binary (`cs`) operating over local markdown sources.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — thesis, pipeline, data model, invariants
- [System Design](docs/SYSTEM_DESIGN.md) — system-level technical shape
- [v0.1 Design](docs/design/v0.1.md) — first implementation slice

## Project Structure

```
docs/
  PROJECT_SPEC.md          # durable spec: thesis, pipeline, invariants
  SYSTEM_DESIGN.md         # system-level component and boundary design
  design/
    v0.1.md                # first implementation slice (Go, markdown-only)
  examples/
    sources/               # example source documents for testing
  archive/
    founding_insight.md    # origin document: the observation that motivated the project
    v0.0.1-alpha/          # artifacts from the prior TypeScript iteration
```

## License

See [LICENSE](LICENSE).
