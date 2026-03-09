<div align="center">
  <img width="196" height="150" alt="ChatGPT Image Feb 9, 2026, 07_48_44 PM" src="https://github.com/user-attachments/assets/512488ba-9234-4c3c-a495-9b68ff45542a" />
</div>

---

Context Synth is a governed context runtime that compiles multiple knowledge sources into a single structured context document for humans, tools, and AI systems.

Users declare which sources are allowed to shape the output, how much influence each source carries, what structure the output should follow, and how much can fit within it. Context Synth produces a reviewable context file that shows what was included, what was omitted, and where each part came from.

Context Synth does the following:

- Compiles knowledge from repositories, local files, and external sources
- Applies explicit source influence throughout the assembly process
- Builds a structured context document within a defined budget
- Preserves source references and records omissions for review

```
Sources → Snap → Extract → Rank → Assemble → Verify → Contextfile
```

For details on the inner workings, see the [Project Spec](docs/PROJECT_SPEC.md) and [System Design](docs/SYSTEM_DESIGN.md).

## What It Is Not

- Not a wiki or documentation system
- Not a RAG layer
- Not an agent
- Not a prompt manager or model provider

## Status

Pre-implementation. The project specification, system design, and v0.1 design are complete. The v0.1 implementation target is a single Go binary (`cs`) operating over local markdown sources.

## Direction

Long term, Context Synth is intended to become the infrastructure layer that governs how context is constructed, versioned, and supplied to AI systems.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — what Context Synth is for and the rules it follows
- [System Design](docs/SYSTEM_DESIGN.md) — how the system is put together
- [v0.1 Design](docs/design/v0.1.md) — the first build target

## License

See [LICENSE](LICENSE).
