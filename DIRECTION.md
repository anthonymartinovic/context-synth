# Direction

## Current State

Context Synth is in active development. The repository contains a snapshot of progress toward the architecture described in VISION.md.

What exists today:

- Core engine and protocol artifacts (in progress)
- Foundational module structure
- Initial adapter and routing implementations

## v1.0 Direction

- VS Code / Cursor extension (driver)
- CLI driver
- Additional adapters (new source types and ingestion paths)
- Expanded routing strategies
- Improved Context Unit normalization logic
- Budget enforcement (token count constraining output size)

## Immediate Priorities

- Migrate repository layout to match architectural contracts
- Enforce core/driver boundary (no driver-specific imports in core)
- Align terminology in codebase with VISION.md definitions
- Validate TypeScript compilation (`npx tsc --noEmit`)

Work advances directly toward the architecture described in VISION.md.
