# Context Synth -- Technical Design Document

## 1. Problem Statement

AI agents are gaining real autonomy in software engineering workflows. But there is no governed layer around **what those agents know**. Agent context today is invisible, ad-hoc, and unauditable.

The problem is not that agents lack access to code -- they can read files, search repos, inspect git history. The problem is that agents are blind to the **external knowledge** that makes code make sense: business domain knowledge, product requirements, architectural decisions, technical constraints, and active work context. This knowledge lives in scattered documents, wikis, tickets, and people's heads.

Governing this context is a critical infrastructure concern -- analogous to how source control governs code, container images govern build artifacts, and Terraform state governs infrastructure.

## 2. Product Definition

Context Synth is a **context governance tool** that bridges external knowledge sources to a repository via a standardized, traceable artifact.

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
                context.md (governed artifact, committed to repo)
                      |
                AI agents consume it alongside the code
```

Without Context Synth, agents see code but are blind to why it exists. With it, the "why" is governed and travels with the repo.

### What Context Synth Is

- Context governance infrastructure
- A bridge from scattered knowledge to a standardized, traceable artifact
- A CLI tool that captures, extracts, ranks, assembles, and verifies context

### What Context Synth Is Not

- Not a content generator -- it extracts and organizes, never fabricates
- Not a code analysis tool
- Not an agent framework or IDE plugin
- Not a prompt manager

## 3. Design Principles

### Source Traceability

Every piece of information in the output artifact traces back to where it came from: a specific source file, a specific paragraph, and a content hash. If an AI agent makes a decision based on context, you can trace exactly where that context originated.

### Content Integrity

Context Synth uses LLMs to extract key information from source documents. A separate model independently verifies that each extracted piece faithfully represents its source material. Extractions that fail verification are excluded. This is a provable guarantee, not a trust-based one.

### Five-Stage Pipeline

| Stage | What it does | Uses LLM | Deterministic |
|-------|-------------|-----------|---------------|
| **Snap** | Captures sources | No | Yes |
| **Extract** | Pulls key information from sources | Yes | No (cached) |
| **Rank** | Prioritizes by relevance | Yes | No (cached) |
| **Assemble** | Builds artifact within budget | No | Yes |
| **Verify** | Proves faithfulness to sources | Yes | No |

LLMs are contained to three stages. Each has a clearly defined role. Capture and assembly are fully deterministic.

### Multi-Model Architecture

Each LLM stage can use a different model optimized for its task:

| Stage | Strength needed | Ideal model class |
|-------|----------------|-------------------|
| Extract | Large context ingestion | Gemini (1M+ context) |
| Rank | Reasoning and judgment | Claude (nuanced reasoning) |
| Verify | Structured faithfulness check | GPT (structured output) |

This is the v1 target. Earlier versions use a single model across all stages.

### Graceful Degradation

If no LLM is configured, the tool falls back to **config-order assembly**: source file contents are included verbatim in the order they appear in configuration. This mode still provides governance value (standardized format, source tracking, token budgeting) but without extraction or intelligent ranking.

## 4. Core Concepts

### Source

A markdown file containing external knowledge relevant to the repository. Referenced by path in the configuration file.

### Section

A named category of context in the output artifact. Sections are **user-defined** in configuration -- there are no built-in or mandatory sections. Each section maps to one or more source files and has a proportional token budget.

### Extraction

The process of pulling discrete pieces of key information from source documents -- facts, decisions, constraints, requirements, concepts, rationale. Each extracted piece carries a reference back to its source file and paragraph.

### Synthesis Plan

The combined output of the extract and rank stages. Contains all extracted information with priority rankings, organized by section. Cached so that identical inputs don't require repeated LLM calls. Stored in `.contextsynth/cache/` and is auditable.

### Context Artifact

The output file (`context.md` by default). A structured markdown document containing ranked, verified extractions organized by section, with full source tracking. Committed to the repository and consumed by AI agents.

## 5. CLI Interface

Distributed as a single binary: `cs`.

### `cs init`

Generates a starter `contextsynth.yml`.

```
$ cs init
Created contextsynth.yml
```

If `contextsynth.yml` already exists, exits with an error.

### `cs synth`

Primary command. Runs the full pipeline and emits the context artifact.

```
$ cs synth
Snapped 8 sources (12,340 tokens)
Extracted 47 items across 3 sections
Ranked by relevance
Assembled 3 sections (7,420 tokens)
Verified 47 items against sources (47 passed, 0 failed)
Wrote context.md
```

**Flags:**
- `--config <path>` -- path to config file (default: `contextsynth.yml`)
- `--output <path>` -- override output file path
- `--verbose` -- log detailed pipeline information
- `--dry-run` -- show what would be assembled without writing output
- `--no-llm` -- skip extract/rank/verify, use config-order verbatim assembly

### `cs snap`

Diagnostic command. Shows what sources would be captured.

```
$ cs snap
Sources (8 files, 12,340 tokens):
  docs/domain/overview.md        2,100 tokens  sha256:a3f2c1...
  docs/domain/glossary.md        1,450 tokens  sha256:b7d4e2...
  docs/adrs/001-use-postgres.md    820 tokens  sha256:c9e8f3...
  ...
