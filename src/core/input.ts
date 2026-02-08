import { readFile } from "fs/promises";
import { resolve } from "path";
import fg from "fast-glob";
import { CSConfig, LoadedFile, Chunk } from "./types";
import { sha256 } from "./util";

export async function loadFiles(
  cfg: CSConfig,
  cwd?: string
): Promise<LoadedFile[]> {
  let allLoadedFiles: LoadedFile[] = [];

  for (let idx = 0; idx < cfg.sources.length; idx++) {
    const source = cfg.sources[idx];
    const matchingFiles = await fg(source.path, { dot: true, cwd });
    matchingFiles.sort(); // Deterministic order

    for (const filePath of matchingFiles) {
      const resolvedPath = cwd ? resolve(cwd, filePath) : filePath;
      const content = await readFile(resolvedPath, "utf-8");
      const hash = sha256(content);
      const normalizedFilePath = filePath.replace(/\\/g, "/");

      allLoadedFiles.push({
        sourceId: source.id,
        sourceWeight: source.weight!,
        sourceOrder: idx,
        filePath: normalizedFilePath,
        content,
        hash,
      });
    }
  }

  allLoadedFiles.sort((a, b) => {
    if (a.sourceOrder !== b.sourceOrder) return a.sourceOrder - b.sourceOrder;
    if (a.filePath < b.filePath) return -1;
    if (a.filePath > b.filePath) return 1;
    return 0;
  });

  if (allLoadedFiles.length === 0) {
    throw new Error("Glob matches across all sources yielded zero files.");
  }

  return allLoadedFiles;
}

export function chunkFiles(files: LoadedFile[]): Chunk[] {
  const chunks: Chunk[] = [];

  for (const file of files) {
    const lines = file.content.split(/\r?\n/);
    let currentChunkContent: string[] = [];
    let currentHeading = "Document";
    let currentLevel = 0;

    const emitChunk = () => {
      const content = currentChunkContent.join("\n").trim();
      if (content.length > 0 || currentHeading === "Document") {
        chunks.push({
          chunkId: sha256(`${file.sourceId}:${file.filePath}:${currentHeading}:${content}`),
          sourceId: file.sourceId,
          sourceWeight: file.sourceWeight,
          sourceOrder: file.sourceOrder,
          filePath: file.filePath,
          heading: currentHeading,
          level: currentLevel,
          content: content,
          contentHash: sha256(content),
        });
      }
    };

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const headingMatch = line.match(/^(#{1,6})\s+(.*)/);

      if (headingMatch) {
        if (currentChunkContent.length > 0 || currentHeading !== "Document") {
          emitChunk();
        }
        currentLevel = headingMatch[1].length;
        currentHeading = headingMatch[2].trim();
        currentChunkContent = [];
      } else {
        currentChunkContent.push(line);
      }
    }
    emitChunk(); // Emit the last chunk
  }

  return chunks;
}
