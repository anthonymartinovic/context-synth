import * as vscode from "vscode";
import { join, isAbsolute } from "path";
import { loadConfig } from "../../core/config";
import { runEngine } from "../../core/engine";
import { createHeuristicRouter } from "../../core/routers/heuristic/router";
import { createCursorRouter } from "./routers/cursor/router";

export type PipelineResult = {
  outputPath: string;          // absolute
  outputDisplayPath: string;   // as configured (usually relative)
  loadedFileCount: number;
  chunkCount: number;
  usedTemplate: boolean;
  filledSlots?: number;
  totalSlots?: number;
  orphanCount?: number;
};

export async function runSynthesisPipeline(
  configPath: string,
  workspaceRoot: string,
  extensionRoot: string,
  progress: vscode.Progress<{ message?: string }>,
  token: vscode.CancellationToken
): Promise<PipelineResult> {
  const config = await loadConfig(configPath);
  if (token.isCancellationRequested) throw new vscode.CancellationError();

  if (!config.template) {
    config.template = join(extensionRoot, "templates", "DEFAULT.md");
  } else if (!isAbsolute(config.template)) {
    // Relative paths are interpreted as workspace paths.
    // (Core will resolve these against workspaceRoot.)
    config.template = config.template;
  }

  const mode = config.routing?.mode ?? "heuristic";
  const router =
    mode === "cursor"
      ? createCursorRouter(token, config.routing?.model)
      : createHeuristicRouter();

  const res = await runEngine(config, workspaceRoot, router, {
    report: (message) => progress.report({ message }),
    isCancelled: () => token.isCancellationRequested,
  });

  return {
    outputPath: res.outputPath,
    outputDisplayPath: res.outputDisplayPath,
    loadedFileCount: res.loadedFileCount,
    chunkCount: res.chunkCount,
    usedTemplate: true,
    filledSlots: res.filledSlots,
    totalSlots: res.totalSlots,
    orphanCount: res.orphanCount,
  };
}
