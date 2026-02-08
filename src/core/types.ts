export type CSConfig = {
  version: 1;
  sources: SourceConfig[];
  template?: string;
  routing?: {
    mode?: "heuristic" | "cursor";
    model?: string; // preferred model name for Cursor routing (substring match)
  };
  context?: { path?: string };
  emit?: { citations?: boolean };
  chunking?: { mode?: "headings" };
};

export type SourceConfig = {
  id: string;
  path: string;
  weight?: number;
};

export type LoadedFile = {
  sourceId: string;
  sourceWeight: number;
  sourceOrder: number;
  filePath: string;
  content: string;
  hash: string;
};

export type Chunk = {
  chunkId: string;
  sourceId: string;
  sourceWeight: number;
  sourceOrder: number;
  filePath: string;
  heading: string;
  level: number;
  content: string;
  contentHash: string;
};

export type Group = {
  heading: string;
  chunks: Chunk[];
};

export type TemplateHeading = {
  slotId?: string;
  heading: string;
  level: number;
};

export type FilledSlot = {
  heading: string;
  level: number;
  primary: Chunk | null;
  supplementary: Chunk[];
};

export type OrphanChunk = Chunk;

export type SynthesisResult = {
  slots: FilledSlot[];
  orphans: OrphanChunk[];
};

export type RoutingMode = "heuristic" | "cursor";

export type ChunkAssignment = { chunkId: string; slotId: string };

export type RouterFn = (
  chunks: Chunk[],
  slots: Array<{ slotId: string; heading: string }>
) => Promise<ChunkAssignment[]>;

/** Deterministic comparator: higher weight wins, then source stack order, then path, then chunkId. */
export function compareChunks(a: Chunk, b: Chunk): number {
  if (a.sourceWeight !== b.sourceWeight) return b.sourceWeight - a.sourceWeight;
  if (a.sourceOrder !== b.sourceOrder) return a.sourceOrder - b.sourceOrder;
  if (a.filePath < b.filePath) return -1;
  if (a.filePath > b.filePath) return 1;
  if (a.chunkId < b.chunkId) return -1;
  if (a.chunkId > b.chunkId) return 1;
  return 0;
}
