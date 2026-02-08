import { resolve } from "path";
import type { CSConfig, RouterFn } from "./types";
import { loadFiles, chunkFiles } from "./input";
import { writeOutput } from "./output";
import { loadTemplate, getRoutableSlots } from "./template";
import { synthesizeWithAssignments, renderSynthesized } from "./synthesis";

/**
 * Core engine (no VS Code APIs).
 *
 * v0.x: sources and template are markdown files on disk.
 * Roadmap: other ingestions (e.g. MCP) should normalize to markdown before chunking/routing.
 */
export async function runEngine(
  cfg: CSConfig,
  workspaceRoot: string,
  router: RouterFn,
  opts?: { report?: (message: string) => void; isCancelled?: () => boolean }
): Promise<{
  outputPath: string;        // absolute
  outputDisplayPath: string; // as configured (usually relative)
  loadedFileCount: number;
  chunkCount: number;
  filledSlots: number;
  totalSlots: number;
  orphanCount: number;
}> {
  const report = opts?.report ?? (() => {});
  const isCancelled = opts?.isCancelled ?? (() => false);

  report("Loading sources…");
  const loadedFiles = await loadFiles(cfg, workspaceRoot);
  if (isCancelled()) throw new Error("Operation cancelled");

  report("Chunking sources…");
  const chunks = chunkFiles(loadedFiles);
  if (isCancelled()) throw new Error("Operation cancelled");

  report("Loading template…");
  if (!cfg.template) {
    throw new Error(
      "Config validation error: template is required at runtime. If omitted in cs.yaml, the adapter must supply the built-in default template."
    );
  }
  const templatePath = resolve(workspaceRoot, cfg.template);
  const template = await loadTemplate(templatePath);
  const slots = getRoutableSlots(template);
  if (isCancelled()) throw new Error("Operation cancelled");

  report("Routing chunks…");
  const nonEmpty = chunks.filter((c) => c.content.trim().length > 0);
  const assignments = await router(nonEmpty, slots);
  if (isCancelled()) throw new Error("Operation cancelled");

  report("Synthesizing…");
  const result = synthesizeWithAssignments(template, chunks, assignments);
  const output = renderSynthesized(result, cfg);

  report("Writing output…");
  const outputPath = resolve(workspaceRoot, cfg.context?.path ?? "CONTEXT.md");
  const absCfg: CSConfig = {
    ...cfg,
    context: { ...cfg.context, path: outputPath },
  };
  await writeOutput(output, absCfg);

  const filled = result.slots.filter((s) => s.primary !== null).length;
  const total = result.slots.filter((s) => s.level > 1).length;

  return {
    outputPath,
    outputDisplayPath: cfg.context?.path ?? "CONTEXT.md",
    loadedFileCount: loadedFiles.length,
    chunkCount: chunks.length,
    filledSlots: filled,
    totalSlots: total,
    orphanCount: result.orphans.length,
  };
}
