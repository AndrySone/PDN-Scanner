package model

import "time"

type FileTask struct {
	ScanID string
	Path   string
	Format string
	Size   int64
}

type Finding struct {
	Category      string   `json:"category"`
	Type          string   `json:"type"`
	Count         int      `json:"count"`
	Confidence    float64  `json:"confidence"`
	MaskedSamples []string `json:"masked_samples"`
}

type FileReport struct {
	Path            string    `json:"path"`
	Format          string    `json:"format"`
	Findings        []Finding `json:"findings"`
	UZLevel         string    `json:"uz_level"`
	UZReasons       []string  `json:"uz_reasons"`

	Status          string   `json:"status"`           
	Errors          []string `json:"errors"`
	ErrorCategories []string `json:"error_categories"` 

	DurationMS int64 `json:"duration_ms"`
}

type ScanReport struct {
	ScanID      string       `json:"scan_id"`
	StartedAt   time.Time    `json:"started_at"`
	FinishedAt  time.Time    `json:"finished_at"`
	FilesTotal  int          `json:"files_total"`
	FilesParsed int          `json:"files_parsed"`
	FilesFailed int          `json:"files_failed"`
	Results     []FileReport `json:"results"`
}