package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg, path); err != nil {
		return Config{}, err
	}

	resolveDefaults(&cfg)
	resolvePaths(&cfg, path)
	return cfg, nil
}

func validate(cfg *Config, path string) error {
	if cfg.Version == "" {
		return fmt.Errorf("config: version is required")
	}
	if cfg.Version != "1" {
		return fmt.Errorf("config: unsupported version %q (expected \"1\")", cfg.Version)
	}
	if cfg.Budget <= 0 {
		return fmt.Errorf("config: budget must be a positive integer")
	}
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("config: at least one source is required")
	}
	for i, s := range cfg.Sources {
		if s.Markdown == "" && s.Audio == "" {
			return fmt.Errorf("config: source %d must have markdown or audio", i)
		}
		if s.Markdown != "" && s.Audio != "" {
			return fmt.Errorf("config: source %d must have only one of markdown or audio", i)
		}
		if s.Weight < 0.0 || s.Weight > 1.0 {
			return fmt.Errorf("config: source %d weight %.2f out of range [0.0, 1.0]", i, s.Weight)
		}
	}
	for i, sec := range cfg.Sections {
		if sec.Name == "" {
			return fmt.Errorf("config: section %d has empty name", i)
		}
		if sec.Budget < 0 {
			return fmt.Errorf("config: section %d budget must be non-negative", i)
		}
	}
	if cfg.LLM != nil {
		if cfg.LLM.Provider == "" {
			return fmt.Errorf("config: llm.provider is required when llm is configured")
		}
		if cfg.LLM.Provider != "llama" {
			return fmt.Errorf("config: unsupported llm provider %q", cfg.LLM.Provider)
		}
	}
	return nil
}

func resolveDefaults(cfg *Config) {
	if cfg.Output.Path == "" {
		cfg.Output.Path = "artifact.json"
	}
	normalizeSectionBudgets(cfg)
}

func normalizeSectionBudgets(cfg *Config) {
	if len(cfg.Sections) == 0 {
		return
	}
	var total float64
	for _, s := range cfg.Sections {
		total += s.Budget
	}
	if total == 0 {
		equal := 1.0 / float64(len(cfg.Sections))
		for i := range cfg.Sections {
			cfg.Sections[i].Budget = equal
		}
		return
	}
	for i := range cfg.Sections {
		cfg.Sections[i].Budget = cfg.Sections[i].Budget / total
	}
}

func resolvePaths(cfg *Config, configPath string) {
	dir := filepath.Dir(configPath)
	for i := range cfg.Sources {
		if cfg.Sources[i].Markdown != "" && !filepath.IsAbs(cfg.Sources[i].Markdown) {
			cfg.Sources[i].Markdown = filepath.Join(dir, cfg.Sources[i].Markdown)
		}
		if cfg.Sources[i].Audio != "" && !filepath.IsAbs(cfg.Sources[i].Audio) {
			cfg.Sources[i].Audio = filepath.Join(dir, cfg.Sources[i].Audio)
		}
	}
}
