import * as vscode from "vscode";
import type { Chunk, ChunkAssignment, RouterFn } from "../../../../core/types";
import { safeParseJsonArray } from "../../../../core/util";

/**
 * This module implements the Cursor router.
 * It uses Cursor's built-in `vscode.lm` API to semantically assign
 * content chunks to template slots.
 */

export function createCursorRouter(
  token?: vscode.CancellationToken,
  preferredModel?: string
): RouterFn {
  return async (chunks, slots) => {
    const models = await vscode.lm.selectChatModels({});
    if (models.length === 0) {
      throw new Error(
        "No language models available in Cursor. " +
          "Make sure you have an active Cursor subscription or configured model provider."
      );
    }

    const want = preferredModel?.trim().toLowerCase();
    const model =
      want
        ? (models.find((m: any) => {
            const hay = `${m?.id ?? ""} ${m?.name ?? ""} ${m?.vendor ?? ""} ${m?.family ?? ""}`.toLowerCase();
            return hay.includes(want);
          }) ?? models[0])
        : models[0];

    const prompt = buildBatchPrompt(chunks, slots);
    const messages = [vscode.LanguageModelChatMessage.User(prompt)];
    const response = await model.sendRequest(messages, {}, token);

    let responseText = "";
    for await (const fragment of response.text) {
      responseText += fragment;
    }

    return parseBatchResponse(responseText, chunks, slots);
  };
}

function buildBatchPrompt(
  chunks: Chunk[],
  slots: Array<{ slotId: string; heading: string }>
): string {
  const chunkDescriptions = chunks
    .filter((c) => c.content.trim().length > 0)
    .map(
      (c, i) =>
        `--- CHUNK ${i} ---\nchunkId: ${c.chunkId}\nsource: ${c.sourceId} | file: ${c.filePath}\nheading: ${c.heading}\ncontent:\n${
          c.content.length > 600 ? c.content.slice(0, 600) + "\n[...truncated]" : c.content
        }`
    )
    .join("\n\n");

  return [
    "You are a document routing engine. Your job is to assign document chunks to the most appropriate template slot.",
    "",
    "TEMPLATE SLOTS (assign to exactly one of these):",
    ...slots.map((s) => `  - ${s.slotId}: ${s.heading}`),
    "",
    "CHUNKS TO ROUTE:",
    chunkDescriptions,
    "",
    "INSTRUCTIONS:",
    "- For each chunk, determine which template slot best matches its content.",
    "- Only include chunks that clearly belong to a slot. Omit chunks that don't fit anywhere.",
    "- Return ONLY a JSON array, no other text. Each element must be:",
    '  {"chunkId": "<exact chunkId>", "slotId": "<exact slotId>"}',
    "",
    "Return the JSON array now:",
  ].join("\n");
}

function parseBatchResponse(
  responseText: string,
  chunks: Chunk[],
  slots: Array<{ slotId: string; heading: string }>
): ChunkAssignment[] {
  const parsed = safeParseJsonArray(responseText);
  if (!Array.isArray(parsed)) {
    throw new Error(
      "Cursor router did not return a valid JSON array for routing assignments.\n" +
        `Response was: ${responseText.slice(0, 500)}`
    );
  }

  const validChunkIds = new Set(chunks.map((c) => c.chunkId));
  const validSlots = new Set(slots.map((s) => s.slotId));

  const assignments: ChunkAssignment[] = [];
  for (const item of parsed) {
    if (
      typeof item?.chunkId === "string" &&
      typeof item?.slotId === "string" &&
      validChunkIds.has(item.chunkId) &&
      validSlots.has(item.slotId)
    ) {
      assignments.push({ chunkId: item.chunkId, slotId: item.slotId });
    }
  }
  return assignments;
}
