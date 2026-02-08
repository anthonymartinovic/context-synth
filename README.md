# Context Synth

Context Synth is a **context synthesis protocol** with a **reference engine** that distills multiple context sources into one structured markdown file.

## Status (experimental)

This repository is in **experimental / active development**.

**Do not rely on this repository yet** for production work. Interfaces, defaults, and behavior may change without notice while the first end-to-end surface is stabilized.

## What this is

- **Protocol**: the stable contracts/invariants for context synthesis in this project (see `cs.schema.json` and the template slot rules).
- **Reference engine**: the TypeScript implementation that executes the protocol (core under `src/core/`).

## What it does (conceptually)

Given:
- multiple **sources** (weighted on a spectrum between 0 and 1)
- a **template** (Markdown headings with `{#slotId}`)
- a **router** (assigns chunks → slots)

Context Synth produces:
- one synthesized **context markdown file** written to your workspace (default output path: `CONTEXT.md`)

## Where this is today (0.0.1-alpha)

This `0.0.1-alpha` snapshot is an **in-progress** protocol + engine:

- **Core engine** exists and runs (markdown ingestion → chunking → routing → synthesis → output).
- **Schema + defaults** are defined in `cs.schema.json` (and applied in `src/core/config.ts`).
- **Routing**:
  - `heuristic` router is available (offline, deterministic).
  - `cursor` router exists in the VS Code adapter (`vscode.lm`), but the extension workflow is not yet polished.

No usage instructions yet: the install/run experience is intentionally not treated as stable.

## v0.1 milestone (first stable workflow)

- A reliable end-to-end run as a **Cursor / VS Code extension** (install → command appears → synthesis runs → output opens).
- Clear example config + example sources that demonstrate real heuristic behavior (no “template-shaped” inputs).
- Cursor routing works when available in the editor runtime.

## v1.0 milestone

- Published, polished **Cursor / VS Code extension**.
- **CLI entrypoint** powered by the same core engine.
- **Budget/token configuration** for non-deterministic routing (Cursor routing).
- **Comprehensive testing** (unit + integration + smoke tests for engine + adapters).

