package report

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"

	"pii-scanner/internal/model"
)

type StreamWriter struct {
	mu sync.Mutex

	csvFile *os.File
	csvW    *csv.Writer

	jsonlFile *os.File
}

func NewStreamWriter(csvPath, jsonlPath string) (*StreamWriter, error) {
	cf, err := os.Create(csvPath)
	if err != nil {
		return nil, err
	}
	cw := csv.NewWriter(cf)
	if err := cw.Write([]string{
		"путь", "категории_ПДн", "количество_находок", "УЗ", "причины_УЗ",
		"статус", "категории_ошибок", "формат_файла", "ошибки",
	}); err != nil {
		_ = cf.Close()
		return nil, err
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		_ = cf.Close()
		return nil, err
	}

	jf, err := os.Create(jsonlPath)
	if err != nil {
		_ = cf.Close()
		return nil, err
	}

	return &StreamWriter{csvFile: cf, csvW: cw, jsonlFile: jf}, nil
}

func (s *StreamWriter) WriteFileReport(fr model.FileReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cats := make([]string, 0, len(fr.Findings))
	total := 0
	for _, f := range fr.Findings {
		cats = append(cats, f.Type)
		total += f.Count
	}

	if err := s.csvW.Write([]string{
		fr.Path,
		strings.Join(cats, "|"),
		strconv.Itoa(total),
		fr.UZLevel,
		strings.Join(fr.UZReasons, " | "),
		fr.Status,
		strings.Join(fr.ErrorCategories, "|"),
		fr.Format,
		strings.Join(fr.Errors, " | "),
	}); err != nil {
		return err
	}
	s.csvW.Flush()
	if err := s.csvW.Error(); err != nil {
		return err
	}

	b, err := json.Marshal(fr)
	if err != nil {
		return err
	}
	_, err = s.jsonlFile.Write(append(b, '\n'))
	return err
}

func (s *StreamWriter) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var e1, e2 error
	if s.csvW != nil {
		s.csvW.Flush()
		e1 = s.csvW.Error()
	}
	if s.csvFile != nil {
		if err := s.csvFile.Close(); e1 == nil {
			e1 = err
		}
	}
	if s.jsonlFile != nil {
		e2 = s.jsonlFile.Close()
	}
	if e1 != nil {
		return e1
	}
	return e2
}