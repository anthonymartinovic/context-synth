<div align="center">
  <img width="196" height="150" alt="ChatGPT Image Feb 9, 2026, 07_48_44 PM" src="https://github.com/user-attachments/assets/512488ba-9234-4c3c-a495-9b68ff45542a" />
</div>

---

Context Synth is a governed context runtime for AI-assisted workflows. It compiles multiple knowledge sources into a single structured context document.

It:

- Accepts sources with declared weights that control how much each source influences the output
- Compiles source material into a structured, budget-bounded context document
- Tracks provenance so every item in the output traces back to its source
- Records what was included, what was omitted, and why

```
Sources → Snap → Extract → Rank → Assemble → Verify → Contextfile
```

See the [Project Spec](docs/PROJECT_SPEC.md) and [System Design](docs/SYSTEM_DESIGN.md) for details on the inner workings.

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
