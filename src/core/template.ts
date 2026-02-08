import { readFile } from "fs/promises";
import { TemplateHeading } from "./types";

/**
 * Loads a markdown template file and parses its headings.
 * @param templatePath The absolute path to the template markdown file.
 * @returns A promise that resolves to an array of TemplateHeading objects.
 * @throws Error if the template file cannot be read.
 */
export async function loadTemplate(
  templatePath: string
): Promise<TemplateHeading[]> {
  const content = await readFile(templatePath, "utf-8");
  return parseTemplateHeadings(content);
}

/**
 * Parses a markdown string and extracts its ATX headings.
 * @param markdown The markdown content as a string.
 * @returns An array of TemplateHeading objects, each representing a heading with its level and text.
 */
function parseTemplateHeadings(markdown: string): TemplateHeading[] {
  const headings: TemplateHeading[] = [];
  const lines = markdown.split(/\r?\n/);
  const seen = new Set<string>();
  const idRe = /\s*\{#([a-z0-9][a-z0-9._-]*)\}\s*$/i;

  for (const line of lines) {
    const match = line.match(/^(#{1,6})\s+(.*)/);
    if (match) {
      const level = match[1].length;
      const raw = match[2].trim();
      const idMatch = raw.match(idRe);
      const slotId = idMatch?.[1];
      const heading = raw.replace(idRe, "").trim();

      if (level > 1) {
        if (!slotId) throw new Error(`Template is missing {#slotId} for heading: "${raw}"`);
        const norm = slotId.toLowerCase();
        if (seen.has(norm)) throw new Error(`Template has duplicate slotId: "${slotId}"`);
        seen.add(norm);
      }

      headings.push({
        level,
        heading,
        slotId,
      });
    }
  }

  return headings;
}

/** Returns routable slots (level > 1) as {slotId, heading}. */
export function getRoutableSlots(template: TemplateHeading[]): Array<{ slotId: string; heading: string }> {
  return template
    .filter((h) => h.level > 1)
    .map((h) => ({ slotId: h.slotId!, heading: h.heading }));
}