```

**Flags:**
- `--config <path>` -- path to config file (default: `contextsynth.yml`)
- `--json` -- output snapshot as JSON

## 6. Configuration

### Schema

```yaml
# Context Synth configuration
version: "1"

# Output settings
output:
  path: context.md

# Total token budget
budget: 10000

# LLM configuration for pipeline stages.
# Each stage can specify its own provider/model.
# If omitted, falls back to config-order assembly.
llm:
  extract:
    provider: gemini
    model: gemini-2.5-pro
    api_key_env: GEMINI_API_KEY
  rank:
    provider: gemini
    model: gemini-2.0-flash
    api_key_env: GEMINI_API_KEY
  verify:
    provider: gemini
    model: gemini-2.0-flash
    api_key_env: GEMINI_API_KEY

# User-defined sections
sections:
  - name: Domain Knowledge
    sources:
      - docs/domain/overview.md
      - docs/domain/glossary.md
    budget: 0.30

  - name: Architecture Decisions
    sources:
      - docs/adrs/001-use-postgres.md
      - docs/adrs/002-event-driven.md
      - docs/architecture.md
    budget: 0.25

  - name: Requirements
    sources:
      - docs/requirements/core-features.md
      - docs/requirements/non-functional.md
    budget: 0.25

  - name: Constraints
    sources:
      - docs/constraints.md
    budget: 0.20
```

### Rules

- `version` must be `"1"`.
- `output.path` defaults to `context.md`.
- `budget` defaults to `10000`.
- `llm` is optional. If omitted, falls back to config-order verbatim assembly.
- Each LLM stage can specify its own provider and model. If only some stages are configured, unconfigured stages inherit from the nearest configured one.
- `api_key_env` names the environment variable containing the API key. Keys are never stored in config.
- Section `budget` values are proportional weights, normalized at runtime.
- A single source file may appear in multiple sections.
- Source paths are relative to the config file's directory.
- Glob patterns are supported. Glob matches are sorted alphabetically.

### Validation Errors

- Missing `contextsynth.yml`
- Unknown `version`
- No sections defined
- Section with no sources
- Source path that doesn't resolve to an existing file
- Zero or negative budget values
- `api_key_env` set but environment variable missing

## 7. Pipeline Design

### Stage 1: Snap (Deterministic)

Captures configured source files.

1. Resolve paths (expand globs, sort alphabetically).
2. Read file contents.
3. Normalize (LF line endings, trim trailing whitespace, ensure final newline).
4. Compute SHA-256 content hash per file.
5. Estimate token count per file (word count / 0.75).
6. Produce snapshot: list of files with content, hashes, and token counts.

### Stage 2: Extract (LLM)

Reads all source documents and extracts key information.

**Input:** All source contents from the snapshot, section configuration.

**Process:**
1. Check cache (keyed by source + config hash). On hit, use cached results.
2. On miss, send all source contents to the extraction model. The prompt instructs it to pull out discrete pieces of information -- facts, decisions, constraints, requirements, concepts, rationale -- each referencing the source file and paragraph it came from.
3. Parse the structured response.
4. Cache the results.

**Output:** A set of extracted items, each with: text, category, source file, source paragraph number.

**Categories:** `fact`, `decision`, `constraint`, `requirement`, `concept`, `rationale`.

### Stage 3: Rank (LLM)

Prioritizes extracted information within each section.

**Input:** All extracted items, section configuration, token budgets.

**Process:**
1. Check cache (keyed by extraction output hash).
2. On miss, send all items to the ranking model. The prompt instructs it to assign priority based on relevance, importance, and uniqueness within each section.
3. Sort items by priority within each section.
4. Cache the results.

**Output:** Synthesis plan -- all items ordered by priority per section.

### Stage 4: Assemble (Deterministic)

Builds the context artifact from ranked items.

1. Within each section, items are in ranked order.
2. Calculate per-section token budgets from proportional weights.
3. Include items in priority order until the section budget is reached. Each item is either fully included or fully omitted.
4. Redistribute surplus from under-budget sections to over-budget sections iteratively.
5. Render structured markdown with source tracking metadata.

### Stage 5: Verify (LLM)

Confirms each extracted piece faithfully represents its source.

1. For each item in the assembled artifact, retrieve the source paragraph it references.
2. Send the (extracted text, source paragraph) pair to the verification model.
3. The model judges: `faithful` or `unfaithful` with a brief reason.
4. Unfaithful items are excluded and recorded in source tracking as `excluded (failed verification: <reason>)`.

For maximum integrity, the verification model should be from a different provider than the extraction model.

### Fallback Mode

When no LLM is configured:

- Stages 2, 3, and 5 are skipped.
- Stage 4 assembles source file contents verbatim in config order.
- Whole files included or omitted based on token budget.
- Source tracking still records everything.

## 8. Token Estimation

Word-based approximation: `word_count / 0.75`. No tokenizer dependency. The budget is inherently approximate since the consuming agent's context window includes other content beyond the Contextfile.

Section budgets are proportional weights normalized against the total. Surplus from under-budget sections redistributes iteratively to sections with omitted items.

Structural overhead (headings, source tracking comments, header) is accounted separately.

## 9. Source Tracking

### Per-Section

```markdown
## Domain Knowledge
<!-- cs:sources
  - text: "The domain model follows event sourcing patterns"
    source: docs/domain/overview.md
    paragraph: 5
    hash: sha256:a3f2c1...
    status: included
    priority: 1
  - text: "CQRS separates read and write models"
    source: docs/domain/overview.md
    paragraph: 12
    hash: sha256:a3f2c1...
    status: omitted (budget exceeded)
    priority: 5
