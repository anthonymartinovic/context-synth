package protocol

import "time"

type Source struct {
	Path        string `json:"path"`
	Content     []byte `json:"-"`
	ContentHash string `json:"content_hash"`
	TokenCount  int    `json:"token_count"`
	Weight      float64 `json:"weight"`
	DeclOrder   int    `json:"decl_order"`
}

type SnapshotItem struct {
	Source     Source `json:"source"`
	Provenance string `json:"provenance"`
}

type Snapshot struct {
	Items []SnapshotItem `json:"items"`
	Hash  string         `json:"hash"`
}

type Extraction struct {
	Content     string  `json:"content"`
	Section     string  `json:"section,omitempty"`
	TokenCount  int     `json:"token_count"`
	Weight      float64 `json:"weight"`
	SourcePath  string  `json:"source"`
	ContentHash string  `json:"content_hash"`
	DeclOrder   int     `json:"-"`
}

type SynthPlan struct {
	Included []Extraction `json:"included"`
	Omitted  []Omission   `json:"omitted"`
}

type Omission struct {
	SourcePath string  `json:"source"`
	Reason     string  `json:"reason"`
	Weight     float64 `json:"weight"`
	TokenCount int     `json:"token_count"`
}

type Artifact struct {
	Meta      Meta              `json:"meta"`
	Sections  []ArtifactSection `json:"sections"`
	Omissions []Omission        `json:"omissions"`
}

type Meta struct {
	Generator    string    `json:"generator"`
	Mode         string    `json:"mode"`
	SnapshotHash string    `json:"snapshot_hash"`
	ConfigHash   string    `json:"config_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Budget       int       `json:"budget"`
	TokensUsed   int       `json:"tokens_used"`
}

type ArtifactSection struct {
	Name  string       `json:"name"`
	Items []Extraction `json:"items"`
}

type SectionBudget struct {
	Name       string
	Proportion float64
}
