package model

import "time"

type Source struct {
	Path        string
	Content     []byte
	ContentHash string
	TokenCount  int
	Weight      float64
	DeclOrder   int
}

type SnapshotItem struct {
	Source      Source
	Provenance string
}

type Snapshot struct {
	Items []SnapshotItem
	Hash  string
}

type Extraction struct {
	Content     string
	Section     string
	TokenCount  int
	Weight      float64
	SourcePath  string
	ContentHash string
	DeclOrder   int
}

type SynthPlan struct {
	Included []Extraction
	Omitted  []Omission
}

type Omission struct {
	SourcePath string
	Reason     string
	Weight     float64
	TokenCount int
}

type Artifact struct {
	FrontMatter FrontMatter
	Sections    []ArtifactSection
	Omissions   []Omission
}

type FrontMatter struct {
	Generator  string
	Mode       string
	Snapshot   string
	Config     string
	Timestamp  time.Time
	Budget     int
	TokensUsed int
}

type ArtifactSection struct {
	Name  string
	Items []Extraction
}

type SectionBudget struct {
	Name       string
	Proportion float64
}
