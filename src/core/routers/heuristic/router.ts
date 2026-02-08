import { Chunk, ChunkAssignment, RouterFn } from "../../types";

/**
 * This module implements the heuristic router.
 * It's a pure, deterministic, offline router that uses word overlap scoring
 * to assign content chunks to template slots.
 */

const STOP_WORDS = new Set([
  "a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "from", "has", "have",
  "in", "into", "is", "it", "of", "on", "or", "our", "that", "the", "their", "then",
  "there", "these", "this", "to", "we", "with", "without", "you", "your",
]);

const tokenize = (s: string): string[] =>
  s
    .toLowerCase()
    .replace(/[^a-z0-9\s]/g, " ")
    .split(/\s+/)
    .map((t) => t.trim())
    .filter((t) => t.length >= 3 && !STOP_WORDS.has(t));

/**
 * Routes a single content chunk to a template slot using a lightweight, deterministic heuristic.
 * This is a local scoring approach: it compares words in the chunk against words in each
 * slot heading and picks the best match.
 * @param chunk The content chunk to route.
 * @param slotHeadings The available template slot headings.
 * @returns The heading of the most appropriate slot, or null if no match is found.
 */
function routeChunkHeuristically(
  chunk: Chunk,
  slotHeadings: string[]
): string | null {
  if (slotHeadings.length === 0) return null;

  const normalizedSlotHeadings = slotHeadings.map((s) => s.trim()).filter(Boolean);
  if (normalizedSlotHeadings.length === 0) return null;

  const chunkText = `${chunk.heading}\n${chunk.content}`.toLowerCase();

  // Fast path: exact heading match (case-insensitive).
  const exact = normalizedSlotHeadings.find((s) => s.toLowerCase() === chunk.heading.toLowerCase());
  if (exact) return exact;

  const chunkTokens = new Set(tokenize(chunkText));
  if (chunkTokens.size === 0) return null;

  let best: { slot: string; score: number } | null = null;
  for (const slot of normalizedSlotHeadings) {
    const slotTokens = tokenize(slot);
    if (slotTokens.length === 0) continue;

    let score = 0;
    for (const tok of slotTokens) {
      if (chunkTokens.has(tok)) score += 1;
    }

    // Prefer matches where the chunk's own heading overlaps with the slot heading.
    const headingTokens = new Set(tokenize(chunk.heading));
    for (const tok of slotTokens) {
      if (headingTokens.has(tok)) score += 2;
    }

    if (!best || score > best.score) best = { slot, score };
  }

  // Require at least some evidence of overlap to avoid random assignment.
  if (!best || best.score < 2) return null;
  return best.slot;
}

export function createHeuristicRouter(): RouterFn {
  return async (chunks, slots) => {
    const assignments: ChunkAssignment[] = [];
    const slotHeadings = slots.map((s) => s.heading);
    const byHeading = new Map(slots.map((s) => [s.heading, s.slotId]));
    for (const chunk of chunks) {
      const slot = routeChunkHeuristically(chunk, slotHeadings);
      if (slot) {
        const slotId = byHeading.get(slot);
        if (slotId) assignments.push({ chunkId: chunk.chunkId, slotId });
      }
    }
    return assignments;
  };
}
