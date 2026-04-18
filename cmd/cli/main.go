package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"pii-scanner/internal/classify"
	"pii-scanner/internal/detect"
	"pii-scanner/internal/extract"
	"pii-scanner/internal/model"
	"pii-scanner/internal/report"
	"pii-scanner/internal/rules"
	"pii-scanner/internal/scan"
)

const (
	maxTextBytes         = 512 * 1024
	maxSamplesPerFinding = 2
)

func truncateText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

func clampFindings(findings []model.Finding) []model.Finding {
	for i := range findings {
		if len(findings[i].MaskedSamples) > maxSamplesPerFinding {
			findings[i].MaskedSamples = findings[i].MaskedSamples[:maxSamplesPerFinding]
		}
	}
	return findings
}

func main() {
	go func() {
		_ = http.ListenAndServe("127.0.0.1:6060", nil)
	}()

	var root string
	var out string
	var pyURL string
	var shmPath string
	var shmSizeMB int64

	flag.StringVar(&root, "root", ".", "root dir to scan")
	flag.StringVar(&out, "out", "report", "output prefix")
	flag.StringVar(&pyURL, "py-url", "http://127.0.0.1:8081", "python worker base url")
	flag.StringVar(&shmPath, "shm-path", "/tmp/pii_shm.dat", "shared mmap file path")
	flag.Int64Var(&shmSizeMB, "shm-size-mb", 256, "shared mmap size in MB")
	flag.Parse()

	scanID := fmt.Sprintf("scan-%d", time.Now().Unix())
	startedAt := time.Now()

	pyClient := extract.NewPythonClient(pyURL)
	if err := pyClient.Health(); err != nil {
		fmt.Printf("WARN: python worker unavailable at startup: %v\n", err)
	}

	shm, err := extract.NewShmManager(shmPath, shmSizeMB*1024*1024)
	if err != nil {
		fmt.Printf("WARN: shm init failed (%v), continue without shm\n", err)
	}
	if shm != nil {
		defer shm.Close()
	}

	stream, err := report.NewStreamWriter(out+".csv", out+".jsonl")
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	tasks := make(chan model.FileTask, 64)
	results := make(chan model.FileReport, 64)

	go func() {
		_ = scan.Walk(root, scanID, tasks)
	}()

	workers := 4
	if runtime.NumCPU() < workers {
		workers = runtime.NumCPU()
		if workers < 1 {
			workers = 1
		}
	}
	pySlots := make(chan struct{}, 1)

	var shmOK, shmFail, httpFallbackOK, httpFail int64
	var filesTotal, filesParsed, filesFailed int64

	uzCfg := classify.DefaultUZConfig()

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()

			for t := range tasks {
				func(task model.FileTask) {
					defer func() {
						if r := recover(); r != nil {
							rel, _ := filepath.Rel(root, task.Path)
							results <- model.FileReport{
								Path:      rel,
								Format:    task.Format,
								UZLevel:   "NO_PD",
								UZReasons: []string{"анализ прерван: panic recovered"},
								Errors:    []string{fmt.Sprintf("panic recovered: %v", r)},
							}
						}
					}()

					start := time.Now()
					rel, rErr := filepath.Rel(root, task.Path)
					if rErr != nil {
						rel = task.Path
					}

					fr := model.FileReport{
						Path:   rel,
						Format: task.Format,
					}

					if detect.NeedsPython(task.Format) {
						pySlots <- struct{}{}
						findings, pyErrors, pyErr := pyClient.Infer(task.Path, task.Format)
						<-pySlots

						if pyErr != nil {
							atomic.AddInt64(&httpFail, 1)
							fr.Errors = append(fr.Errors, pyErr.Error())
						} else {
							atomic.AddInt64(&httpFallbackOK, 1)
						}
						if len(pyErrors) > 0 {
							fr.Errors = append(fr.Errors, pyErrors...)
						}

						fr.Findings = clampFindings(findings)

						if len(fr.Findings) == 0 && len(fr.Errors) > 0 {
							fr.UZLevel = "NO_PD"
							fr.UZReasons = []string{"анализ не выполнен из-за ошибки обработки файла"}
						} else {
							uz, reasons := classify.DetermineUZWithReasons(fr.Findings, uzCfg)
							fr.UZLevel = uz
							fr.UZReasons = reasons
						}

						fr.DurationMS = time.Since(start).Milliseconds()
						fr.ErrorCategories = classify.CategorizeErrors(fr.Errors)
						fr.Status = classify.DetermineStatus(len(fr.Findings), fr.Errors)
						results <- fr
						return
					}

					text, exErr := extract.ExtractText(task.Path, task.Format)
					if exErr != nil {
						fr.Errors = append(fr.Errors, exErr.Error())
						fr.UZLevel = "NO_PD"
						fr.UZReasons = []string{"ошибка извлечения текста"}
						fr.DurationMS = time.Since(start).Milliseconds()
						results <- fr
						return
					}

					text = truncateText(text, maxTextBytes)

					base := rules.DetectPII(text)
					extra := rules.DetectPIIExtra(text)
					fr.Findings = clampFindings(rules.MergeFindings(base, extra))

					uz, reasons := classify.DetermineUZWithReasons(fr.Findings, uzCfg)
					fr.UZLevel = uz
					fr.UZReasons = reasons

					fr.DurationMS = time.Since(start).Milliseconds()
					results <- fr
				}(t)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		atomic.AddInt64(&filesTotal, 1)

		if r.Status == "failed" {
			atomic.AddInt64(&filesFailed, 1)
		} else {
			atomic.AddInt64(&filesParsed, 1)
		}

		if err := stream.WriteFileReport(r); err != nil {
			fmt.Printf("WARN: stream write failed for %s: %v\n", r.Path, err)
		}
	}

	finishedAt := time.Now()

	summary := model.ScanSummary{
		ScanID:         scanID,
		StartedAt:      startedAt,
		FinishedAt:     finishedAt,
		FilesTotal:     int(filesTotal),
		FilesParsed:    int(filesParsed),
		FilesFailed:    int(filesFailed),
		ShmOK:          atomic.LoadInt64(&shmOK),
		ShmFail:        atomic.LoadInt64(&shmFail),
		HTTPFallbackOK: atomic.LoadInt64(&httpFallbackOK),
		HTTPFail:       atomic.LoadInt64(&httpFail),
	}

	b, _ := json.MarshalIndent(summary, "", "  ")
	_ = os.WriteFile(out+".summary.json", b, 0644)

	fmt.Printf("Done. total=%d parsed=%d failed=%d\n", summary.FilesTotal, summary.FilesParsed, summary.FilesFailed)
	fmt.Printf("Reports: %s.csv, %s.jsonl, %s.summary.json\n", out, out, out)
	fmt.Printf("Python stats: shm_ok=%d shm_fail=%d http_fallback_ok=%d http_fail=%d\n",
		summary.ShmOK, summary.ShmFail, summary.HTTPFallbackOK, summary.HTTPFail)
}