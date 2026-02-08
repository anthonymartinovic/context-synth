import { readFile } from "fs/promises";
import yaml from "js-yaml";
import { CSConfig } from "./types";

/**
 * Loads and validates the configuration from a YAML file.
 *
 * @param configPath The path to the YAML configuration file.
 * @returns A promise that resolves to the loaded and validated CSConfig.
 * @throws An error if the file cannot be read, parsed, or if validation fails.
 */
export async function loadConfig(configPath: string): Promise<CSConfig> {
  const fileContents = await readFile(configPath, "utf-8");
  const parsedConfig = yaml.load(fileContents) as any;

  const config: CSConfig = {
    version: 1,
    sources: [],
    routing: { mode: "heuristic" },
    context: { path: "CONTEXT.md" },
    emit: { citations: true },
    chunking: { mode: "headings" },
  };

  // Deep merge parsed config with defaults
  Object.assign(config, parsedConfig);
  if (parsedConfig.context) Object.assign(config.context!, parsedConfig.context);
  if (parsedConfig.emit) Object.assign(config.emit!, parsedConfig.emit);
  if (parsedConfig.chunking) Object.assign(config.chunking!, parsedConfig.chunking);

  validateConfig(config);

  config.sources = config.sources.map(source => ({
    ...source,
    weight: source.weight ?? 1.0
  }));

  return config;
}

/**
 * Validates the given CSConfig object against a set of rules.
 *
 * @param config The configuration object to validate.
 * @throws An error if any validation rule is violated.
 */
function validateConfig(config: CSConfig): void {
  if (config.version !== 1) throw new Error("Config validation error: version must be 1.");
  if (!config.sources || !Array.isArray(config.sources) || config.sources.length === 0) {
    throw new Error("Config validation error: sources must be a non-empty array.");
  }

  if (config.template !== undefined) {
    if (typeof config.template !== "string" || config.template.trim().length === 0) {
      throw new Error("Config validation error: template must be a non-empty string when provided.");
    }
  }

  const sourceIds = new Set<string>();
  for (const source of config.sources) {
    if (!source.id) throw new Error("Config validation error: a source is missing an 'id'.");
    if (sourceIds.has(source.id)) throw new Error(`Config validation error: duplicate source id '${source.id}'.`);
    sourceIds.add(source.id);

    if (!source.path || typeof source.path !== "string" || source.path.trim().length === 0) {
      throw new Error(`Config validation error: source '${source.id}' has an invalid or empty 'path'.`);
    }

    const weight = source.weight;
    if (weight !== undefined && (!Number.isFinite(weight) || weight < 0 || weight > 1)) {
      throw new Error(`Config validation error: source '${source.id}' has an invalid 'weight'. Must be a finite number between 0 and 1.`);
    }
  }

  if (config.chunking?.mode && config.chunking.mode !== "headings") {
    throw new Error("Config validation error: chunking.mode must be 'headings'.");
  }

  if (
    config.routing?.mode &&
    config.routing.mode !== "heuristic" &&
    config.routing.mode !== "cursor"
  ) {
    throw new Error("Config validation error: routing.mode must be 'heuristic' or 'cursor'.");
  }
}
