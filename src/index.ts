// Public core API (intentionally small). Adapters are never exported.

export type { CSConfig, ChunkAssignment, RouterFn } from "./core/types";
export { loadConfig } from "./core/config";
export { runEngine } from "./core/engine";
