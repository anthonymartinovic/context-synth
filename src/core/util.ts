import { createHash } from "crypto";

export function sha256(content: string): string {
  return createHash("sha256").update(content).digest("hex");
}

export function safeParseJsonArray(text: string): any[] | null {
  const trimmed = text.trim();

  try {
    const result = JSON.parse(trimmed);
    if (Array.isArray(result)) return result;
  } catch {
    // fall through
  }

  const fenceMatch = trimmed.match(/```(?:json)?\s*([\s\S]*?)\s*```/i);
  if (fenceMatch?.[1]) {
    try {
      const result = JSON.parse(fenceMatch[1].trim());
      if (Array.isArray(result)) return result;
    } catch {
      // fall through
    }
  }

  const arrayMatch = trimmed.match(/\[[\s\S]*\]/);
  if (arrayMatch) {
    try {
      const result = JSON.parse(arrayMatch[0]);
      if (Array.isArray(result)) return result;
    } catch {
      // fall through
    }
  }

  return null;
}
