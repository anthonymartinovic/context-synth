<div align="center">
  <img width="196" height="150" alt="Context Synth Logo" src="logo.png" />

  **A context graph runtime that turns knowledge into capabilities, delivering real-time user experiences.**
</div>

---

## The v1 Vision

Context Synth is a runtime where applications emerge from context rather than being written ahead of time. A collection of artifacts form a dependency graph — each carrying its own context, capabilities, and declared dependencies — and the application emerges from the live graph.

```
context graph (collection of artifacts with dependencies)
   ↓
application emerges from the graph
```

This is closer to how an operating system works than how any application works. An OS has no single artifact at the end — it has a collection of capabilities forming a dependency graph, and the user experience emerges from that graph at runtime. With AI, this model becomes far more achievable.

---

## Current State: v0.1.0

v0.1.0 is the first end-to-end release. It validates the core building blocks — governed context compilation, capability graph resolution, adapter projection — through a two-phase pipeline:

```
Sources → Snap → Extract → Rank → Assemble → Verify → Context Artifact (JSON)
Context Artifact → Graph Construction → Dependency Resolution → Projection
```

Today, sources compile into a single governed artifact, and a capability graph resolves downstream of it. This proves the pieces work. The v1 direction inverts the architecture so that artifacts form the graph directly rather than feeding into one. See [STATUS.md](STATUS.md) for full details on what's working, what isn't, and what's next.

For details on the v0.1.0 architecture, see the [Project Spec](docs/PROJECT_SPEC.md) and [System Design](docs/SYSTEM_DESIGN.md).

#### What v0.1.0 supports

- **Markdown and MP3 audio** source types
- Per-source weight declaration (markdown supports globs)
- Global token budget with strict enforcement
- User-defined section structure with LLM-backed classification
- Deterministic fallback mode with no LLM required
- Machine-readable JSON artifact with full governance metadata (provenance, weights, hashes, omissions)
- Capability graph construction and topological resolution
- LLM and subprocess executor mechanisms
- Adapter interface for output projection
- Llama as the default LLM provider (local inference via ollama)
- [DJ-V](examples/dj-v/) — first example application (documents + audio → governed artifact → capability graph → generative music)

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

# Produce a context artifact with LLM-backed extraction (requires ollama)
cs synth

# Resolve a capability graph over the artifact
cs run --artifact artifact.json --capabilities capabilities.json --adapter "your-adapter-cmd"
```

LLM-backed extraction requires ollama running locally with a Llama model:

```bash
brew install ollama
ollama serve &
ollama pull llama3.1:8b
```

## Usage

### `cs init`

Creates a starter `contextsynth.yml` in the current directory.

### `cs snap`

Resolves sources and prints a snapshot table showing each source's weight, token count, and content hash.

```
Source                                   Weight   Tokens Hash
sources/md/vibe.md                         1.00      340 a1b2c3d4
sources/audio/midnight_drive.mp3           0.80    1,200 c9d0e1f2
4 sources | 3,440 tokens estimated | budget: 10,000
```

### `cs synth`

Runs the compilation pipeline and writes a JSON context artifact.

| Flag | Description |
|------|-------------|
| `--config` | Path to config file (default: `contextsynth.yml`) |
| `--output` | Output file path (overrides config) |
| `--no-llm` | Deterministic mode — no LLM, flat weight-ordered output |
| `--dry-run` | Print output to stdout instead of writing a file |
| `--verbose` | Print pipeline summary to stderr |

### `cs run`

Loads a context artifact, resolves a capability graph, and projects output through an adapter.

| Flag | Description |
|------|-------------|
| `--artifact` | Path to context artifact (default: `artifact.json`) |
| `--capabilities` | Path to capability declarations (default: `capabilities.json`) |
| `--adapter` | Adapter command (e.g., `deno run adapter/src/index.ts`) |
| `--verbose` | Verbose output |

`cs synth` and `cs run` operate independently across the artifact boundary. The context artifact is the stable interface between them.

## Configuration

```yaml
version: "1"

output:
  path: artifact.json

budget: 10000

sources:
  - markdown: docs/domain/overview.md
    weight: 1.0
  - markdown: docs/adrs/*.md
    weight: 0.7
  - audio: sources/audio/reference.mp3
    weight: 0.8

sections:
  - name: Domain Knowledge
    budget: 0.30
  - name: Architecture Decisions
    budget: 0.25

llm:
  provider: llama
  model: llama3.1:8b
```

- **Sources** declare markdown file paths (globs supported) or audio file paths, each with a weight from 0.0 to 1.0.
- **Sections** define the output structure with proportional budget allocation.
- **LLM** is optional. Without it, the pipeline runs in deterministic fallback mode. Only `llama` is supported as a provider.
- Section classification requires the LLM. In `--no-llm` mode, output is flat weight-ordered.
- Audio sources require the `CS_AUDIO_ANALYZER` environment variable pointing to a librosa-based analysis script.

## Documentation

- [Project Spec](docs/PROJECT_SPEC.md) — what Context Synth is for and the rules it follows
- [System Design](docs/SYSTEM_DESIGN.md) — how the system is put together
- [Technology](docs/TECHNOLOGY.md) — foundational technology choices
- [Design](docs/design/) — versioned design documents (scope, architecture, acceptance criteria)
- [Plan](docs/plan/) — versioned implementation plans and milestones
- [Reflections](docs/reflections/) — versioned retrospectives on each release

## License

See [LICENSE](LICENSE).
