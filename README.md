<div align="center">
  <img width="196" height="150" alt="ChatGPT Image Feb 9, 2026, 07_48_44 PM" src="https://github.com/user-attachments/assets/512488ba-9234-4c3c-a495-9b68ff45542a" />
</div>

---

A weighted epistemic runtime for governed AI reasoning.

AI frames probability as certainty. Context Synth governs what shapes the answer, and by how much.

## What It Does

Context Synth compiles knowledge from diverse sources into a bounded, source-backed context artifact. Every included item traces to its origin. Every omission is recorded. Weight declares how much each source influences the output.

```
Sources → Snap → Extract → Rank → Assemble → Verify → Contextfile
```

Four inputs govern a synthesis run:

- **Sources** — knowledge inputs to the pipeline
- **Weight** — each source's declared authority (0.0–1.0)
- **Structure** — the shape of the output (user-defined sections)
- **Budget** — a token bound on the output

## Status

Pre-implementation. The project specification, architecture, and v0.1 design are complete. Implementation in Go has not yet started.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — thesis, pipeline, data model, invariants
- [Architecture](docs/ARCHITECTURE.md) — system-level technical shape
- [v0.1 Design](docs/design/v0.1.md) — first implementation slice

## License

See [LICENSE](LICENSE).
