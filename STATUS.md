# Status

_Last updated: v0.1.0_

## What's Working

**The core loop is proven end-to-end.** Governed context compiled from diverse sources drives a capability graph that resolves into dynamic output. The framework provided the infrastructure; DJ-V (the first example application) provided the declarations and adapter.

### Framework

**Three-tier architecture holds.** Protocol defines contracts, engine compiles context, runtime resolves capabilities. The tiers communicate through the protocol and do not know each other's internals.

**Two source types.** Markdown documents and MP3 audio files compile into a single governed artifact. Audio analysis runs locally via librosa with HPSS-based feature extraction (tempo, key, spectral characteristics, timbre). Adding audio was a new integration, not a framework modification — the source provider interface extends without pipeline changes.

**Machine-readable JSON artifact.** The context artifact carries full governance metadata — provenance, weights, content hashes, omissions, budget, mode — from compilation through resolution.

**`cs synth` and `cs run` operate independently.** The context artifact is the stable interface between compilation and resolution. The protocol boundary is visible at the CLI level.

**Engine pipeline is complete.** Snap → Extract → Rank → Assemble → Verify → Render (JSON Context Artifact). The determinism boundary is explicit: Snap deterministic, Extract LLM-backed, Rank/Assemble/Verify deterministic.

**Runtime resolves capability graphs.** Graph construction from static capability declarations, topological sort for dependency ordering, and executor dispatch (LLM and subprocess) all function as designed.

**Source provider interface extends without pipeline changes.** Markdown and MP3 audio integrations ship. Adding audio required a new integration, not a framework modification.

**Llama as default LLM provider.** Local inference via ollama — no proprietary API dependency. Used in both the engine (extraction/classification) and the runtime (context interpretation). The `llm.Client` interface remains provider-agnostic.

**Deterministic mode (`--no-llm`) is fully functional.** Core governance properties — boundedness, provenance, reviewability, omission recording — hold without an LLM.

**Source traceability and weight influence survive end-to-end.** Every included item traces to a source path and content hash. Changing source weights visibly changes artifact contents and downstream capability outputs.

**Budget enforcement is strict.** The artifact never exceeds the configured token budget. Omitted items are recorded with reason, not silently dropped.

**Parallel LLM extraction.** Extraction calls run concurrently — one goroutine per source — cutting wall-clock time roughly in half.

### DJ-V

**Audio reference tracks and markdown documents compile into a governed artifact.** Librosa analyzes each MP3 using HPSS-based feature extraction (tempo, key, spectral characteristics, timbre). Document content and audio features occupy the same governed artifact, ranked and assembled by weight.

**The artifact drives a two-node capability graph.** `interpret_context` (Llama derives music parameters from governed context) → `generate_audio` (MusicGen produces audio from those parameters). The graph resolves in dependency order via the framework's runtime.

**Different context produces different music.** Different source documents, different reference tracks, and different weight configurations produce different music parameters, which produce audibly different audio. Governed context materially shapes the output.

**The full pipeline runs end-to-end.** `cs synth` → artifact → `cs run` → music playback. The adapter (a thin Deno script) plays the generated audio — it renders resolved state without influencing how that state was produced.

## What's Not Working Well

**LLM extraction is non-deterministic.** Run `cs synth` twice with the same config and you'll get different chunking and different section classifications. This is inherent to LLM-backed extraction but sits uncomfortably against governance goals. There is no verification that items were classified into the *right* section — only that they were classified into *a* valid section.

**Section classification has no fallback.** If the LLM classifies an extraction into a section name that doesn't match any declared section (e.g., due to hallucination or casing), the item is silently omitted with reason "no matching section." Misclassified items disappear without warning.

**The verification surface is minimal.** It reports included/omitted items with source paths and weights. Does not yet support policy-aware verification or structured review workflows.

**The artifact format is explicitly unstable.** It is a working format for v0.1.0, not a protocol commitment. It will change.

**No retry or timeout on LLM calls.** A transient failure aborts the entire run. A hung ollama request blocks indefinitely.

**Audio analysis is slow.** Librosa analysis of 4 MP3 files takes ~2 minutes. Acceptable for a demo; needs caching or parallelism for real workflows.

**MusicGen on Apple Silicon MPS.** EnCodec decoder hits a conv1d channel limitation (>65536 channels). Generation falls back to CPU, taking ~2 minutes for 30 seconds of audio. This is a PyTorch/MPS constraint, not a framework issue.

## What's Deferred

- TypeScript SDK
- Dynamic capability discovery (capabilities are statically declared)
- Multiple simultaneous adapters
- Reactive or continuous runtime mode
- Stable protocol format commitment
- Capability failure handling and partial resolution
- Graph visualization and debugging tools
- Sub-file provenance (paragraph, span, byte range)
- MCP-backed source ingestion
- Caching and incremental regeneration
- Per-section LLM configuration
- Artifact diffing and drift detection
- CI integration

## Known Rough Edges

- `cs init` always writes to `contextsynth.yml` in the current directory — the `--config` flag has no effect on it.
- Source paths in the artifact are relative to the config file location, which can look odd when running `cs` from a different working directory.
