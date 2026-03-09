# Status

_Last updated: v0.1.0_

## What's Working

**The core thesis holds.** Declare sources with weights, set a budget, run `cs synth` — you get a bounded, traceable context artifact where weight visibly shapes what's included and what's omitted.

**Deterministic mode (`--no-llm`) is fully functional.** Sources are resolved, hashed, ranked by weight, assembled within budget, and rendered with provenance annotations and an omissions table. No LLM required. This path is stable and predictable.

**LLM-backed mode works end-to-end.** Content from N sources is extracted, classified into user-defined sections, ranked, and assembled. The sections don't have to map 1:1 to source files — the LLM synthesizes across source boundaries. Tested against Gemini 2.5 Flash.

**`cs snap` makes the source set inspectable.** Before running synthesis, you can see exactly what files will be included, their weights, token estimates, and content hashes.

**Budget enforcement is strict.** The artifact never exceeds the configured token budget. Omitted items are recorded with reason in an `# Omissions` appendix.

**Provenance is end-to-end.** Every included item traces back to a source file path and SHA-256 content hash taken at snapshot time. Source annotations are grouped by consecutive source in the output to reduce noise.

**Parallel LLM extraction.** Extraction calls are made concurrently — one goroutine per source file — rather than sequentially, cutting wall-clock time roughly in half.

## What's Not Working Well

**LLM extraction is non-deterministic.** Run `cs synth` twice with the same config and you'll get different chunking and different section classifications. This is inherent to LLM-backed extraction, but it sits uncomfortably against the tool's governance goals. There is currently no verification that items were classified into the *right* section — only that they were classified into *a* valid section.

**Runtime is slow for LLM mode.** With 3 source files, synthesis takes ~40 seconds. With larger source sets, this will get painful. Parallel extraction helps but doesn't solve the fundamental cost of one LLM call per source file.

**The `verify` package is a placeholder.** It computes a verification surface (included/omitted items with metadata) that is currently discarded after computation. Verification is defined as the artifact itself in v0.1 — front matter, inline annotations, and the omissions appendix — which is real but limited.

**Section classification has no fallback.** If the LLM classifies an extraction into a section name that doesn't match any declared section (e.g., due to hallucination or casing), the item is silently omitted with reason "no matching section." This could cause confusion.

**The artifact format is explicitly unstable.** The output format is a working format, not a contract. It will change.

## What's Deferred

- Stable artifact format and versioning
- Sub-file provenance (paragraph, span, byte range)
- MCP-backed source ingestion
- Caching and incremental regeneration
- Per-section LLM configuration
- Artifact diffing and drift detection
- CI integration

## Known Rough Edges

- `cs init` always writes to `contextsynth.yml` in the current directory — the `--config` flag has no effect on it.
- Source paths in the Contextfile are relative to the config file location, which can look odd when running `cs` from a different working directory.
- No retry logic on LLM API calls. A transient failure aborts the entire run.
