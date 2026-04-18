package extract

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pii-scanner/internal/model"
)

type pyInferRequest struct {
	Path   string `json:"path"`
	Format string `json:"format"`
}

type pyInferResponse struct {
	Findings []model.Finding `json:"findings"`
	Errors   []string        `json:"errors"`
}

type PythonClient struct {
	BaseURL    string
	HTTPClient *http.Client
	MaxRetries int
}

func NewPythonClient(baseURL string) *PythonClient {
	return &PythonClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, 
		},
		MaxRetries: 3,
	}
}

func (p *PythonClient) Infer(path, format string) ([]model.Finding, []string, error) {
	reqBody := pyInferRequest{Path: path, Format: format}
	payload, _ := json.Marshal(reqBody)

	var timeout time.Duration
	switch format {
	case "pdf":
		timeout = 90 * time.Second
	case "tif", "tiff", "jpg", "jpeg", "png", "gif":
		timeout = 120 * time.Second
	case "docx", "xlsx", "xlsm":
		timeout = 60 * time.Second
	default:
		timeout = 45 * time.Second
	}

	var lastErr error
	for attempt := 0; attempt <= p.MaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.BaseURL+"/infer", bytes.NewReader(payload))
		if err != nil {
			cancel()
			return nil, nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			cancel()
			if attempt < p.MaxRetries {
				time.Sleep(backoff(attempt))
				continue
			}
			return nil, nil, fmt.Errorf("python infer failed after retries: %w", lastErr)
		}

		var out pyInferResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&out)
		_ = resp.Body.Close()
		cancel()

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("python worker 5xx: %s", resp.Status)
			if attempt < p.MaxRetries {
				time.Sleep(backoff(attempt))
				continue
			}
			return nil, nil, fmt.Errorf("python infer failed after retries: %w", lastErr)
		}
		if resp.StatusCode >= 400 {
			return nil, nil, fmt.Errorf("python worker 4xx: %s", resp.Status)
		}
		if decodeErr != nil {
			lastErr = decodeErr
			if attempt < p.MaxRetries {
				time.Sleep(backoff(attempt))
				continue
			}
			return nil, nil, fmt.Errorf("python response decode failed: %w", lastErr)
		}

		return out.Findings, out.Errors, nil
	}

	return nil, nil, fmt.Errorf("python infer failed: %w", lastErr)
}

func backoff(attempt int) time.Duration {
	switch attempt {
	case 0:
		return 200 * time.Millisecond
	case 1:
		return 500 * time.Millisecond
	case 2:
		return 1 * time.Second
	default:
		return 2 * time.Second
	}
}

func (p *PythonClient) Health() error {
	req, _ := http.NewRequest(http.MethodGet, p.BaseURL+"/health", nil)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("python health bad status: %s", resp.Status)
	}
	return nil
}