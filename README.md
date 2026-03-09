<div align="center">
  <img width="196" height="150" alt="Context Synth Logo" src="logo.png" />
</div>

---

## Status

**v0.1.0-alpha.1** — early release, not yet stable. The `cs` binary compiles local markdown sources into bounded, traceable context artifacts with weight-based precedence, provenance tracking, and recorded omissions. LLM-backed extraction runs against Gemini (the only supported provider in this release). The artifact format is unstable and will change.

For a detailed account of what's working, what isn't, and what's deferred, see [STATUS.md](STATUS.md).

---

Context Synth is a governed context runtime that compiles multiple knowledge sources into a single structured context document for humans, tools, and AI systems.

Users declare which sources are allowed to shape the output, how much influence each source carries, what structure the output should follow, and how much can fit within it. Context Synth produces a reviewable context file that shows what was included, what was omitted, and where each part came from.

```
Sources → Snap → Extract → Rank → Assemble → Verify → Contextfile
```

For details on the inner workings, see the [Project Spec](docs/PROJECT_SPEC.md) and [System Design](docs/SYSTEM_DESIGN.md).

### What v0.1 supports today

- **Local markdown files** as the only source type
- Per-source weight declaration with glob expansion
- Global token budget with strict enforcement
- User-defined section structure (with LLM-backed classification)
- Deterministic fallback mode with no LLM required
- Provenance and omissions recorded in every artifact

### Long-term direction

Context Synth is intended to become an infrastructure layer that governs how context is constructed, versioned, and supplied to AI systems — ingesting from repositories, MCP servers, and other external sources; supporting stable artifact versioning, diffing, and drift detection; and integrating into CI pipelines. None of that is in the current version.

## What It Is Not

- Not a wiki or documentation system
- Not a RAG layer
- Not an agent
- Not a prompt manager or model provider

## Installation

Requires Go 1.26+.

```bash
go install github.com/anthonymartinovic/context-synth/cmd/cs@latest
```

Or build from source:

```bash
git clone https://github.com/anthonymartinovic/context-synth.git
cd context-synth
go build -o cs ./cmd/cs/
```

## Quick Start

```bash
# Generate a starter config
cs init

# Edit contextsynth.yml to declare your sources, weights, and budget

# Preview resolved sources
cs snap

# Produce a context artifact (deterministic, no LLM)
cs synth --no-llm

# Produce a context artifact with LLM-backed extraction
export GEMINI_API_KEY="your-key"
cs synth
```

## Usage

### `cs init`

Creates a starter `contextsynth.yml` in the current directory.

### `cs snap`

Resolves sources and prints a snapshot table showing each source's weight, token count, and content hash.

```
Source                                   Weight   Tokens Hash
docs/domain/overview.md                    1.00    2,340 a1b2c3d4
docs/adr/001-rest.md                       0.70    1,100 c9d0e1f2
3 sources | 3,440 tokens estimated | budget: 10,000
```

### `cs synth`

Runs the full pipeline and writes a context artifact.

| Flag | Description |
|------|-------------|
| `--config` | Path to config file (default: `contextsynth.yml`) |
| `--output` | Output file path (overrides config) |
| `--no-llm` | Deterministic mode — no LLM, flat weight-ordered output |
| `--dry-run` | Print output to stdout instead of writing a file |
| `--verbose` | Print pipeline summary to stderr |

## Configuration

```yaml
version: "1"

output:
  path: Contextfile

budget: 10000

sources:
  - path: docs/domain/overview.md
    weight: 1.0
  - path: docs/adrs/*.md
    weight: 0.7

sections:
  - name: Domain Knowledge
    budget: 0.30
  - name: Architecture Decisions
    budget: 0.25

llm:
  provider: gemini
  model: gemini-2.5-pro
  api_key_env: GEMINI_API_KEY
```

- **Sources** declare file paths (globs supported) with weights from 0.0 to 1.0.
- **Sections** define the output structure with proportional budget allocation.
- **LLM** is optional. Without it, the pipeline runs in deterministic fallback mode.
- Section classification requires the LLM. In `--no-llm` mode, output is flat weight-ordered.
- **Only Gemini is supported in v0.1.** The `provider` field exists for future extensibility but is currently ignored — the pipeline always uses Gemini.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — what Context Synth is for and the rules it follows
- [System Design](docs/SYSTEM_DESIGN.md) — how the system is put together
- [v0.1 Design](docs/design/v0.1.md) — the first build target
- [v0.1 Plan](docs/plan/v0.1.md) — implementation milestones

## License

See [LICENSE](LICENSE).
