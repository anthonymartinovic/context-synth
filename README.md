<div align="center">
  <img width="196" height="150" alt="ChatGPT Image Feb 9, 2026, 07_48_44 PM" src="https://github.com/user-attachments/assets/512488ba-9234-4c3c-a495-9b68ff45542a" />
</div>

---

Context Synth is a governed context runtime. It compiles knowledge from declared sources into a bounded, traceable context artifact shaped by explicit source weight. Every inclusion, omission, and influence is inspectable.

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
- Not a RAG layer
- Not an agent
- Not a prompt manager or model provider

## Status

Pre-implementation. The project specification, system design, and v0.1 design are complete. The v0.1 implementation target is a single Go binary (`cs`) operating over local markdown sources.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — thesis, pipeline, data model, invariants
- [System Design](docs/SYSTEM_DESIGN.md) — system-level technical shape
- [v0.1 Design](docs/design/v0.1.md) — first implementation slice

## License

See [LICENSE](LICENSE).
