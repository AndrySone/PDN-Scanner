package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pii-scanner/internal/model"
)

type pyInferShmRequest struct {
	ShmPath string `json:"shm_path"`
	Offset  int64  `json:"offset"`
	Length  int64  `json:"length"`
	Format  string `json:"format"`
	Path    string `json:"path"`
}

type pyInferShmResponse struct {
	Findings []model.Finding `json:"findings"`
	Errors   []string        `json:"errors"`
}

func (p *PythonClient) InferFromShm(shmPath string, offset, length int64, format, path string) ([]model.Finding, []string, error) {
	reqBody := pyInferShmRequest{
		ShmPath: shmPath,
		Offset:  offset,
		Length:  length,
		Format:  format,
		Path:    path,
	}
	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, p.BaseURL+"/infer_shm", bytes.NewReader(b))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("python infer_shm status: %s", resp.Status)
	}

	var out pyInferShmResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, nil, err
	}
	return out.Findings, out.Errors, nil
}