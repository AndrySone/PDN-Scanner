package model

import "time"

type ScanSummary struct {
	ScanID      string    `json:"scan_id"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
	FilesTotal  int       `json:"files_total"`
	FilesParsed int       `json:"files_parsed"`
	FilesFailed int       `json:"files_failed"`

	ShmOK          int64 `json:"shm_ok"`
	ShmFail        int64 `json:"shm_fail"`
	HTTPFallbackOK int64 `json:"http_fallback_ok"`
	HTTPFail       int64 `json:"http_fail"`
}