import * as vscode from "vscode";
import { join } from "path";
import { readFile } from "fs/promises";
import { runSynthesisPipeline } from "./pipeline";

export function activate(context: vscode.ExtensionContext) {
  const extensionRoot = context.extensionUri.fsPath;
  const command = vscode.commands.registerCommand("context-synth.synth", () => runSynthesis(extensionRoot));
  context.subscriptions.push(command);
}

export function deactivate() {}

async function runSynthesis(extensionRoot: string) {
  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  if (!workspaceRoot) {
    vscode.window.showErrorMessage("Context Synth: Open a workspace folder first.");
    return;
  }

  const configPath = await findConfigFile(workspaceRoot);
  if (!configPath) {
    vscode.window.showWarningMessage("Context Synth: No cs.yaml found in workspace root.");
    return;
  }

  await vscode.window.withProgress({
    location: vscode.ProgressLocation.Notification,
    title: "Context Synth",
    cancellable: true,
  }, async (progress, token) => {
    try {
      const res = await runSynthesisPipeline(
        configPath,
        workspaceRoot,
        extensionRoot,
        progress,
        token
      );
      const outputPath = res.outputPath;

      const doc = await vscode.workspace.openTextDocument(outputPath);
      await vscode.window.showTextDocument(doc, { preview: false });

      vscode.window.showInformationMessage(
        `Context Synth complete — ${res.loadedFileCount} files, ${res.chunkCount} chunks → ${res.outputDisplayPath}`
      );
    } catch (err: any) {
      if (token.isCancellationRequested) return;
      vscode.window.showErrorMessage(`Context Synth error: ${err.message}`);
    }
  });
}

// ---------------------------------------------------------------------------
// Config file discovery (local to extension.ts for UI flow)
// ---------------------------------------------------------------------------

async function findConfigFile(root: string): Promise<string | null> {
  const candidates = ["cs.yaml", "cs.yml", "context-synth.yaml", "context-synth.yml"];
  for (const name of candidates) {
    const fullPath = join(root, name);
    try {
      await readFile(fullPath, "utf-8");
      return fullPath;
    } catch {
      // not found, try next
    }
  }
  return null;
}
