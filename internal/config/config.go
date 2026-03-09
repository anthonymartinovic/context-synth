package config

type Config struct {
	Version  string        `yaml:"version"`
	Output   OutputConfig  `yaml:"output"`
	Budget   int           `yaml:"budget"`
	Sources  []SourceDecl  `yaml:"sources"`
	Sections []SectionDecl `yaml:"sections"`
	LLM      *LLMConfig    `yaml:"llm,omitempty"`
}

type OutputConfig struct {
	Path string `yaml:"path"`
}

type SourceDecl struct {
	Path   string  `yaml:"path"`
	Weight float64 `yaml:"weight"`
}

type SectionDecl struct {
	Name   string  `yaml:"name"`
	Budget float64 `yaml:"budget"`
}

type LLMConfig struct {
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
}