-->
```

Every item -- included, omitted, or failed verification -- is recorded with its source reference and status.

### Artifact Header

```markdown
<!-- cs:metadata
  version: "1"
  generator: cs/0.1.0
  config_hash: sha256:f1e2d3c4...
  source_hash: sha256:a1b2c3d4...
  total_tokens: 7420
  sections: 3
  sources: 8
  mode: ranked
  models:
    extract: gemini-2.5-pro
    rank: gemini-2.0-flash
    verify: gemini-2.0-flash
-->
```

## 10. Project Structure

```
context-synth/
├── cmd/
│   └── cs/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── init.go
│   │   ├── synth.go
│   │   └── snap.go
│   ├── config/
│   │   └── config.go
│   ├── source/
│   │   └── source.go
│   ├── snap/
│   │   └── snap.go
│   ├── extract/
│   │   └── extract.go
│   ├── rank/
│   │   └── rank.go
│   ├── assemble/
│   │   ├── assemble.go
│   │   ├── budget.go
│   │   └── verify.go
│   ├── llm/
│   │   ├── provider.go
│   │   └── gemini.go
│   ├── cache/
│   │   └── cache.go
│   ├── token/
│   │   └── estimate.go
│   ├── render/
│   │   └── markdown.go
│   └── model/
│       └── model.go
├── testdata/
│   └── fixtures/
├── .contextsynth/
│   └── cache/
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

## 11. Dependencies

| Dependency | Purpose |
|-----------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML config parsing |
| `github.com/google/generative-ai-go` | Gemini API client |
| `google.golang.org/api` | Google API support |

Future versions add provider SDKs as multi-model support expands.

## 12. Roadmap

### v0.1 -- Validate extraction (build now)

Single Gemini call combining extract + rank. Structural verification only. Prove the concept works.

### v0.2 -- Split stages, semantic verification (projected)

Separate extract and rank into independent calls. Add LLM-powered semantic verification. Add caching.

### v0.3 -- Multi-model support (projected)

Per-stage model configuration. Add Claude and GPT providers. Cross-model verification.

### v1.0 -- Stable release

Paragraph-level extraction. Stable artifact schema. `cs validate` and `cs diff` commands. CI integration. Documentation.

The path from v0.1 to v1 will be shaped by what we learn from building. Intermediate milestones are directional, not guaranteed.

## 13. Success Criteria

Context Synth v1 succeeds if:

1. The tool extracts meaningful information from markdown documentation and assembles it into a governed context artifact.
2. Every piece of information in the artifact traces back to a specific source file and paragraph.
3. Cross-model verification proves extractions are faithful to their sources.
4. The artifact respects configured token budgets.
5. The multi-model architecture allows each pipeline stage to use the optimal model.
6. The tool falls back gracefully when no LLM is configured.
7. The tool runs as a single binary with no runtime dependencies beyond the filesystem and (optionally) network access to LLM APIs.
