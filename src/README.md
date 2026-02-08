## `src/` — Code map

This folder contains the reference engine (core) and the VS Code/Cursor adapter surface.

### Layout

```
src/
  core/
    engine.ts        ← orchestrates the pipeline (core entrypoint)
    config.ts        ← load + validate `cs.yaml` + apply defaults
    input.ts         ← load markdown sources (glob) + chunk by headings
    template.ts      ← parse template headings + enforce `{#slotId}`
    synthesis.ts     ← fill slots + orphan filtering + render markdown
    output.ts        ← write output file to disk
    types.ts         ← shared types + deterministic chunk ordering
    util.ts          ← small pure helpers (sha256, etc.)
    routers/
      heuristic/
        router.ts    ← offline deterministic router (word overlap)

  adapters/
    vscode/
      extension.ts   ← registers command + UI/progress
      pipeline.ts    ← selects router + injects built-in template + runs core
      routers/
        cursor/
          router.ts  ← Cursor router (uses `vscode.lm`)

  index.ts           ← public exports (core only)
```

### Key rule

Core code must **never** import `vscode`. Router selection happens only in the adapter pipeline.

---

### End-to-end flow (what happens when you run “Synthesize”)

1. **VS Code command** (`adapters/vscode/extension.ts`)
   - Finds `cs.yaml` in the workspace root.
   - Runs the pipeline with progress + cancellation.

2. **Adapter pipeline** (`adapters/vscode/pipeline.ts`)
   - Calls `loadConfig(...)` (core) to parse + validate config.
   - Chooses the router based on `routing.mode`:
     - `heuristic` → `core/routers/heuristic/router.ts`
     - `cursor` → `adapters/vscode/routers/cursor/router.ts`
   - If `template` is omitted, injects the built-in template shipped with the extension (`templates/DEFAULT.md`).
   - Calls `runEngine(...)` (core).

3. **Core engine** (`core/engine.ts`)
   - Loads source files (`core/input.ts`).
   - Chunks by ATX headings (`core/input.ts`).
   - Loads + validates the template (`core/template.ts`).
   - Calls the injected router to produce assignments (`RouterFn`).
   - Synthesizes + renders the output markdown (`core/synthesis.ts`).
   - Writes the output file (`core/output.ts`).

---

### “Where do I look?” (debugging map)

- **Config problems** (missing fields, bad weights, wrong routing mode):
  - `core/config.ts` (`validateConfig`)
  - `cs.schema.json` (canonical config surface + defaults)

- **Template problems** (missing/duplicate/invalid `{#slotId}`):
  - `core/template.ts`
  - `templates/DEFAULT.md`

- **“Why did this chunk land there?”** (routing decisions):
  - Heuristic: `core/routers/heuristic/router.ts`
  - Cursor: `adapters/vscode/routers/cursor/router.ts`

- **“Why is the output shaped like this?”** (rendering/orphans/citations):
  - `core/synthesis.ts`

---

### Core vs adapter boundary (non-negotiable)

- **Core (`src/core/`)**:
  - No `vscode` imports, deterministic once assignments are provided.
  - Reads from disk using Node FS (v0.x constraint: markdown files on disk).

- **Adapters (`src/adapters/`)**:
  - May import `vscode`.
  - Own router selection and runtime integration (progress, cancellation, `vscode.lm`).

