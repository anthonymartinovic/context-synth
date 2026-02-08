import { mkdir, writeFile } from "fs/promises";
import { dirname } from "path";
import { CSConfig } from "./types";

/** Writes the synthesized output markdown file to disk. */
export async function writeOutput(
  synthesizedMd: string,
  cfg: CSConfig
): Promise<void> {
  const outputPath = cfg.context?.path ?? "CONTEXT.md";

  await mkdir(dirname(outputPath), { recursive: true });

  await writeFile(outputPath, synthesizedMd, "utf-8");
}
