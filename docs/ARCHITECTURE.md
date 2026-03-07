# Context Synth -- v1 Architecture

## Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                     CONTEXT SOURCES                         │
│                                                             │
│   Markdown files: domain docs, ADRs, requirements,         │
│   specs, constraints -- configured in contextsynth.yml      │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  1. SNAP                                         no LLM    │
│                                                             │
│  Read configured source files, normalize content,           │
│  compute hashes, estimate token counts.                     │
│                                                             │
│  Output: snapshot (files + hashes + tokens)                 │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  2. EXTRACT                                       Gemini    │
│                                                             │
│  LLM reads all source contents and pulls out key            │
│  information: facts, decisions, constraints, requirements,  │
│  concepts, rationale. Each item references its source       │
│  file and paragraph.                                        │
│                                                             │
│  Output: extracted items with source references             │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  3. RANK                                          Claude    │
│                                                             │
│  LLM assigns priority to each extracted item within its     │
│  section based on relevance, importance, and uniqueness.    │
│                                                             │
│  Output: items ordered by priority per section              │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  4. ASSEMBLE                                      no LLM   │
│                                                             │
│  Include items in priority order within each section's      │
│  token budget. Redistribute surplus. Render markdown.       │
│                                                             │
│  Output: assembled context artifact                         │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  5. VERIFY                                          GPT     │
│                                                             │
│  For each item: check it against its source paragraph.      │
│  A different model than extraction -- independent check.    │
│  Unfaithful items are excluded.                             │
│                                                             │
│  Output: verified artifact                                  │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                       context.md                            │
│                                                             │
│  Ranked, verified context organized by section.             │
│  Full source tracking. Committed to repo.                   │
│  Consumed by AI agents alongside the code.                  │
└─────────────────────────────────────────────────────────────┘
```

## Models

```
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│    GEMINI      │    │    CLAUDE      │    │     GPT       │
│                │    │               │    │               │
│  Reads all     │    │  Judges what  │    │  Checks each  │
│  source docs   │    │  matters most │    │  extraction   │
│  (1M+ tokens)  │    │               │    │  against its  │
│                │    │  Ranks items  │    │  source       │
│  Extracts key  │    │  by priority  │    │               │
│  information   │    │               │    │  Independent  │
│                │    │               │    │  cross-check  │
│                │    │               │    │               │
│  STAGE 2       │    │  STAGE 3      │    │  STAGE 5      │
└───────────────┘    └───────────────┘    └───────────────┘
```

## Data Flow

```
Markdown files
    │
    ▼
Snapshot (files, hashes, token counts)
    │
    ▼
Extracted items (text, category, source reference)
    │
    ▼
Ranked items (priority order per section)
    │
    ▼
Assembled artifact (within token budgets)
    │
    ▼
Verified artifact (unfaithful items removed)
    │
    ▼
context.md (with source tracking metadata)
```

## Traceability

At any point, you can:

- Trace any piece of information in the artifact back to its source file and paragraph
- Verify the source hasn't changed (content hash comparison)
- See why something was included or excluded (priority, budget, verification result)
- Know which models made which decisions

## Fallback

When no LLM is configured, stages 2, 3, and 5 are skipped. Source file contents are included verbatim in config order within token budgets. Still governed. Still traceable. Less intelligent.

## Version Progression

```
v0.1  Snap → [Extract+Rank] → Assemble → Verify(structural) → context.md
      Single model (Gemini), combined call, basic checks

v0.2  Snap → Extract → Rank → Assemble → Verify(semantic) → context.md
      Split stages, LLM-powered verification, caching (projected)

v0.3  Same pipeline, multi-model (Gemini/Claude/GPT) (projected)

v1.0  Full pipeline, stable schema, production-ready
```
